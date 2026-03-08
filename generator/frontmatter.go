package main

import "strings"

// TODO: add full text search
// TODO: add graph view
// TODO: LAYOUT

// TODO: USE THIS!!!!!!
func frontmatterFor(title string, tags []string) string {
	b := strings.Builder{}
	b.WriteString("---\n")
	b.WriteString("title: " + title + "\n")
	b.WriteString("draft: false\n")
	if len(tags) > 0 {
		b.WriteString("tags:\n")
		for _, tag := range tags {
			b.WriteString("  - " + tag + "\n")
		}
	}
	b.WriteString("---\n")
	return b.String()
}
