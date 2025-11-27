# sheek - Ctrl+T keybinding for zsh
# Add this to your ~/.zshrc: source /path/to/sheek/script/keybinds.zsh

SHEEK_BIN=""
_SHEEK_BIN_WARNED=0
: "${SHEEK_AUTO_RECORD:=1}"

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

_sheek_record_history_entry() {
    if [ "${SHEEK_AUTO_RECORD:-1}" -eq 0 ]; then
        return 0
    fi

    emulate -L zsh

    local entry="${1:-}"
    if [ -z "$entry" ]; then
        return 0
    fi

    case "$entry" in
        sheek|sheek\ *|*SheekRecorderSkip* )
            return 0
            ;;
    esac

    if ! _sheek_require_binary; then
        return 0
    fi

    {
        "$SHEEK_BIN" record \
            --cmd "$entry" \
            --cwd "$PWD" \
            --ts "$(date +%s)" >/dev/null 2>&1
    } &!

    return 0
}

sheek-widget() {
    if ! _sheek_require_binary; then
        return 1
    fi

    local selected
    local cursor_pos=${CURSOR:-0}
    local prompt_fragment=""
    if (( cursor_pos > 0 )); then
        prompt_fragment="${BUFFER[1,cursor_pos]}"
    fi
    # Preserve TERM and COLORTERM for color support and pass current prompt fragment
    # Capture stdout for command, but allow stderr to show (for command preview)
    selected=$(SHEEK_INITIAL_QUERY="$prompt_fragment" TERM="${TERM:-xterm-256color}" COLORTERM="${COLORTERM:-truecolor}" "$SHEEK_BIN")
    
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
    if ! _sheek_require_binary; then
        return 1
    fi

    local selected
    # Preserve TERM and COLORTERM for color support
    selected=$(TERM="${TERM:-xterm-256color}" COLORTERM="${COLORTERM:-truecolor}" "$SHEEK_BIN")
    
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

# Register recorder hook once
if [[ -z "${zshaddhistory_functions[(r)_sheek_record_history_entry]}" ]]; then
    zshaddhistory_functions+=(_sheek_record_history_entry)
fi

