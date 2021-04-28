// This code is available on the terms of the project LICENSE.md file,
// also available online at https://blueoakcouncil.org/license/1.0.0.

package bolt

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	dexdb "decred.org/dcrdex/client/db"
	"decred.org/dcrdex/dex/encode"
	"go.etcd.io/bbolt"
)

type upgradefunc func(db *bbolt.DB) error

// Each database upgrade function should be keyed by the database
// version it upgrades.
var (
	upgrades = [...]upgradefunc{
		// v0 => v1 adds a version key. Upgrades the MatchProof struct to
		// differentiate between server revokes and self revokes.
		v1Upgrade,
		// v1 => v2 adds a MaxFeeRate field to the OrderMetaData, used for match
		// validation.
		v2Upgrade,
		// v2 => v3 adds a tx data field to the match proof.
		v3Upgrade,
	}
	batchDoneErr = errors.New("batched finished")
)

// DBVersion is the latest version of the database that is understood. Databases
// with recorded versions higher than this will fail to open (meaning any
// upgrades prevent reverting to older software).
const DBVersion = uint32(len(upgrades))

func fetchDBVersion(tx *bbolt.Tx) (uint32, error) {
	bucket := tx.Bucket(appBucket)
	if bucket == nil {
		return 0, fmt.Errorf("app bucket not found")
	}

	versionB := bucket.Get(versionKey)
	if versionB == nil {
		return 0, fmt.Errorf("database version not found")
	}

	return intCoder.Uint32(versionB), nil
}

func setDBVersion(tx *bbolt.Tx, newVersion uint32) error {
	bucket := tx.Bucket(appBucket)
	if bucket == nil {
		return fmt.Errorf("app bucket not found")
	}

	return bucket.Put(versionKey, encode.Uint32Bytes(newVersion))
}

// upgradeDB checks whether any upgrades are necessary before the database is
// ready for application usage.  If any are, they are performed.
func (db *BoltDB) upgradeDB() error {
	var version uint32
	version, err := db.getVersion()
	if err != nil {
		return err
	}

	if version > DBVersion {
		return fmt.Errorf("unknown database version %d, "+
			"client recognizes up to %d", version, DBVersion)
	}

	if version == DBVersion {
		// No upgrades necessary.
		return nil
	}

	fmt.Printf("Upgrading database from version %d to %d\n", version, DBVersion)

	// Backup the current version's DB file before processing the upgrades to
	// DBVersion. Note that any intermediate versions are not stored.
	currentFile := filepath.Base(db.Path())
	backupPath := fmt.Sprintf("%s.v%d.bak", currentFile, version) // e.g. dexc.db.v1.bak
	if err = db.backup(backupPath); err != nil {
		return fmt.Errorf("failed to backup DB prior to upgrade: %w", err)
	}

	for i, upgrade := range upgrades[version:] {
		err = doUpgrade(db.DB, upgrade, version+uint32(i)+1)
		if err != nil {
			return err
		}
	}
	return nil

	// All upgrades in a single tx.
	// return db.Update(func(tx *bbolt.Tx) error {
	// 	// Execute all necessary upgrades in order.
	// 	for i, upgrade := range upgrades[version:] {
	// 		err := doUpgrade(tx, upgrade, version+uint32(i)+1)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// 	return nil
	// })
}

// Get the currently stored DB version.
func (db *BoltDB) getVersion() (version uint32, err error) {
	return version, db.View(func(tx *bbolt.Tx) error {
		version, err = getVersionTx(tx)
		return err
	})
}

// Get the uint32 stored in the appBucket's versionKey entry.
func getVersionTx(tx *bbolt.Tx) (uint32, error) {
	bucket := tx.Bucket(appBucket)
	if bucket == nil {
		return 0, fmt.Errorf("appBucket not found")
	}
	versionB := bucket.Get(versionKey)
	if versionB == nil {
		// A nil version indicates a version 0 database.
		return 0, nil
	}
	return intCoder.Uint32(versionB), nil
}

func v1Upgrade(db *bbolt.DB) error {
	const oldVersion = 0

	if err := ensureVersion(db, oldVersion); err != nil {
		return err
	}

	return reloadMatchProofs(db)
}

