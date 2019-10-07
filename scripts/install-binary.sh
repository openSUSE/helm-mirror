#!/bin/bash

# Shamelessly copied from https://github.com/databus23/helm-diff

version=$1

PROJECT_NAME="helm-mirror"
PROJECT_GH="openSUSE/$PROJECT_NAME"

: ${HELM_PLUGIN_PATH:="$(helm home --debug=false)/plugins/helm-mirror"}

# Convert the HELM_PLUGIN_PATH to unix if cygpath is
# available. This is the case when using MSYS2 or Cygwin
# on Windows where helm returns a Windows path but we
# need a Unix path

if type cygpath > /dev/null 2>&1; then
  HELM_PLUGIN_PATH=$(cygpath -u $HELM_PLUGIN_PATH)
fi

if [[ $SKIP_BIN_INSTALL == "1" ]]; then
  echo "Skipping binary install"
  exit
fi

# initArch discovers the architecture for this system.
initArch() {
  ARCH=$(uname -m)
  case $ARCH in
    armv5*) ARCH="armv5";;
    armv6*) ARCH="armv6";;
    armv7*) ARCH="armv7";;
    aarch64) ARCH="arm64";;
    x86) ARCH="386";;
    x86_64) ARCH="amd64";;
    i686) ARCH="386";;
    i386) ARCH="386";;
  esac
}

# initOS discovers the operating system for this system.
initOS() {
  OS=$(echo `uname`|tr '[:upper:]' '[:lower:]')

  case "$OS" in
    # Msys support
    msys*) OS='windows';;
    # Minimalist GNU for Windows
    mingw*) OS='windows';;
    darwin) OS='macos';;
  esac
}

# verifySupported checks that the os/arch combination is supported for
# binary builds.
verifySupported() {
  local supported="linux-amd64"
  if ! echo "${supported}" | grep -q "${OS}-${ARCH}"; then
    echo "No prebuild binary for ${OS}-${ARCH}."
    exit 1
  fi

  if ! type "curl" > /dev/null && ! type "wget" > /dev/null; then
    echo "Either curl or wget is required"
    exit 1
  fi
}

# getDownloadURL checks the latest available version.
getDownloadURL() {
  if [ -n "$version" ] && [ "$version" != 'master' ]; then
    DOWNLOAD_URL="https://github.com/$PROJECT_GH/releases/download/$version/helm-mirror-$OS.tgz"
  else
    # Use the GitHub API to find the download url for this project.
    local url="https://github.com/$PROJECT_GH/releases/latest"
    if type "curl" > /dev/null; then
      version=$(curl -s $url | grep -o -E 'v.+\"' | awk '{split($0,a,/"/); print a[1]}')
    elif type "wget" > /dev/null; then
      version=$(wget -qSO- $url --max-redirect 0 2>&1 | grep Location: | awk '{split($0,a,/\//); print a[8]}')
    fi
    DOWNLOAD_URL="https://github.com/$PROJECT_GH/releases/download/$version/helm-mirror-$OS.tgz"
  fi
  echo "using download url ${DOWNLOAD_URL}"
}

# downloadFile downloads the latest binary package and also the checksum
# for that binary.
downloadFile() {
  PLUGIN_TMP_FILE="/tmp/${PROJECT_NAME}.tgz"
  echo "Downloading $DOWNLOAD_URL"
  if type "curl" > /dev/null; then
    curl -L "$DOWNLOAD_URL" -o "$PLUGIN_TMP_FILE"
  elif type "wget" > /dev/null; then
    wget -q -O "$PLUGIN_TMP_FILE" "$DOWNLOAD_URL"
  fi
}

# installFile verifies the SHA256 for the file, then unpacks and
# installs it.
installFile() {
  HELM_TMP="/tmp/$PROJECT_NAME"
  rm -rf "$HELM_TMP"
  mkdir -p "$HELM_TMP"
  tar xf "$PLUGIN_TMP_FILE" -C "$HELM_TMP" --strip-components=1
  echo "Preparing to install into ${HELM_PLUGIN_PATH}"
  mkdir -p "$HELM_PLUGIN_PATH/bin"
  pushd "$HELM_TMP"
  cp -r $HELM_TMP/* "$HELM_PLUGIN_PATH"
  popd
}

# fail_trap is executed if an error occurs.
fail_trap() {
  result=$?
  if [ "$result" != "0" ]; then
    echo "Failed to install $PROJECT_NAME"
    echo "For support, go to https://github.com/openSUSE/helm-mirror."
  fi
  exit $result
}

# testVersion tests the installed client to make sure it is working.
testVersion() {
  set +e
  echo "$PROJECT_NAME installed into $HELM_PLUGIN_PATH"
  $HELM_PLUGIN_PATH/bin/helm-mirror version
  set -e
}

# Execution

#Stop execution on any error
trap "fail_trap" EXIT
set -e
initArch
initOS
verifySupported
getDownloadURL
downloadFile
installFile
testVersion
