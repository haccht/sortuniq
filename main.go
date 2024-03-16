package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"sort"
	"text/tabwriter"

	flags "github.com/jessevdk/go-flags"
	"golang.org/x/crypto/ssh/terminal"
)

type options struct {
	WithCount bool `short:"c" long:"count" description:"Prefix the number of occurrences"`
	WithRate  bool `short:"R" long:"rate" description:"Prefix the percentage of occurrences"`
	SortByNum bool `short:"n" long:"order" description:"Sort by the number of occurrences"`
	Reverse   bool `short:"r" long:"reverse" description:"Reverse the order"`
}

func run() error {
	var opts options
	args, err := flags.Parse(&opts)
	if err != nil {
		if fe, ok := err.(*flags.Error); ok && fe.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	readers := make([]io.Reader, 0, len(args)+1)
	for _, arg := range args {
		f, err := os.Open(arg)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %s", arg, err)
		}
		defer f.Close()
		readers = append(readers, f)
	}
	if !terminal.IsTerminal(0) {
		readers = append(readers, os.Stdin)
	}
	if len(readers) == 0 {
		os.Exit(0)
	}

	total := 0
	hash := make(map[string]int, 1024)
	scanner := bufio.NewScanner(io.MultiReader(readers...))
	for scanner.Scan() {
		total++
		hash[scanner.Text()]++
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	keys := make([]string, 0, len(hash))
	for k := range hash {
		keys = append(keys, k)
	}

	if opts.SortByNum {
		sort.Slice(keys, func(i, j int) bool {
			return hash[keys[i]] > hash[keys[j]]
		})
	} else {
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})
	}
	if opts.Reverse {
		slices.Reverse(keys)
	}

	w := tabwriter.NewWriter(os.Stdout, 4, 0, 1, ' ', tabwriter.AlignRight)
	for _, k := range keys {
		if opts.WithCount {
			fmt.Fprintf(w, "%d\t  ", hash[k])
		}
		if opts.WithRate {
			rate := float32(hash[k]) / float32(total) * 100
			fmt.Fprintf(w, "%6.2f%%\t  ", rate)
		}
		fmt.Fprintf(w, "%s\n", k)
	}
	w.Flush()

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
