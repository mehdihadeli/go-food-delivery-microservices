#!/bin/bash

# Download MesloLGM Nerd Font
wget https://github.com/ryanoasis/nerd-fonts/releases/download/v3.0.2/Meslo.zip -O MesloLGM.zip

# Extract the font files
unzip MesloLGM.zip -d MesloLGM

# Create the fonts directory if it doesn't exist
mkdir -p ~/.local/share/fonts

# Move the font files to the fonts directory
mv MesloLGM/*.ttf ~/.local/share/fonts/

# Update the font cache
fc-cache -fv

# Clean up
rm -rf MesloLGM.zip MesloLGM

# Verify installation
fc-list | grep "MesloLGM"

echo "Font setup completed."