// v2Upgrade adds a MaxFeeRate field to the OrderMetaData. The upgrade sets the
// MaxFeeRate field for all historical orders to the max uint64. This avoids any
// chance of rejecting a pre-existing active match.
func v2Upgrade(db *bbolt.DB) error {
	const oldVersion = 1

	if err := ensureVersion(db, oldVersion); err != nil {
		return err
	}

	fmt.Println("Adding max fee rates to orders in the database.  This may take a while...")
	start := time.Now()

	// For each order, set a maxfeerate of max uint64.
	maxFeeB := uint64Bytes(^uint64(0))

	// doBatch contains the primary logic for updating max fee in batches.
	// This is done because attempting to migrate in a single database
	// transaction could result in massive memory usage and could potentially
	// crash on many systems due to ulimits.
	const maxEntries = 20000
	var (
		resumeOffset uint32
		numUpdated   uint32
		totalUpdated uint64
	)
	doBatch := func() error {
		return db.Update(func(tx *bbolt.Tx) error {
			master := tx.Bucket(ordersBucket)
			if master == nil {
				return fmt.Errorf("failed to open orders bucket")
			}
			numUpdated = 0
			var numIterated uint32
			return master.ForEach(func(oid, _ []byte) error {

				if numUpdated >= maxEntries {
					return batchDoneErr
				}

				// Skip entries that have already been migrated in previous batches.
				numIterated++
				if numIterated-1 < resumeOffset {
					return nil
				}
				resumeOffset++

				oBkt := master.Bucket(oid)
				if oBkt == nil {
					return fmt.Errorf("order %x bucket is not a bucket", oid)
				}
				numUpdated++
				return oBkt.Put(maxFeeRateKey, maxFeeB)
			})
		})
	}

	// Update all entries in batches.
	for {
		err := doBatch()
		totalUpdated += uint64(numUpdated)
		if err != nil {
			if errors.Is(err, batchDoneErr) {
				fmt.Printf("Updated %d entries (%d total)\n", numUpdated, totalUpdated)
				continue
			}
			return err
		}
		break
	}

	elapsed := time.Since(start).Round(time.Millisecond)
	fmt.Printf("Done updating orders.  Total entries: %d in %v\n",
		totalUpdated, elapsed)
	return nil
}

func v3Upgrade(db *bbolt.DB) error {
	const oldVersion = 2

	if err := ensureVersion(db, oldVersion); err != nil {
		return err
	}

	// Upgrade the match proof. We just have to retrieve and re-store the
	// buckets. The decoder will recognize the the old version and add the new
	// field.
	return reloadMatchProofs(db)
}

func ensureVersion(db *bbolt.DB, ver uint32) error {
	return db.View(func(tx *bbolt.Tx) error {
		dbVersion, err := fetchDBVersion(tx)
		if err != nil {
			return fmt.Errorf("error fetching database version: %w", err)
		}

		if dbVersion != ver {
			return fmt.Errorf("wrong version for upgrade. expected %d, got %d", ver, dbVersion)
		}
		return nil
	})
}

func reloadMatchProofs(db *bbolt.DB) error {

	fmt.Println("Reloading match proofs in the database.  This may take a while...")
	start := time.Now()

	// doBatch contains the primary logic for updating match proofs in
	// batches.  This is done because attempting to migrate in a single database
	// transaction could result in massive memory usage and could potentially
	// crash on many systems due to ulimits.
	const maxEntries = 20000
	var (
		resumeOffset uint32
		numUpdated   uint32
		totalUpdated uint64
	)
	doBatch := func() error {
		return db.Update(func(tx *bbolt.Tx) error {
			matches := tx.Bucket(matchesBucket)
			numUpdated = 0
			var numIterated uint32
			return matches.ForEach(func(k, _ []byte) error {

				if numUpdated >= maxEntries {
					return batchDoneErr
				}

				// Skip entries that have already been migrated in previous batches.
				numIterated++
				if numIterated-1 < resumeOffset {
					return nil
				}
				resumeOffset++

				mBkt := matches.Bucket(k)
				if mBkt == nil {
					return fmt.Errorf("match %x bucket is not a bucket", k)
				}
				proofB := mBkt.Get(proofKey)
				if len(proofB) == 0 {
					return fmt.Errorf("empty match proof")
				}
				proof, err := dexdb.DecodeMatchProof(proofB)
				if err != nil {
					return fmt.Errorf("error decoding proof: %w", err)
				}
				err = mBkt.Put(proofKey, proof.Encode())
				if err != nil {
					return fmt.Errorf("error re-storing match proof: %w", err)
				}

				numUpdated++
				return nil
			})
		})
	}

	// Update all entries in batches.
	for {
		err := doBatch()
		totalUpdated += uint64(numUpdated)
		if err != nil {
			if errors.Is(err, batchDoneErr) {
				fmt.Printf("Updated %d entries (%d total)\n", numUpdated, totalUpdated)
				continue
			}
			return err
		}
		break
	}

	elapsed := time.Since(start).Round(time.Millisecond)
	fmt.Printf("Done updating match proofs.  Total entries: %d in %v\n",
		totalUpdated, elapsed)
	return nil
}

func doUpgrade(db *bbolt.DB, upgrade upgradefunc, newVersion uint32) error {
	err := upgrade(db)
	if err != nil {
		return fmt.Errorf("error upgrading DB: %v", err)
	}
	// Persist the database version.
	return db.Update(func(tx *bbolt.Tx) error {
		err := setDBVersion(tx, newVersion)
		if err != nil {
			return fmt.Errorf("error setting DB version: %v", err)
		}
		return nil
	})
}
