#!/usr/bin/env bash

set -e
set -o nounset
set -o pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# shellcheck source=./go.helpers
source "${SCRIPT_DIR}/go.helpers"

requirements=(gpg gopass)
declare -A deps=(["gpg"]="gpg2")

goal_check-binaries() {
  echob "*** Checking required binaries ***\\n"
  for cmd in "${requirements[@]}"; do
    check "$cmd" "$cmd is not installed" "$cmd is installed"
  done
}

goal_install-binaries() {
  echob "*** Install required binaries ***\\n"
  for cmd in "${requirements[@]}"; do
    if  ! silent_check "$cmd" ; then
      local dep=${deps[$cmd]:-$cmd}

      echob "**** Installing ${dep} ****"
      brew install "${dep}"
    fi
  done
}

goal_check-tty() {
  echob "*** Checking GPG_TTY ***\\n"

  if [ -z ${GPG_TTY+x} ]; then
    echo_error "The variable GPG_TTY is not set\\n"
    echo "Not having the GPG_TTY variable set up will lead to hard-to-debug errors when calling gpg. To set it, do:

    export GPG_TTY=\$(tty)

You should add this line to your ~/.bashrc file to ensure this is set every session.
"
    exit 1
  else
    echo_check "GPG_TTY is set"
  fi
}

goal_verify() {
  goal_check-binaries
  goal_check-tty
}

TARGET=${1:-}
if type -t "goal_${TARGET}" &>/dev/null; then
  goal_"${TARGET}" "${@:2}"
else
  echo "usage: $0 <goal>

goal:
    check-binaries           -- checks if you have the required binaries installed
    install-binaries         -- install the binaries required for gopass

    check-tty                -- checks whether you have the _tty_ set up correctly to use GPG

    verify                   -- verifies that you can use gopass correctly
"
  exit 1
fi
