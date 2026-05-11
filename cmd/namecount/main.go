package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Sene4ka/test-task-yadro-tatlin.object.core/internal/counter"
	"github.com/Sene4ka/test-task-yadro-tatlin.object.core/internal/output"
)

func main() {
	sortByFlag := flag.String("sort-by", "freq", "sorting type: 'alph' for alphabetical, 'freq' for frequency")
	orderFlag := flag.String("order", "", "sort order: 'asc' or 'desc' (default 'asc' for alphabetical and 'desc' for frequency)")
	preserveCaseFlag := flag.Bool("preserve-case", false, "preserve original case when sorting")
	outputFileFlag := flag.String("o", "", "output file (default stdout)")

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <filename>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	filename := flag.Arg(0)

	var w io.Writer = os.Stdout
	if *outputFileFlag != "" {
		f, err := os.Create(*outputFileFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		w = f
	}

	sortBy := strings.ToLower(*sortByFlag)
	orderVal := strings.ToLower(*orderFlag)

	if orderVal != "" && orderVal != "asc" && orderVal != "desc" {
		fmt.Fprintf(os.Stderr, "Invalid order: %s (expected 'asc' or 'desc')\n", *orderFlag)
		os.Exit(1)
	}

	var order output.SortOrder
	switch sortBy {
	case "alph":
		if orderVal == "desc" {
			order = output.Desc
		} else {
			order = output.Asc
		}
	case "freq":
		if orderVal == "asc" {
			order = output.Asc
		} else {
			order = output.Desc
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown sorting type: %s\n", *sortByFlag)
		os.Exit(1)
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	counts, err := counter.CountNames(file, *preserveCaseFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while counting: %v\n", err)
		os.Exit(1)
	}

	var writers = map[string]func(io.Writer, map[string]int, output.SortOrder) error{
		"alph": output.WriteMapAlphabetical,
		"freq": output.WriteMapByFrequency,
	}

	if err := writers[sortBy](w, counts, order); err != nil {
		fmt.Fprintf(os.Stderr, "Output error: %v\n", err)
		os.Exit(1)
	}
}
