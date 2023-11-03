package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	flags "github.com/jessevdk/go-flags"
	"golang.org/x/crypto/ssh/terminal"
)

type options struct {
	Count   bool `short:"c" long:"count" description:"Prefix lines by the number of occurrences"`
	Reverse bool `short:"r" long:"reverse" description:"Reverse the result"`
}

func main() {
	var opts options
	args, err := flags.Parse(&opts)
	if err != nil {
		if fe, ok := err.(*flags.Error); ok && fe.Type == flags.ErrHelp {
			os.Exit(0)
		}
		log.Fatal(err)
	}

	readers := make([]io.Reader, 0, len(args)+1)
	if !terminal.IsTerminal(0) {
		readers = append(readers, os.Stdin)
	}
	for _, arg := range args {
		f, err := os.Open(arg)
		if err != nil {
			log.Fatalf("fail to open file %s: %s", arg, err)
		}
		defer f.Close()
		readers = append(readers, f)
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
		log.Fatalln(scanner.Err())
	}

	keys := make([]string, 0, len(hash))
	for k, _ := range hash {
		keys = append(keys, k)
	}

	if opts.Count {
		sort.Slice(keys, func(i, j int) bool {
			if opts.Reverse {
				return hash[keys[i]] < hash[keys[j]]
			} else {
				return hash[keys[i]] > hash[keys[j]]
			}
		})

		for _, k := range keys {
			fmt.Printf("%8d\t%s\n", hash[k], k)
		}
	} else {
		sort.Slice(keys, func(i, j int) bool {
			if opts.Reverse {
				return keys[i] > keys[j]
			} else {
				return keys[i] < keys[j]
			}
		})

		for _, k := range keys {
			fmt.Printf("%s\n", k)
		}
	}
}
