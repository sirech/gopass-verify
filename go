#!/usr/bin/env bash

set -e
set -o nounset
set -o pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# shellcheck source=./go.helpers
source "${SCRIPT_DIR}/go.helpers"

requirements=(gpg gopass)

goal_check-binaries() {
  echob "*** Checking required binaries ***"
  for cmd in "${requirements[@]}"; do
    check "$cmd" "$cmd is not installed" "$cmd is installed"
  done
}

TARGET=${1:-}
if type -t "goal_${TARGET}" &>/dev/null; then
  goal_"${TARGET}" "${@:2}"
else
  echo "usage: $0 <goal>

goal:
    check-binaries           -- checks if you have the required binaries installed
"
  exit 1
fi
