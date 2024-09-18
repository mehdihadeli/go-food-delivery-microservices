#!/usr/bin/env bash

set -eax

# Run the setup-fonts.sh script
echo "Installing Nerd Fonts..."
./setup-fonts.sh

# echo "Installing tools for Go development..."
echo "set permission for running scripts"
./setup-devcontainers.sh