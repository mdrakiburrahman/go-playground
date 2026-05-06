#!/bin/bash
#
#
#       Sets up a dev env with all pre-reqs. This script is idempotent, it will
#       only attempt to install dependencies, if not exists.   
#
# ---------------------------------------------------------------------------------------
#

set -e
set -m

echo ""
echo "┌───────────────────────────────────┐"
echo "│ Checking for package dependencies │"
echo "└───────────────────────────────────┘"
echo ""

PACKAGES=""
if ! command -v make &> /dev/null; then PACKAGES="$PACKAGES make"; fi
if ! command -v tree &> /dev/null; then PACKAGES="$PACKAGES tree"; fi
if ! command -v jq &> /dev/null; then PACKAGES="$PACKAGES jq"; fi
if [ ! -z "$PACKAGES" ]; then
    echo "Packages $PACKAGES not found - installing..."
    sudo apt-get update 2>&1 > /dev/null
    sudo DEBIAN_FRONTEND=noninteractive apt-get install -y $PACKAGES 2>&1 > /dev/null
fi

DOCKER_VERSION="5:27.5.1-1~ubuntu.24.04~noble"

if ! [ -x "$(command -v docker)" ]; then
  echo "docker is not installed on your devbox, installing..."
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
  sudo add-apt-repository -y "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
  sudo apt-get update -q
  sudo apt-get install -y apt-transport-https ca-certificates curl
  sudo apt-get install -y --allow-downgrades docker-ce="$DOCKER_VERSION" docker-ce-cli="$DOCKER_VERSION" containerd.io
else
  echo "docker is already installed."
fi

sudo mkdir -p /etc/docker
echo '{"max-concurrent-downloads": 32}' | sudo tee /etc/docker/daemon.json > /dev/null

echo "docker is installed, restarting..."
sudo systemctl restart docker

sudo chmod 666 /var/run/docker.sock
docker container ls
docker ps -q | xargs -r docker kill

echo ""
echo "┌────────────────────────────────────┐"
echo "│ Checking for language dependencies │"
echo "└────────────────────────────────────┘"
echo ""

GO_VERSION="1.22.0"

if ! command -v go &> /dev/null; then
    GO_DOWNLOAD_DIR=`mktemp -d`
    wget -o- "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -P $GO_DOWNLOAD_DIR
    sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf $GO_DOWNLOAD_DIR/go${GO_VERSION}.linux-amd64.tar.gz

    USER_HOME=$(eval echo ~$USER)
    LINES="export PATH=\$PATH:/usr/local/go/bin:$USER_HOME/go/bin"
    echo -e "$LINES" >> ~/.bashrc
    source ~/.bashrc
fi

/usr/local/go/bin/go install github.com/go-delve/delve/cmd/dlv@latest
/usr/local/go/bin/go install github.com/open-telemetry/opentelemetry-collector-contrib/cmd/telemetrygen@latest
/usr/local/go/bin/go install github.com/apache/arrow/go/v13/parquet/cmd/parquet_reader@latest

echo ""
echo "┌────────────────────────┐"
echo "│ Checking for CLI tools │"
echo "└────────────────────────┘"
echo ""

export PATH=$(echo "$PATH" | tr ':' '\n' | grep -v "/mnt/c" | tr '\n' ':' | sed 's/:$//')
AZ_PATH=$(which az 2>/dev/null)
if [[ -z "$AZ_PATH" || "$AZ_PATH" == *"/mnt/c"* ]]; then
  echo "Native Linux Azure CLI not found, installing..."
  curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash
  export PATH="$HOME/bin:$PATH"
  [[ -f "$HOME/.bashrc" ]] && source "$HOME/.bashrc"
else
  echo "Native Linux Azure CLI already installed at: $AZ_PATH"
fi

az account get-access-token --query "expiresOn" -o tsv >/dev/null 2>&1
if [[ $? -ne 0 ]]; then
    echo "az is not logged in, logging in..."
    az login >/dev/null
fi

if ! command -v duckdb &> /dev/null; then
    echo "duckdb not found - installing..."
    curl https://install.duckdb.org | sh
    export PATH='/home/boor/.duckdb/cli/latest':$PATH
fi

echo ""
echo "┌───────────────────────────────┐"
echo "│ Installing VS Code extensions │"
echo "└───────────────────────────────┘"
echo ""

code --install-extension golang.go@0.45.0
code --install-extension ms-vscode.makefile-tools

echo ""
echo "┌──────────┐"
echo "│ Versions │"
echo "└──────────┘"
echo ""

echo "Docker: $(docker --version)"
echo "Go: $(/usr/local/go/bin/go version)"
echo "Azure CLI: $(az version)"