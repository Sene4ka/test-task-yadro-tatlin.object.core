package output

import (
	"cmp"
	"fmt"
	"io"
	"slices"
)

type SortOrder int

const (
	Asc SortOrder = iota
	Desc
)

type countsEntry struct {
	name  string
	count int
}

func flattenMap(counts map[string]int) []countsEntry {
	result := make([]countsEntry, 0, len(counts))
	for k, v := range counts {
		result = append(result, countsEntry{k, v})
	}
	return result
}

func writeEntries(w io.Writer, entries []countsEntry) error {
	for _, e := range entries {
		if _, err := fmt.Fprintf(w, "%s:%d\n", e.name, e.count); err != nil {
			return err
		}
	}
	return nil
}

func WriteMapAlphabetical(w io.Writer, counts map[string]int, order SortOrder) error {
	entries := flattenMap(counts)
	slices.SortFunc(entries, func(a, b countsEntry) int {
		if order == Desc {
			return -cmp.Compare(a.name, b.name)
		}
		return cmp.Compare(a.name, b.name)
	})
	return writeEntries(w, entries)
}

func WriteMapByFrequency(w io.Writer, counts map[string]int, order SortOrder) error {
	entries := flattenMap(counts)
	slices.SortFunc(entries, func(a, b countsEntry) int {
		if a.count != b.count {
			if order == Asc {
				return cmp.Compare(a.count, b.count)
			}
			return -cmp.Compare(a.count, b.count)
		}
		return cmp.Compare(a.name, b.name)
	})
	return writeEntries(w, entries)
}
