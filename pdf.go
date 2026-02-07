package main

import (
	"log"
	"strings"

	"github.com/ledongthuc/pdf"
)

func ExtractPDFText(path string) []string {
	f, r, err := pdf.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var lines []string
	total := r.NumPage()
	for i := 1; i <= total; i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}
		content, _ := p.GetPlainText(nil)
		for _, l := range splitLines(content) {
			lines = append(lines, l)
		}
	}
	return lines
}

func splitLines(s string) []string {
	var out []string
	for _, l := range strings.Split(s, "\n") {
		l = strings.TrimSpace(l)
		if l != "" {
			out = append(out, l)
		}
	}
	return out
}
