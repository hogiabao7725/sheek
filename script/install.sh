#!/bin/bash
# install.sh - Build and install sheek

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

echo -e "${BOLD}=== sheek Installation ===${NC}"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed. Please install Go first.${NC}" >&2
    exit 1
fi

# Build sheek
echo "Building sheek..."
cd "$PROJECT_ROOT"
if ! go build -o bin/sheek ./cmd/main-sheek; then
    echo -e "${RED}Error: Failed to build sheek.${NC}" >&2
    exit 1
fi

echo -e "${GREEN}Build successful!${NC}"
echo ""

# Ask for installation location
echo "Where would you like to install sheek?"
echo "1) User-local (~/go/bin or ~/.local/bin) - recommended"
echo "2) System-wide (/usr/local/bin) - requires sudo"
echo ""
read -p "Choose option [1-2]: " install_choice

INSTALL_PATH=""
BINARY_PATH=""

case $install_choice in
    1)
        # User-local installation
        if [ -d "$HOME/.local/bin" ]; then
            INSTALL_PATH="$HOME/.local/bin"
        elif [ -d "$HOME/go/bin" ]; then
            INSTALL_PATH="$HOME/go/bin"
        else
            mkdir -p "$HOME/.local/bin"
            INSTALL_PATH="$HOME/.local/bin"
        fi
        BINARY_PATH="$INSTALL_PATH/sheek"
        # Remove old binary if exists (in case it's busy)
        if [ -f "$BINARY_PATH" ]; then
            # Try to kill any running sheek processes
            if command -v pkill &> /dev/null; then
                pkill -x sheek 2>/dev/null || true
                sleep 0.5
            fi
            rm -f "$BINARY_PATH"
        fi
        cp "$PROJECT_ROOT/bin/sheek" "$BINARY_PATH"
        chmod +x "$BINARY_PATH"
        echo -e "${GREEN}Installed to $BINARY_PATH${NC}"
        ;;
    2)
        # System-wide installation
        if [ "$EUID" -ne 0 ]; then
            echo -e "${YELLOW}System-wide installation requires sudo privileges.${NC}"
            BINARY_PATH="/usr/local/bin/sheek"
            # Remove old binary if exists
            if [ -f "$BINARY_PATH" ]; then
                if command -v pkill &> /dev/null; then
                    sudo pkill -x sheek 2>/dev/null || true
                    sleep 0.5
                fi
                sudo rm -f "$BINARY_PATH"
            fi
            sudo cp "$PROJECT_ROOT/bin/sheek" "$BINARY_PATH"
            sudo chmod +x "$BINARY_PATH"
            echo -e "${GREEN}Installed to $BINARY_PATH${NC}"
        else
            BINARY_PATH="/usr/local/bin/sheek"
            # Remove old binary if exists
            if [ -f "$BINARY_PATH" ]; then
                if command -v pkill &> /dev/null; then
                    pkill -x sheek 2>/dev/null || true
                    sleep 0.5
                fi
                rm -f "$BINARY_PATH"
            fi
            cp "$PROJECT_ROOT/bin/sheek" "$BINARY_PATH"
            chmod +x "$BINARY_PATH"
            echo -e "${GREEN}Installed to $BINARY_PATH${NC}"
        fi
        ;;
    *)
        echo -e "${RED}Invalid choice. Exiting.${NC}"
        exit 1
        ;;
esac

echo ""
echo -e "${BOLD}=== Installation Complete ===${NC}"
echo ""
echo -e "${GREEN}sheek binary installed at: $BINARY_PATH${NC}"
echo ""

# Install default config file
echo -e "${BOLD}=== Config Setup ===${NC}"
CONFIG_DIR="$HOME/.config/sheek"
CONFIG_FILE="$CONFIG_DIR/config.json"
DEFAULT_CONFIG="$PROJECT_ROOT/internal/config/default.json"

if [ ! -f "$DEFAULT_CONFIG" ]; then
    echo -e "${YELLOW}Warning: Default config file not found at $DEFAULT_CONFIG${NC}"
    echo -e "${YELLOW}Skipping config installation.${NC}"
else
    # Create config directory if it doesn't exist
    mkdir -p "$CONFIG_DIR"
    
    # Only copy if config.json doesn't exist (don't overwrite user's config)
    if [ ! -f "$CONFIG_FILE" ]; then
        cp "$DEFAULT_CONFIG" "$CONFIG_FILE"
        echo -e "${GREEN}Default config installed at: $CONFIG_FILE${NC}"
        echo -e "${BLUE}You can edit this file to customize sheek's appearance and behavior.${NC}"
    else
        echo -e "${YELLOW}Config file already exists at: $CONFIG_FILE${NC}"
        echo -e "${BLUE}Skipping config installation to preserve your existing configuration.${NC}"
        echo -e "${BLUE}To reset to defaults, delete the file and run install.sh again.${NC}"
    fi
fi
echo ""

# Detect shell and provide instructions
CURRENT_SHELL="${SHELL##*/}"
if [ -z "$CURRENT_SHELL" ]; then
    CURRENT_SHELL="unknown"
fi

echo -e "${BOLD}=== Shell Integration ===${NC}"
echo ""
echo -e "${BLUE}To use sheek with Ctrl+T keybind:${NC}"
echo ""

if [ "$CURRENT_SHELL" = "zsh" ]; then
    echo -e "${YELLOW}For zsh:${NC}"
    echo "  Add this line to ~/.zshrc:"
    echo -e "    ${BOLD}source $SCRIPT_DIR/keybinds.zsh${NC}"
    echo ""
    echo "  Then run:"
    echo -e "    ${BOLD}source ~/.zshrc${NC}"
    echo ""
    echo "  Or restart your terminal."
elif [ "$CURRENT_SHELL" = "bash" ]; then
    echo -e "${YELLOW}For bash:${NC}"
    echo "  Add this line to ~/.bashrc:"
    echo -e "    ${BOLD}source $SCRIPT_DIR/keybinds.bash${NC}"
    echo ""
    echo "  Then run:"
    echo -e "    ${BOLD}source ~/.bashrc${NC}"
    echo ""
    echo "  Or restart your terminal."
else
    echo -e "${YELLOW}Detected shell: $CURRENT_SHELL${NC}"
    echo ""
    echo -e "${YELLOW}For zsh:${NC}"
    echo "  Add this line to ~/.zshrc:"
    echo -e "    ${BOLD}source $SCRIPT_DIR/keybinds.zsh${NC}"
    echo ""
    echo -e "${YELLOW}For bash:${NC}"
    echo "  Add this line to ~/.bashrc:"
    echo -e "    ${BOLD}source $SCRIPT_DIR/keybinds.bash${NC}"
    echo ""
fi

echo ""
echo "Keybind files location:"
echo "  - $SCRIPT_DIR/keybinds.zsh (for zsh)"
echo "  - $SCRIPT_DIR/keybinds.bash (for bash)"
echo ""
echo -e "${GREEN}After setting up, press Ctrl+T to search command history!${NC}"
