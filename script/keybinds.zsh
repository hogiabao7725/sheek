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
    return 1
fi

sheek-widget() {
    local selected
    selected=$($SHEEK_BIN 2>/dev/null)
    
    if [ -n "$selected" ]; then
        # Remove trailing newline if any
        selected=$(echo "$selected" | tr -d '\n')
        BUFFER="$selected"
        CURSOR=${#BUFFER}
    fi
    
    zle redisplay
}

zle -N sheek-widget
bindkey '^T' sheek-widget

