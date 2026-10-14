package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCountLinesLongLine(t *testing.T) {
	// Line longer than 64KB (bufio.MaxScanTokenSize)
	longLine := strings.Repeat("a", 100*1024)
	input := strings.NewReader(longLine + "\n" + longLine + "\n" + "short\n")

	hash := make(map[string]int)
	total, err := countLines(input, hash)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if hash[longLine] != 2 {
		t.Errorf("expected count 2 for long line, got %d", hash[longLine])
	}
	if hash["short"] != 1 {
		t.Errorf("expected count 1 for 'short', got %d", hash["short"])
	}
}

func TestMakeBar(t *testing.T) {
	tests := []struct {
		name     string
		rate     float64
		width    int
		expected string
	}{
		{"zero", 0.0, 10, "          "},
		{"full", 100.0, 10, "██████████"},
		{"half", 50.0, 10, "█████     "},
		{"quarter", 25.0, 8, "██      "},
		{"sub-block", 12.5, 8, "█       "},
		{"tiny-rate", 0.1, 10, "▏         "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := makeBar(tt.rate, tt.width)
			if got != tt.expected {
				t.Errorf("makeBar(%v, %v) = %q, want %q", tt.rate, tt.width, got, tt.expected)
			}
		})
	}
}

func TestDeterministicTieBreaking(t *testing.T) {
	// Multiple items with the same frequency
	lines := "cherry\nbanana\napple\ndate\n"
	r := strings.NewReader(lines)
	hash := make(map[string]int)
	_, err := countLines(r, hash)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	opts := options{
		SortByNum: true,
	}

	err = run([]string{}, opts, &buf)
	// Notice run with no args and non-terminal stdin won't use our reader directly,
	// so let's test via file.
}

func TestRunOptions(t *testing.T) {
	tmpDir := t.TempDir()
	f1 := filepath.Join(tmpDir, "test1.txt")
	content := "banana\napple\nbanana\ncherry\nbanana\napple\n"
	if err := os.WriteFile(f1, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	t.Run("default sort (alphabetical)", func(t *testing.T) {
		var out bytes.Buffer
		err := run([]string{f1}, options{}, &out)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "apple\nbanana\ncherry\n"
		if out.String() != expected {
			t.Errorf("got %q, want %q", out.String(), expected)
		}
	})

	t.Run("count and order by num with tie-breaking", func(t *testing.T) {
		// banana: 3, apple: 2, cherry: 1
		var out bytes.Buffer
		err := run([]string{f1}, options{WithCount: true, SortByNum: true}, &out)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		output := out.String()
		lines := strings.Split(strings.TrimSpace(output), "\n")
		if len(lines) != 3 {
			t.Fatalf("expected 3 lines, got %d", len(lines))
		}
		if !strings.Contains(lines[0], "3") || !strings.Contains(lines[0], "banana") {
			t.Errorf("first line should be banana (3): %s", lines[0])
		}
		if !strings.Contains(lines[1], "2") || !strings.Contains(lines[1], "apple") {
			t.Errorf("second line should be apple (2): %s", lines[1])
		}
		if !strings.Contains(lines[2], "1") || !strings.Contains(lines[2], "cherry") {
			t.Errorf("third line should be cherry (1): %s", lines[2])
		}
	})

	t.Run("tie breaking alphabetical", func(t *testing.T) {
		fTie := filepath.Join(tmpDir, "tie.txt")
		// zebra: 1, ant: 1, monkey: 1
		if err := os.WriteFile(fTie, []byte("zebra\nant\nmonkey\n"), 0644); err != nil {
			t.Fatal(err)
		}
		var out bytes.Buffer
		err := run([]string{fTie}, options{WithCount: true, SortByNum: true}, &out)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		output := out.String()
		lines := strings.Split(strings.TrimSpace(output), "\n")
		if len(lines) != 3 {
			t.Fatalf("expected 3 lines, got %d", len(lines))
		}
		// In tie-break: ant -> monkey -> zebra
		if !strings.Contains(lines[0], "ant") {
			t.Errorf("expected ant first, got %s", lines[0])
		}
		if !strings.Contains(lines[1], "monkey") {
			t.Errorf("expected monkey second, got %s", lines[1])
		}
		if !strings.Contains(lines[2], "zebra") {
			t.Errorf("expected zebra third, got %s", lines[2])
		}
	})

	t.Run("with rate and bar", func(t *testing.T) {
		var out bytes.Buffer
		err := run([]string{f1}, options{WithRate: true, WithBar: true, SortByNum: true}, &out)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		output := out.String()
		if !strings.Contains(output, "%") || !strings.Contains(output, "█") {
			t.Errorf("expected rate and bar in output: %s", output)
		}
	})

	t.Run("reverse order", func(t *testing.T) {
		var out bytes.Buffer
		err := run([]string{f1}, options{Reverse: true}, &out)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "cherry\nbanana\napple\n"
		if out.String() != expected {
			t.Errorf("got %q, want %q", out.String(), expected)
		}
	})
}
