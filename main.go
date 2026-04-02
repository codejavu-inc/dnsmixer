package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/publicsuffix"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: dnsmixer [options]\n\nOptions:\n")
		fmt.Fprintf(os.Stderr, "  -w <file>   Wordlist file (if omitted, only permutes existing subdomains)\n")
		fmt.Fprintf(os.Stderr, "  -t <file>   Target list of subdomains (Required)\n")
		fmt.Fprintf(os.Stderr, "  -l1         Process only first-level subdomains\n")
		fmt.Fprintf(os.Stderr, "  -sc <sub>   Scope execution to only this subdomain and its children\n")
		fmt.Fprintf(os.Stderr, "  -o <file>   Output file path to save results (default is stdout)\n")
	}

	wordlistPtr := flag.String("w", "", "")
	targetPtr := flag.String("t", "", "")
	l1Ptr := flag.Bool("l1", false, "")
	scPtr := flag.String("sc", "", "")
	outPtr := flag.String("o", "", "")

	flag.Parse()

	if *targetPtr == "" {
		fmt.Fprintf(os.Stderr, "Error: -t is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	targets, err := readLines(*targetPtr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading targets: %v\n", err)
		os.Exit(1)
	}

	var words []string
	if *wordlistPtr != "" {
		words, err = readLines(*wordlistPtr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading wordlist: %v\n", err)
			os.Exit(1)
		}
	}

	var outFile *os.File
	var writer *bufio.Writer
	if *outPtr != "" {
		outFile, err = os.Create(*outPtr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer outFile.Close()
		writer = bufio.NewWriter(outFile)
		defer writer.Flush()
	}

	filtered := filterTargets(targets, *l1Ptr, *scPtr)

	if len(words) > 0 {
		var total int64
		for _, domain := range filtered {
			base, _ := publicsuffix.EffectiveTLDPlusOne(domain)
			sub := strings.TrimSuffix(domain, "."+base)
			if sub == domain || sub == "" {
				total += int64(len(words))
			} else {
				total += int64(len(words)) * 5
			}
		}

		var count int64
		for _, domain := range filtered {
			base, _ := publicsuffix.EffectiveTLDPlusOne(domain)
			sub := strings.TrimSuffix(domain, "."+base)
			if sub == domain {
				sub = ""
			}
			for _, w := range words {
				out(writer, fmt.Sprintf("%s.%s", w, base))
				count++
				if sub != "" {
					out(writer, fmt.Sprintf("%s-%s.%s", w, sub, base))
					out(writer, fmt.Sprintf("%s.%s.%s", w, sub, base))
					out(writer, fmt.Sprintf("%s-%s.%s", sub, w, base))
					out(writer, fmt.Sprintf("%s.%s.%s", sub, w, base))
					count += 4
				}
				fmt.Fprintf(os.Stderr, "\rLeft to generate: %d", total-count)
			}
		}
	} else {
		groups := group(filtered)
		var total int64
		for _, subs := range groups {
			n := int64(len(subs))
			total += n * (n - 1) * 2
		}
		var count int64
		for base, subs := range groups {
			for i := 0; i < len(subs); i++ {
				for j := 0; j < len(subs); j++ {
					if i == j {
						continue
					}
					out(writer, fmt.Sprintf("%s-%s.%s", subs[i], subs[j], base))
					out(writer, fmt.Sprintf("%s.%s.%s", subs[i], subs[j], base))
					count += 2
					fmt.Fprintf(os.Stderr, "\rLeft to generate: %d", total-count)
				}
			}
		}
	}
	fmt.Fprintln(os.Stderr, "\nDone.")
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text != "" {
			lines = append(lines, text)
		}
	}
	return lines, scanner.Err()
}

func filterTargets(targets []string, l1 bool, sc string) []string {
	var filtered []string
	for _, domain := range targets {
		domain = strings.TrimSpace(domain)
		if domain == "" {
			continue
		}
		base, err := publicsuffix.EffectiveTLDPlusOne(domain)
		if err != nil {
			continue
		}
		if sc != "" {
			if domain != sc && !strings.HasSuffix(domain, "."+sc) {
				continue
			}
		}
		if l1 {
			sub := strings.TrimSuffix(domain, "."+base)
			if sub == domain {
				sub = ""
			}
			if strings.Count(sub, ".") > 0 {
				continue
			}
		}
		filtered = append(filtered, domain)
	}
	return filtered
}

func group(targets []string) map[string][]string {
	m := make(map[string][]string)
	for _, domain := range targets {
		base, _ := publicsuffix.EffectiveTLDPlusOne(domain)
		sub := strings.TrimSuffix(domain, "."+base)
		if sub == domain {
			sub = ""
		}
		if sub != "" {
			m[base] = append(m[base], sub)
		}
	}
	return m
}

func out(w *bufio.Writer, s string) {
	if w != nil {
		w.WriteString(s + "\n")
	} else {
		fmt.Println(s)
	}
}
