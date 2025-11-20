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
    echo "Run: cd /path/to/sheek && ./script/install.sh" >&2
    return 1
fi

# Verify binary is executable
if [ ! -x "$SHEEK_BIN" ]; then
    echo "Warning: sheek binary is not executable: $SHEEK_BIN" >&2
    return 1
fi

sheek-bash-widget() {
    local selected
    local cursor_pos="${READLINE_POINT:-0}"
    local current_line="${READLINE_LINE:-}"
    local prompt_fragment="${current_line:0:$cursor_pos}"
    # Preserve TERM and COLORTERM for color support and pass the current prompt
    # Capture stdout for command, but allow stderr to show (for command preview)
    selected=$(SHEEK_INITIAL_QUERY="$prompt_fragment" TERM="${TERM:-xterm-256color}" COLORTERM="${COLORTERM:-truecolor}" $SHEEK_BIN)
    
    if [ -n "$selected" ]; then
        # Remove trailing newline if any
        selected=$(echo "$selected" | tr -d '\n')
        # Insert the selected command at cursor position
        READLINE_LINE="${READLINE_LINE:0:$READLINE_POINT}$selected${READLINE_LINE:$READLINE_POINT}"
        READLINE_POINT=$(($READLINE_POINT + ${#selected}))
    fi
}

# Shell function to run sheek and add selected command to history
# Note: In bash, we can't directly modify the prompt buffer from a function,
# so we add it to history and user can press Up arrow to retrieve it.
# For full functionality (command in prompt buffer), use Ctrl+T keybind instead.
sheek() {
    local selected
    # Preserve TERM and COLORTERM for color support
    selected=$(TERM="${TERM:-xterm-256color}" COLORTERM="${COLORTERM:-truecolor}" $SHEEK_BIN)
    
    if [ -n "$selected" ]; then
        # Remove trailing newline if any
        selected=$(echo "$selected" | tr -d '\n')
        # Add to history so user can press Up arrow to get it
        if [[ $- == *i* ]] && [ -n "$BASH_VERSION" ]; then
            history -s "$selected"
            echo "Command added to history. Press Up arrow to retrieve: $selected"
        else
            # Just print the command
            echo "$selected"
        fi
    fi
}

# Bind Ctrl+T to sheek widget
# Unbind any existing Ctrl+T binding first
bind -r '\C-t' 2>/dev/null || true
# Bind Ctrl+T to sheek widget
bind -x '"\C-t": sheek-bash-widget'

