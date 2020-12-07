package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		log.Fatalf("usage: %s <costs.csv>", os.Args[0])
	}

	filename := args[0]
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := csv.NewWriter(os.Stdout)
	r := csv.NewReader(f)

	var dateMatch = regexp.MustCompile(`\A\d{4}-\d{2}-\d{2}\z`)

	var headers []string
	var names []string
	var wroteHeader bool

	out := make([]string, 0, 4)
	for i := 0; ; i++ {
		row, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		if i == 0 {
			for _, id := range row {
				id = strings.TrimSuffix(id, "($)")
				id = strings.TrimSpace(id)
				headers = append(headers, id)
			}
			continue
		}

		if i == 1 || i == 2 && !dateMatch.MatchString(row[0]) {
			// The second and third rows can have things like
			// - Linked Account Name
			// - Linked Account Total
			// Keep the name as additional useful metadata
			if strings.Index(row[0], "Name") >= 0 {
				for _, name := range row {
					name = strings.TrimSuffix(name, "($)")
					name = strings.TrimSpace(name)
					names = append(names, name)
				}
			}
			continue
		}

		// date,account_id,account_name,amount

		date := row[0]
		if !dateMatch.MatchString(date) {
			log.Fatalf("Row %d Expected date in first column but was: %q", i, date)
		}

		if !wroteHeader {
			if len(names) > 0 {
				w.Write([]string{"date", "id", "name", "amt"})
			} else {
				w.Write([]string{"date", "id", "amt"})
			}
			wroteHeader = true
		}

		for i := 1; i < len(row)-1; i++ {
			out = out[:0]
			out = append(out, date)
			out = append(out, headers[i]) // id
			if len(names) > 0 {
				out = append(out, names[i])
			}
			out = append(out, row[i]) // amt
			err := w.Write(out)
			if err != nil {
				panic(err)
			}
		}
	}
	w.Flush()
}
