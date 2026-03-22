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
		WriteFile(fmt.Sprintf(`note/%s.md`, withoutSpaces(note.TitleText)), string(note.Bytes()))
	}
}

type Note struct {
	TitleText    string
	CreationDate dayMonthYr
	ModifiedDate *dayMonthYr
	ContentFile  string
	Tags         []string
}

func (note *Note) Bytes() []byte {
	b := strings.Builder{}
	b.WriteString("---\n")
	b.WriteString("title: " + note.TitleText + "\n")
	b.WriteString("draft: false\n")
	if len(note.Tags) > 0 {
		b.WriteString("tags:\n")
		for _, tag := range note.Tags {
			b.WriteString("  - " + tag + "\n")
		}
	}
	b.WriteString("---\n")
	contentBytes, err := os.ReadFile(notesDir + note.ContentFile)
	if err != nil {
		panic("failed to read note file " + note.ContentFile + ": " + err.Error())
	}
	b.Write(contentBytes)
	// TODO: DATES!!!
	return []byte(b.String())
}

func (note *Note) Link() string {
	return fmt.Sprintf(`[%s](note/%s)`, note.Title, withoutSpaces(note.TitleText))
}
func (pg *Note) Dst() string {
	return dstFor("note", withoutSpaces(pg.TitleText))
}

func (pg *Note) Title() string {
	if pg == nil {
		return noLinkText
	}
	return pg.TitleText
}

const notesDir = "notes/"

var notes []Note
var noteNames = utils.Set[string]{}

func newNote(title string, creationDate dayMonthYr, modifiedDate *dayMonthYr, contentFileName string, tags ...string) {
	if noteNames.Contains(title) {
		panic("note " + contentFileName + " already exists")
	}
	notes = append(notes, Note{
		TitleText:    title,
		CreationDate: creationDate,
		ModifiedDate: modifiedDate,
		ContentFile:  contentFileName,
		Tags:         tags,
	})
	blogPostNames.Add(title)
}
func initNotes() {
	newNote("Example Note", *NewPostDate(3, 12, 2026), nil, "exampleNote.md")
}
