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
    check "$cmd" "$cmd is not installed\\nRun ./go install-binaries" "$cmd is installed"
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

  if [ -z "${GPG_TTY+x}" ]; then
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

goal_check-secret-key() {
  echob "*** Checking if there is a secret key ***\\n"

  if [ -z "$(gpg --list-secret-keys)" ]; then
    echo_error "There is no secret key\\n"
    echo "Not having a secret key means that you will not be able to read secrets from gopass, even if you are added as a recipient to it"
    exit 1
  else
    echo_check "Secret key is present"
  fi
}

goal_check-encryption() {
  echob "*** Checking if the gpg key can encrypt/decrypt a text ***\\n"

  read -rp "Which key should be used? (Enter the name, the id or the email): " keyid

  if echo 'example' | gpg --encrypt --recipient "${keyid}" | gpg --decrypt ; then
    echo_check "\\nCheck successful"
  else
    echo ''
    echo_error "Could not encrypt and decrypt text\\n"
    echo "This usually means that the key is not imported correctly, or that there is no secret key, or that the key is not trusted. Gopass will not be able to decrypt secrets until you can use your key"
    exit 1
  fi
}

goal_verify() {
  goal_check-binaries
  goal_check-tty
  goal_check-secret-key
  goal_check-encryption
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
    check-secret-key         -- checks that you have a GPG secret key
    check-encryption         -- checks that you can encrypt and decrypt a text

    verify                   -- verifies that you can use gopass correctly
"
  exit 1
fi
