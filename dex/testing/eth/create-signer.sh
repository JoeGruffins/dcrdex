#!/bin/sh
# Script for creating eth signers.
set -e

# The following are required script arguments
TMUX_WIN_ID=$1
NAME=$2
PASSWORD=$3

MASTER_PASSWORD="abcdefghij"
NAME_CLEF="${NAME}-clef"
GROUP_DIR="${NODES_ROOT}/${NAME}"
KEYSTORE_DIR="${GROUP_DIR}/node/keystore"
CLEF_DIR="${GROUP_DIR}/clef"
RULES="${GROUP_DIR}/rules.js"
mkdir -p "${CLEF_DIR}"

# rules.js allows us to used signer functions based on certain curcumstances.
cat > "${RULES}" <<EOF
function ApproveListing() {
    return "Approve"
}
function ApproveSignData(r){
  if(r.message.indexOf("${PASSWORD}") >= 0){
    return "Approve"
  }
}
EOF
# RULES_SHA is the SHA256 of the above file.
RULES_SHA=$(sha256sum "${RULES}" | awk '{print $1;}')
echo "${RULES_SHA}"

# Create a tmux window.
tmux new-window -t "$TMUX_WIN_ID" -n "${NAME_CLEF}"
tmux send-keys -t "$TMUX_WIN_ID" "set +o history" C-m
tmux send-keys -t "$TMUX_WIN_ID" "cd ${CLEF_DIR}" C-m

echo "Initializing simnet ${NAME} clef"
tmux send-keys -t "$TMUX_WIN_ID" "clef --configdir ${CLEF_DIR} init; tmux wait-for -S ${NAME_CLEF}" C-m
sleep 0.5
tmux send-keys -t "$TMUX_WIN_ID" ok C-m
sleep 0.5
tmux send-keys -t "$TMUX_WIN_ID" ${MASTER_PASSWORD} C-m
sleep 0.5
tmux send-keys -t "$TMUX_WIN_ID" ${MASTER_PASSWORD} C-m
tmux wait-for "${NAME_CLEF}"

echo "Loading simnet ${NAME} clef rules sha"
tmux send-keys -t "$TMUX_WIN_ID" "clef --configdir ${CLEF_DIR} attest \"${RULES_SHA}\"; tmux wait-for -S ${NAME_CLEF}" C-m
sleep 0.5
tmux send-keys -t "$TMUX_WIN_ID" ok C-m
sleep 0.5
tmux send-keys -t "$TMUX_WIN_ID" ${MASTER_PASSWORD} C-m
tmux wait-for "${NAME_CLEF}"

echo "Starting simnet ${NAME} clef"
tmux send-keys -t "$TMUX_WIN_ID" "clef --configdir ${CLEF_DIR} --chainid 42 --keystore ${KEYSTORE_DIR} --rules ${RULES}" C-m
sleep 0.5
tmux send-keys -t "$TMUX_WIN_ID" ok C-m
sleep 0.5
tmux send-keys -t "$TMUX_WIN_ID" ${MASTER_PASSWORD} C-m
sleep 1
