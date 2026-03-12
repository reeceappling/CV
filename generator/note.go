package main

import (
	"appli.ng/cv/generator/utils"
	"fmt"
	"os"
	"slices"
	"strings"
)

func createNotesPages() {
	initNotes()
	b := strings.Builder{}
	b.WriteString(frontmatterFor("Notes"))
	for _, note := range slices.Backward(notes) {
		b.WriteString(fmt.Sprintf("- %s\n", note.Link())) // TODO: DATES
	}
	WriteFile("Notes.md", b.String())
	// create all note pages
	for _, note := range notes {
		WriteFile(fmt.Sprintf(`note/%s.md`, withoutSpaces(note.Title)), string(note.Bytes()))
	}
}

type Note struct {
	Title        string
	CreationDate dayMonthYr  // TODO: ok? use
	ModifiedDate *dayMonthYr // TODO: ok? use
	Content      []byte
	Tags         []string
}

func (note Note) Bytes() []byte {
	b := strings.Builder{}
	b.WriteString("---\n")
	b.WriteString("title: " + note.Title + "\n")
	b.WriteString("draft: false\n")
	if len(note.Tags) > 0 {
		b.WriteString("tags:\n")
		for _, tag := range note.Tags {
			b.WriteString("  - " + tag + "\n")
		}
	}
	b.WriteString("---\n")
	// TODO: DATES!!!
	return []byte(b.String())
}

func (note Note) Link() string {
	return fmt.Sprintf(`[%s](note/%s)`, note.Title, withoutSpaces(note.Title))
}

const notesDir = "notes/"

var notes []Note
var noteNames = utils.Set[string]{}

func newNote(title string, creationDate dayMonthYr, modifiedDate *dayMonthYr, contentFileName string, tags ...string) {
	contentBytes, err := os.ReadFile(notesDir + contentFileName)
	if err != nil {
		panic("failed to read note file " + contentFileName + ": " + err.Error())
	}
	if noteNames.Contains(title) {
		panic("note " + contentFileName + " already exists")
	}
	notes = append(notes, Note{
		Title:        title,
		CreationDate: creationDate,
		ModifiedDate: modifiedDate,
		Content:      contentBytes,
		Tags:         tags,
	})
	blogPostNames.Add(title)
}
func initNotes() {
	newNote("Example Note", dayMonthYr{3, 12, 2026}, nil, "exampleNote.md")
}
