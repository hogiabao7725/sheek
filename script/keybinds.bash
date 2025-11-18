# sheek - Ctrl+T keybinding for bash
# Add this to your ~/.bashrc: source /path/to/sheek/script/keybinds.bash

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

sheek-bash-widget() {
    local selected
    selected=$($SHEEK_BIN 2>/dev/null)
    
    if [ -n "$selected" ]; then
        # Remove trailing newline if any
        selected=$(echo "$selected" | tr -d '\n')
        # Insert the selected command at cursor position
        READLINE_LINE="${READLINE_LINE:0:$READLINE_POINT}$selected${READLINE_LINE:$READLINE_POINT}"
        READLINE_POINT=$(($READLINE_POINT + ${#selected}))
    fi
}

# Bind Ctrl+T to sheek widget
bind -x '"\C-t": sheek-bash-widget'

