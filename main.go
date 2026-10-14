package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"strings"
	"text/tabwriter"

	flags "github.com/jessevdk/go-flags"
	"golang.org/x/term"
)

type options struct {
	WithCount bool `short:"c" long:"count" description:"Prefix the number of occurrences"`
	WithRate  bool `short:"R" long:"rate" description:"Prefix the percentage of occurrences"`
	WithBar   bool `short:"b" long:"bar" description:"Prefix a bar chart showing the percentage"`
	SortByNum bool `short:"n" long:"order" description:"Sort by the number of occurrences"`
	Reverse   bool `short:"r" long:"reverse" description:"Reverse the order"`
}

func countLines(r io.Reader, hash map[string]int) (int, error) {
	reader := bufio.NewReader(r)
	total := 0
	for {
		line, err := reader.ReadBytes('\n')
		if len(line) > 0 {
			if line[len(line)-1] == '\n' {
				line = line[:len(line)-1]
				if len(line) > 0 && line[len(line)-1] == '\r' {
					line = line[:len(line)-1]
				}
			}
			total++
			hash[string(line)]++
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return total, err
		}
	}
	return total, nil
}

func makeBar(rate float64, width int) string {
	if rate <= 0 {
		return strings.Repeat(" ", width)
	}
	if rate > 100 {
		rate = 100
	}
	fractionChars := []rune{' ', '▏', '▎', '▍', '▌', '▋', '▊', '▉'}

	totalEighths := int(math.Round(rate / 100.0 * float64(width*8)))
	if totalEighths == 0 && rate > 0 {
		totalEighths = 1
	}

	fullBlocks := totalEighths / 8
	remainder := totalEighths % 8

	var sb strings.Builder
	for i := 0; i < fullBlocks && i < width; i++ {
		sb.WriteRune('█')
	}
	curLen := fullBlocks
	if curLen < width && remainder > 0 {
		sb.WriteRune(fractionChars[remainder])
		curLen++
	}
	for curLen < width {
		sb.WriteRune(' ')
		curLen++
	}
	return sb.String()
}

func run(args []string, opts options, stdout io.Writer) error {
	hash := make(map[string]int, 1024)
	total := 0

	if len(args) == 0 {
		if term.IsTerminal(int(os.Stdin.Fd())) {
			return nil
		}
		n, err := countLines(os.Stdin, hash)
		if err != nil {
			return err
		}
		total += n
	} else {
		for _, arg := range args {
			if arg == "-" {
				n, err := countLines(os.Stdin, hash)
				if err != nil {
					return err
				}
				total += n
			} else {
				if err := func() error {
					f, err := os.Open(arg)
					if err != nil {
						return fmt.Errorf("failed to open file %s: %w", arg, err)
					}
					defer f.Close()

					n, err := countLines(f, hash)
					if err != nil {
						return err
					}
					total += n
					return nil
				}(); err != nil {
					return err
				}
			}
		}
	}

	if len(hash) == 0 {
		return nil
	}

	keys := make([]string, 0, len(hash))
	for k := range hash {
		keys = append(keys, k)
	}

	if opts.SortByNum {
		slices.SortFunc(keys, func(a, b string) int {
			if hash[a] != hash[b] {
				if hash[a] > hash[b] {
					return -1
				}
				return 1
			}
			if a < b {
				return -1
			} else if a > b {
				return 1
			}
			return 0
		})
	} else {
		slices.Sort(keys)
	}
	if opts.Reverse {
		slices.Reverse(keys)
	}

	w := tabwriter.NewWriter(stdout, 4, 0, 1, ' ', tabwriter.AlignRight)
	for _, k := range keys {
		if opts.WithCount {
			fmt.Fprintf(w, "%d\t  ", hash[k])
		}
		rate := 0.0
		if total > 0 {
			rate = float64(hash[k]) / float64(total) * 100
		}
		if opts.WithRate {
			fmt.Fprintf(w, "%6.2f%%\t  ", rate)
		}
		if opts.WithBar {
			fmt.Fprintf(w, "%s\t  ", makeBar(rate, 20))
		}
		fmt.Fprintf(w, "%s\n", k)
	}
	w.Flush()

	return nil
}

func main() {
	var opts options
	args, err := flags.Parse(&opts)
	if err != nil {
		if fe, ok := err.(*flags.Error); ok && fe.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	if err := run(args, opts, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
