#!/usr/bin/env bash

set -e
set -o nounset
set -o pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
GOPASS_VERSION="1.9.2"

# shellcheck source=./go.helpers
source "${SCRIPT_DIR}/go.helpers"

requirements=(gpg git)

tools_mac=(gpg git)
tools_linux=(gnupg git)

example_text="Lorem ipsum dolor sit amet, ..."

goal_check-tools() {
  echo ""
  echob "*** Checking required tools ***\\n"
  for cmd in "${requirements[@]}"; do
    check "$cmd" "$cmd is not installed\\nRun ./go install-tools" "$cmd is installed"
  done
}

goal_check-gopass() {
  echo ""
  echob "*** Checking version of Gopass ***\\n"

  if ! type "gopass" > /dev/null 2>&1; then
    echo_error "Gopass is missing.\\n"
    echo "Run:"
    echo ""
    echob "./go install-gopass"
    echo ""
    exit 1
  

  elif [[ "$(gopass -v)" != "gopass ${GOPASS_VERSION}"* ]]; then
    echo_error "The version of Gopass is too old"
    exit 1

  else
    echo_check "Gopass has a compatible version"
    echo ""
  fi
}

goal_install-tools() {
  echob "*** Installing required tools ***\\n"
  if [[ "$OSTYPE" == "darwin"* ]]; then
    goal_install-tools-mac
  elif [[ "$OSTYPE" == "linux-gnu" ]]; then
    goal_install-tools-linux
  fi
}

goal_install-tools-mac(){
  for cmd in "${tools_mac[@]}"; do
  if ! silent_check "$cmd" ; then
      if [ "$cmd" = "gpg" ]; then
        cmd="gpg2"
      fi
      echob "**** Installing ${cmd} ****"
      brew install "${cmd}"
    fi
  done
}

goal_install-tools-linux(){
  for cmd in "${tools_linux[@]}"; do
    sudo apt install "${cmd}"
  done
}

goal_install-gopass() {
  echob "*** Installing Gopass ***\\n"
  if [[ "$OSTYPE" == "darwin"* ]]; then
    goal_install-gopass-mac
  elif [[ "$OSTYPE" == "linux-gnu" ]]; then
    goal_install-gopass-linux
  fi
}

goal_install-gopass-linux(){
  echob "Downloading Gopass .deb for Linux"
  curl -Lo ./gopass-${GOPASS_VERSION}-linux-amd64.deb https://github.com/gopasspw/gopass/releases/download/v${GOPASS_VERSION}/gopass-${GOPASS_VERSION}-linux-amd64.deb
  echob "Downloading Gopass SHA256SUMS"
  curl -Lo ./gopass_${GOPASS_VERSION}_SHA256SUMS https://github.com/gopasspw/gopass/releases/download/v${GOPASS_VERSION}/gopass_${GOPASS_VERSION}_SHA256SUMS
  echob "Checking SHA256"
  grep linux-amd64.deb gopass_${GOPASS_VERSION}_SHA256SUMS | tee /proc/self/fd/2 | sha256sum --check -
  echob "Installing Gopass"
  sudo dpkg -i gopass-${GOPASS_VERSION}-linux-amd64.deb
  echob "Removing temporary files"
  rm -rf gopass_${GOPASS_VERSION}_*
  gopass -version
  echob "Gopass is successfully installed"
}

goal_install-gopass-mac(){
  brew install gopass
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
    echo_check "GPG_TTY is set\\n"
  fi
}

goal_check-secret-key() {
  echob "*** Checking if there is a secret key ***\\n"

  if [ -z "$(gpg --list-secret-keys)" ]; then
    echo_error "There is no secret key\\n"
    echo "Not having a secret key means that you will not be able to read secrets from gopass, even if you are added as a recipient to it"
    exit 1
  else
    echo_check "Secret key is present\\n"
  fi
}

goal_check-encryption() {
  echob "*** Checking if your GPG key can encrypt/decrypt an example text ***\\n"

  read -rp "Which key should be used? (Enter the name, the key ID or the email of your GPG key): " keyid

  if echo "$example_text" | gpg --encrypt --recipient "${keyid}" | gpg --decrypt ; then
    echo ""
    echo_check "Check successful\\n"
  else
    echo ""
    echo_error "Could not encrypt and decrypt text\\n"
    echo "This usually means that the key is not imported correctly, or that there is no secret key, or that the key is not trusted. Gopass will not be able to decrypt secrets until you can use your key"
    exit 1
  fi
}

goal_check-path() {
  echob "*** Checking if path is configured correctly ***\\n"

  if [ -z "$(gopass config path)" ]; then
    echo_error "No Gopass store has been initialized"

  else 
    for store in $(gopass config path); do
      if echo "$store" | grep -q gpgcli-noop ; then
        echo_error "Path config is set to noop\\n"
        echo "If you run 'gopass config path' you will see that one of your backends is configured as gpgcli-noop, which won't work. It should be set to gpgcli-gitcli instead. Check ~/.config/gopass/config.yml"
        exit 1
      elif echo "$store" | grep -q gpgcli-gitcli ; then
        echo_check "Correct backend gpgcli-gitcli is in use"
        echo ""
      fi
    done
  fi
}

goal_verify() {
  goal_check-tools
  goal_check-gopass
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
    check-tools              -- Check if you have the required tools installed
    check-gopass             -- Check if Gopass is correctly installed

    install-tools            -- Install the required tools for Gopass
    install-gopass           -- Install Gopass

    check-tty                -- Check whether you have the _tty_ set up correctly to use GPG
    check-secret-key         -- Check that you have a GPG secret key
    check-encryption         -- Check that you can encrypt and decrypt a text

    check-path               -- Check that the path to the Gopass store is configured correctly

    verify                   -- Verify that you can use Gopass correctly
"
  exit 1
fi
