#!/bin/sh
# Script for creating eth nodes.
set -e

# The following are required script arguments. Non-mining nodes only require
# up to IS_MINER.
TMUX_NODE_WIN_ID=$1
TMUX_SIGNER_WIN_ID=$2
CLEF_PASSWORD=$3
NAME=$4
NODE_PORT=$5
NODE_KEY=$6
# IS_MINER is 0 or 1.
IS_MINER=$7
CHAIN_ADDRESS=$8
CHAIN_PASSWORD=$9
CHAIN_ADDRESS_JSON=${10}
CHAIN_ADDRESS_JSON_FILE_NAME=${11}
ADDRESS=${12}
ADDRESS_PASSWORD=${13}
ADDRESS_JSON=${14}
ADDRESS_JSON_FILE_NAME=${15}

GROUP_DIR="${NODES_ROOT}/${NAME}"
NODE_DIR="${GROUP_DIR}/node"
MINE_JS="${GROUP_DIR}/mine.js"
SEND_JS="${GROUP_DIR}/send.js"
SIGNER="${GROUP_DIR}/clef/clef.ipc"
mkdir -p "${NODE_DIR}"

# Write node ctl script.
cat > "${NODES_ROOT}/harness-ctl/${NAME}" <<EOF
#!/bin/sh
geth --datadir="${NODE_DIR}" \$*
EOF
chmod +x "${NODES_ROOT}/harness-ctl/${NAME}"

if [ "${IS_MINER}" = 1 ]; then
  # Write mining javascript.
  # NOTE: This sometimes mines more than one block. It is a race. This returns
  # the number of blocks mined within the lifespan of the function, but one more
  # MAY be mined after returning.
  cat > "${MINE_JS}" <<EOF
function mine() {
  blkN = eth.blockNumber;
  miner.start();
  miner.stop();
  admin.sleep(1.1);
  return eth.blockNumber - blkN;
}
EOF

  # Write mine script.
  cat > "${NODES_ROOT}/harness-ctl/mine-${NAME}" <<EOF
#!/bin/sh
case \$1 in
    ''|*[!0-9]*)  ;;
    *) NUM=\$1 ;;
esac
for i in \$(seq \$NUM) ; do
  ./${NAME} attach --preload "${MINE_JS}" --exec 'mine()'
done
EOF
  chmod +x "${NODES_ROOT}/harness-ctl/mine-${NAME}"

  cat > "${SEND_JS}" <<EOF
function send(dest, amt) {
  eth.sendTransaction({from: ${ADDRESS}, to: \$dest, amount: \$amt})
}
EOF

  cat > "${NODES_ROOT}/harness-ctl/${NAME}-send-to" <<EOF
#!/bin/sh
case \$2 in
    ''|*[!0-9]*)  ;;
    *) AMT=\$2 ;;
esac
DEST=\$1
./${NAME} attach --preload "${SEND_JS}" --exec 'send(DEST, AMT)'
EOF
  chmod +x "${NODES_ROOT}/harness-ctl/${NAME}-send-to"

  # Write password file to unlock accounts later.
  cat > "${GROUP_DIR}/password" <<EOF
$CHAIN_PASSWORD
$ADDRESS_PASSWORD
EOF
fi

# Create a node tmux window.
tmux new-window -t "$TMUX_NODE_WIN_ID" -n "${NAME}"
tmux send-keys -t "$TMUX_NODE_WIN_ID" "set +o history" C-m
tmux send-keys -t "$TMUX_NODE_WIN_ID" "cd ${NODE_DIR}" C-m

# Create and wait for a node initiated with a predefined genesis json.
echo "Creating simnet ${NAME} node"
tmux send-keys -t "$TMUX_NODE_WIN_ID" "${NODES_ROOT}/harness-ctl/${NAME} init "\
	"$GENESIS_JSON_FILE_LOCATION; tmux wait-for -S ${NAME}" C-m
tmux wait-for "${NAME}"

if [ "${IS_MINER}" = 1 ]; then
  # Create two accounts. The first is used to mine blocks. The second contains
  # funds.
  echo "Creating account"
  cat > "${NODE_DIR}/keystore/$CHAIN_ADDRESS_JSON_FILE_NAME" <<EOF
$CHAIN_ADDRESS_JSON
EOF
  cat > "${NODE_DIR}/keystore/$ADDRESS_JSON_FILE_NAME" <<EOF
$ADDRESS_JSON
EOF
fi

# The node key lets us control the enode address value.
echo "Setting node key"
cat > "${NODE_DIR}/geth/nodekey" <<EOF
$NODE_KEY
EOF

if [ "${IS_MINER}" = 1 ]; then
  # Start the eth node with both accounts unlocked, listening restricted to
  # localhost, and syncmode set to full.
  echo "Starting simnet ${NAME} node"
  tmux send-keys -t "$TMUX_NODE_WIN_ID" "${NODES_ROOT}/harness-ctl/${NAME} --port " \
	  "${NODE_PORT} --nodiscover --unlock ${CHAIN_ADDRESS},${ADDRESS} " \
	  "--password ${GROUP_DIR}/password --miner.etherbase ${CHAIN_ADDRESS} " \
	  "--syncmode full --netrestrict 127.0.0.1/32" C-m
else
  # Create the signer.
  "${HARNESS_DIR}/create-signer.sh" "$TMUX_SIGNER_WIN_ID" "${NAME}" "${CLEF_PASSWORD}"

  sleep 1

  # Start the eth wtih listening restricted to localhost, syncmode set to full,
  # and using the signer.
  echo "Starting simnet ${NAME} node"
  tmux send-keys -t "$TMUX_NODE_WIN_ID" "${NODES_ROOT}/harness-ctl/${NAME} --port " \
	  "${NODE_PORT} --nodiscover --syncmode full --netrestrict 127.0.0.1/32 " \
	  "--signer ${SIGNER}" C-m
fi
