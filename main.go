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
	Count   bool `short:"c" long:"count" description:"Prefix the number of occurrences"`
	Sort    bool `short:"n" long:"order" description:"Sort by the number of occurrences"`
	Reverse bool `short:"r" long:"reverse" description:"Reverse the order"`
}

func run() error {
	var opts options
	args, err := flags.Parse(&opts)
	if err != nil {
		if fe, ok := err.(*flags.Error); ok && fe.Type == flags.ErrHelp {
			os.Exit(0)
		}
		return err
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

	hash := make(map[string]int, 1024)
	scanner := bufio.NewScanner(io.MultiReader(readers...))
	for scanner.Scan() {
		hash[scanner.Text()]++
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	keys := make([]string, 0, len(hash))
	for k, _ := range hash {
		keys = append(keys, k)
	}

	if opts.Sort {
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

	if opts.Count {
		w := tabwriter.NewWriter(os.Stdout, 4, 0, 1, ' ', tabwriter.AlignRight)
		for _, k := range keys {
			fmt.Fprintf(w, "%d\t  %s\n", hash[k], k)
		}
		w.Flush()
	} else {
		for _, k := range keys {
			fmt.Println(k)
		}
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
