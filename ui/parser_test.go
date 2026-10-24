package ui

import (
	"errors"
	"slices"
	"testing"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		line     string
		wantName string
		wantArgs []string
	}{
		{"list", "list", []string{}},
		{"  ADD  scooter\tVK-1 ", "add", []string{"scooter", "VK-1"}},
		{`add scooter VK-1 "Ninebot Max"`, "add", []string{"scooter", "VK-1", "Ninebot Max"}},
		{`add "" x`, "add", []string{"", "x"}},
		{`say "a \"b\" \\c"`, "say", []string{`a "b" \c`}},
		{`path C:\dir`, "path", []string{`C:\dir`}},
	}
	for _, tt := range tests {
		cmd, err := ParseLine(tt.line)
		if err != nil {
			t.Fatalf("ParseLine(%q): %v", tt.line, err)
		}
		if cmd.Name != tt.wantName || !slices.Equal(cmd.Args, tt.wantArgs) {
			t.Errorf("ParseLine(%q) = %q %q, want %q %q",
				tt.line, cmd.Name, cmd.Args, tt.wantName, tt.wantArgs)
		}
	}
}

func TestParseLine_Errors(t *testing.T) {
	if _, err := ParseLine("   "); !errors.Is(err, ErrEmptyLine) {
		t.Errorf("пустая строка: got %v", err)
	}
	if _, err := ParseLine(`add "abc`); !errors.Is(err, ErrUnclosedQuote) {
		t.Errorf("незакрытая кавычка: got %v", err)
	}
}

func TestCommandArgs(t *testing.T) {
	cmd := Command{Name: "x", Args: []string{"5", "0,4", "abc", "NaN"}}

	if err := cmd.RequireArgs(3); !errors.Is(err, ErrArgsCount) {
		t.Errorf("RequireArgs: got %v", err)
	}
	if n, err := cmd.IntArg(0); err != nil || n != 5 {
		t.Errorf("IntArg(0) = %d, %v", n, err)
	}
	if f, err := cmd.FloatArg(1); err != nil || f != 0.4 {
		t.Errorf("FloatArg(1) = %v, %v", f, err)
	}
	if _, err := cmd.IntArg(2); !errors.Is(err, ErrBadArg) {
		t.Errorf("IntArg(2): got %v", err)
	}
	if _, err := cmd.FloatArg(3); !errors.Is(err, ErrBadArg) {
		t.Errorf("FloatArg(NaN): got %v", err)
	}
	if _, err := cmd.Arg(10); !errors.Is(err, ErrArgsCount) {
		t.Errorf("Arg(10): got %v", err)
	}
}
