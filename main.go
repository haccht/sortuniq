package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

func main() {
	var (
		count   = flag.Bool("c", false, "prefix lines by the number of occurrences")
		reverse = flag.Bool("r", false, "reverse the result")
	)
	flag.Parse()

	readers := make([]io.Reader, 0, flag.NArg()+1)
	readers = append(readers, os.Stdin)
	for _, arg := range flag.Args() {
		f, err := os.Open(arg)
		if err != nil {
			log.Fatalf("fail to open file %s: %s", arg, err)
		}
		defer f.Close()

		readers = append(readers, f)
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

	if *count {
		sort.Slice(keys, func(i, j int) bool {
			if *reverse {
				return hash[keys[i]] < hash[keys[j]]
			} else {
				return hash[keys[i]] > hash[keys[j]]
			}
		})

		for _, k := range keys {
			fmt.Printf("%10d\t%s\n", hash[k], k)
		}
	} else {
		sort.Slice(keys, func(i, j int) bool {
			if *reverse {
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
