package output

import (
	"bytes"
	"testing"
)

func TestWriteMapAlphabetical(t *testing.T) {
	counts := map[string]int{"Миша": 1, "Алёна": 2, "Дима": 1}
	var buf bytes.Buffer
	if err := WriteMapAlphabetical(&buf, counts, Asc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	expected := "Алёна:2\nДима:1\nМиша:1\n"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestWriteMapAlphabeticalDesc(t *testing.T) {
	counts := map[string]int{"Миша": 1, "Алёна": 2, "Дима": 1}
	var buf bytes.Buffer
	if err := WriteMapAlphabetical(&buf, counts, Desc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	expected := "Миша:1\nДима:1\nАлёна:2\n"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestWriteMapByFrequency(t *testing.T) {
	counts := map[string]int{"Миша": 1, "Алёна": 2, "Дима": 1}
	var buf bytes.Buffer
	if err := WriteMapByFrequency(&buf, counts, Desc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	expected := "Алёна:2\nДима:1\nМиша:1\n"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestWriteMapByFrequencyAsc(t *testing.T) {
	counts := map[string]int{"Миша": 1, "Алёна": 2, "Дима": 1}
	var buf bytes.Buffer
	if err := WriteMapByFrequency(&buf, counts, Asc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	expected := "Дима:1\nМиша:1\nАлёна:2\n"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestWriteEmptyMap(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteMapAlphabetical(&buf, map[string]int{}, Asc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got %q", buf.String())
	}
}
