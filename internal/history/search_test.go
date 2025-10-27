package history

import "testing"

func TestSearchExact(t *testing.T) {
	cmds, _ := LoadAndParseZshHistory()

	tests := []struct {
		input string
		want  int
	}{
		{"kitty", 7},
		{"repeat", 2},
	}

	for _, test := range tests {
		got := SearchExact(cmds, test.input)
		if len(got) != test.want {
			t.Errorf("SearchExact(%q) = %d, want %d", test.input, len(got), test.want)
		}
	}
}
