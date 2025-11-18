# sheek - Ctrl+T keybinding for zsh
# Add this to your ~/.zshrc: source /path/to/sheek/script/keybinds.zsh

# Find sheek binary - try multiple locations
SHEEK_BIN=""
if [ -f "$HOME/.local/bin/sheek" ]; then
    SHEEK_BIN="$HOME/.local/bin/sheek"
elif [ -f "$HOME/go/bin/sheek" ]; then
    SHEEK_BIN="$HOME/go/bin/sheek"
elif [ -f "/usr/local/bin/sheek" ]; then
    SHEEK_BIN="/usr/local/bin/sheek"
elif [ -f "/usr/bin/sheek" ]; then
    SHEEK_BIN="/usr/bin/sheek"
elif command -v sheek &> /dev/null; then
    SHEEK_BIN=$(command -v sheek)
fi

if [ -z "$SHEEK_BIN" ] || [ ! -f "$SHEEK_BIN" ]; then
    echo "Warning: sheek binary not found. Please install sheek first." >&2
    echo "Run: cd /path/to/sheek && ./script/install.sh" >&2
    return 1
fi

# Verify binary is executable
if [ ! -x "$SHEEK_BIN" ]; then
    echo "Warning: sheek binary is not executable: $SHEEK_BIN" >&2
    return 1
fi

sheek-widget() {
    local selected
    # Preserve TERM and COLORTERM for color support
    # Capture stdout for command, but allow stderr to show (for command preview)
    selected=$(TERM="${TERM:-xterm-256color}" COLORTERM="${COLORTERM:-truecolor}" $SHEEK_BIN)
    
    if [ -n "$selected" ]; then
        # Remove trailing newline if any
        selected=$(echo "$selected" | tr -d '\n')
        BUFFER="$selected"
        CURSOR=${#BUFFER}
    fi
    
    zle redisplay
}

# Shell function to run sheek and put selected command in prompt buffer
# This allows running 'sheek' directly from command line (like fzf)
sheek() {
    local selected
    # Preserve TERM and COLORTERM for color support
    selected=$(TERM="${TERM:-xterm-256color}" COLORTERM="${COLORTERM:-truecolor}" $SHEEK_BIN)
    
    if [ -n "$selected" ]; then
        # Remove trailing newline if any
        selected=$(echo "$selected" | tr -d '\n')
        # Use print -z to put command in prompt buffer (like fzf does)
        # This works when called from command line, not just from zle widget
        print -z "$selected"
    fi
}

# Register widget and bind key
zle -N sheek-widget
# Unbind any existing Ctrl+T binding first
bindkey -r '^T' 2>/dev/null || true
# Bind Ctrl+T to sheek widget
bindkey '^T' sheek-widget

