# sheek - Ctrl+T keybinding for bash
# Add this to your ~/.bashrc: source /path/to/sheek/script/keybinds.bash

SHEEK_BIN=""
_SHEEK_BIN_WARNED=0

_sheek_require_binary() {
    if [ -n "$SHEEK_BIN" ] && [ -x "$SHEEK_BIN" ]; then
        return 0
    fi

    local candidate
    for candidate in "$HOME/.local/bin/sheek" "$HOME/go/bin/sheek" "/usr/local/bin/sheek" "/usr/bin/sheek"; do
        if [ -x "$candidate" ]; then
            SHEEK_BIN="$candidate"
            return 0
        fi
    done

    if command -v sheek &> /dev/null; then
        candidate="$(command -v sheek)"
        if [ -n "$candidate" ] && [ -x "$candidate" ]; then
            SHEEK_BIN="$candidate"
            return 0
        fi
    fi

    if [ "${_SHEEK_BIN_WARNED:-0}" -eq 0 ]; then
        echo "sheek: binary not found. Install it with ./script/install.sh and reload your shell." >&2
        _SHEEK_BIN_WARNED=1
    fi

    return 1
}

sheek-bash-widget() {
    if ! _sheek_require_binary; then
        return 1
    fi

    local selected
    local cursor_pos="${READLINE_POINT:-0}"
    local current_line="${READLINE_LINE:-}"
    local prompt_fragment="${current_line:0:$cursor_pos}"
    # Preserve TERM and COLORTERM for color support and pass the current prompt
    # Capture stdout for command, but allow stderr to show (for command preview)
    selected=$(SHEEK_INITIAL_QUERY="$prompt_fragment" TERM="${TERM:-xterm-256color}" COLORTERM="${COLORTERM:-truecolor}" "$SHEEK_BIN")
    
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
    if ! _sheek_require_binary; then
        return 1
    fi

    local selected
    # Preserve TERM and COLORTERM for color support
    selected=$(TERM="${TERM:-xterm-256color}" COLORTERM="${COLORTERM:-truecolor}" "$SHEEK_BIN")
    
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

