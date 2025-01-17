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
echo "┌────────────────────────────────────┐"
echo "│ Checking for language dependencies │"
echo "└────────────────────────────────────┘"
echo ""

GO_VERSION="1.22.0"

if ! command -v go &> /dev/null; then
    GO_DOWNLOAD_DIR=`mktemp -d`
    wget -o- "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -P $GO_DOWNLOAD_DIR
    sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf $GO_DOWNLOAD_DIR/go${GO_VERSION}.linux-amd64.tar.gz

    LINES="export PATH=\$PATH:/usr/local/go/bin"
    echo -e "$LINES" >> ~/.bashrc
    source ~/.bashrc
fi

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

echo ""
echo "┌──────────┐"
echo "│ Versions │"
echo "└──────────┘"
echo ""

echo "Docker: $(docker --version)"
echo "Go: $(go version)"
echo "Azure CLI: $(az version)"