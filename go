#!/usr/bin/env bash

set -e
set -o nounset
set -o pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
GOPASS_VERSION="1.8.5"

# shellcheck source=./go.helpers
source "${SCRIPT_DIR}/go.helpers"

requirements=(gpg gopass)

goal_check-binaries() {
  echob "*** Checking required binaries ***\\n"
  for cmd in "${requirements[@]}"; do
    check "$cmd" "$cmd is not installed\\nRun ./go install-binaries" "$cmd is installed"
  done
}

goal_check-version() {
  echob "*** Checking versions of the binaries ***\\n"

  if [[ "$(gopass -v)" != "gopass ${GOPASS_VERSION}"* ]]; then
    echo_error "the version of gopass is too old"
    exit 1
  else
    echo_check "gopass has a compatible version"
  fi
}

goal_install-binaries() {
  echob "*** Install required binaries ***\\n"
  for cmd in "${requirements[@]}"; do
    if ! silent_check "$cmd" ; then
      if [ "$cmd" = "gpg" ]; then
        cmd="gpg2"
      fi

      echob "**** Installing ${cmd} ****"
      brew install "${cmd}"
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

goal_check-path() {
  echob "*** Checking if path is configured correctly ***\\n"

  for store in $(gopass config path); do
    if echo "$store" | grep -q gpgcli-noop ; then
      echo_error "Path config is set to noop\\n"
      echo "If you run 'gopass config path' you will see that one of your backends is configured as gpgcli-noop, which won't work. It should be set to gpgcli-gitcli instead. Check ~/.config/gopass/config.yml"
      exit 1
    fi
  done
}

goal_verify() {
  goal_check-binaries
  goal_check-version
  goal_check-tty
  goal_check-secret-key
  goal_check-encryption
  goal_check-path
}

TARGET=${1:-}
if type -t "goal_${TARGET}" &>/dev/null; then
  goal_"${TARGET}" "${@:2}"
else
  echo "usage: $0 <goal>

goal:
    check-binaries           -- checks if you have the required binaries installed
    check-version            -- checks if the versions of the required binaries match
    install-binaries         -- install the binaries required for gopass

    check-tty                -- checks whether you have the _tty_ set up correctly to use GPG
    check-secret-key         -- checks that you have a GPG secret key
    check-encryption         -- checks that you can encrypt and decrypt a text

    check-path               -- checks that the path to the gopass store is configured correctly

    verify                   -- verifies that you can use gopass correctly
"
  exit 1
fi
