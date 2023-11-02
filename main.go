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

type item struct {
	k string
	v int
}

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

	lines := make(map[string]int, 1024)
	scanner := bufio.NewScanner(io.MultiReader(readers...))

	for scanner.Scan() {
		line := scanner.Text()
		lines[line]++
	}
	if scanner.Err() != nil {
		log.Fatalln(scanner.Err())
	}

	items := make([]item, 0, len(lines))
	for k, v := range lines {
		items = append(items, item{k, v})
	}

	if *count {
		sort.Slice(items, func(i, j int) bool {
			if *reverse {
				return items[i].v < items[j].v
			} else {
				return items[i].v > items[j].v
			}
		})

		for _, item := range items {
			fmt.Printf("%10d\t%s\n", item.v, item.k)
		}
	} else {
		sort.Slice(items, func(i, j int) bool {
			if *reverse {
				return items[i].k > items[j].k
			} else {
				return items[i].k < items[j].k
			}
		})

		for _, item := range items {
			fmt.Printf("%s\n", item.k)
		}
	}
}
