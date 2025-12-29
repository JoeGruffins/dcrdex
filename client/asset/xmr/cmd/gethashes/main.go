package main

import (
	"encoding/json"
	"fmt"

	"decred.org/dcrdex/client/asset/xmr/toolsdl"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

func run() error {
	d := toolsdl.Download{}
	hashes, err := d.DownloadHashesFile()
	if err != nil {
		return err
	}
	fmt.Println(hashes)
	if err = d.GetHashedZips(hashes); err != nil {
		return err
	}
	if err = d.CheckHashedZips(); err != nil {
		return err
	}
	b, err := json.MarshalIndent(d.Hashes(), "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}
