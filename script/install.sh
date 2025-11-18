#!/bin/bash
# install.sh - Build and install sheek

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "=== sheek Installation ==="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go first." >&2
    exit 1
fi

# Build sheek
echo "Building sheek..."
cd "$PROJECT_ROOT"
if ! go build -o bin/sheek ./cmd/main-sheek; then
    echo "Error: Failed to build sheek." >&2
    exit 1
fi

echo "Build successful!"
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
        echo "Installed to $BINARY_PATH"
        ;;
    2)
        # System-wide installation
        if [ "$EUID" -ne 0 ]; then
            echo "System-wide installation requires sudo privileges."
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
            echo "Installed to $BINARY_PATH"
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
            echo "Installed to $BINARY_PATH"
        fi
        ;;
    *)
        echo "Invalid choice. Exiting."
        exit 1
        ;;
esac

echo ""
echo "=== Keybind files ==="
echo "Keybind files are ready at:"
echo "  - $SCRIPT_DIR/keybinds.zsh (for zsh)"
echo "  - $SCRIPT_DIR/keybinds.bash (for bash)"

echo ""
echo "=== Installation Complete ==="
echo ""
echo "To use sheek with Ctrl+T:"
echo ""
echo "For zsh:"
echo "  Add to ~/.zshrc:"
echo "    source $SCRIPT_DIR/keybinds.zsh"
echo ""
echo "For bash:"
echo "  Add to ~/.bashrc:"
echo "    source $SCRIPT_DIR/keybinds.bash"
echo ""
echo "Then run: source ~/.zshrc  (or source ~/.bashrc)"
echo ""
echo "sheek binary installed at: $BINARY_PATH"

