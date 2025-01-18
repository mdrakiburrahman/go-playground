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
if [ ! -z "$PACKAGES" ]; then
    echo "Packages $PACKAGES not found - installing..."
    sudo apt-get update 2>&1 > /dev/null
    sudo DEBIAN_FRONTEND=noninteractive apt-get install -y $PACKAGES 2>&1 > /dev/null
fi

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

echo ""
echo "┌────────────────────────┐"
echo "│ Checking for CLI tools │"
echo "└────────────────────────┘"
echo ""

if ! command -v docker &> /dev/null; then
    echo "docker not found - installing..."
    curl -sL https://get.docker.com | sudo bash
fi
sudo chmod 666 /var/run/docker.sock

if ! command -v az &> /dev/null; then
    echo "az not found - installing..."
    curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash
fi

echo ""
echo "┌───────────────────────────────┐"
echo "│ Installing VS Code extensions │"
echo "└───────────────────────────────┘"
echo ""

code --install-extension github.copilot
code --install-extension eamodio.gitlens
code --install-extension golang.go
code --install-extension ms-vscode.makefile-tools

echo ""
echo "┌──────────┐"
echo "│ Versions │"
echo "└──────────┘"
echo ""

echo "Docker: $(docker --version)"
echo "Go: $(/usr/local/go/bin/go version)"
echo "Azure CLI: $(az version)"