package main

import (
	"appli.ng/cv/generator/utils"
	"fmt"
	"os"
	"slices"
	"strings"
)

func createBlogPages() {
	initBlogPosts()
	b := strings.Builder{}
	b.WriteString(frontmatterFor("Blog"))
	for _, post := range slices.Backward(blogPosts) {
		b.WriteString(fmt.Sprintf("- %s\n", post.Link())) // TODO: DATES
	}
	WriteFile("Blog.md", b.String())
	// create all blog post pages
	for _, post := range blogPosts {
		WriteFile(fmt.Sprintf(`blog/%s.md`, withoutSpaces(post.Title)), string(post.Bytes()))
	}
}

type BlogPost struct {
	Title        string
	CreationDate dayMonthYr  // TODO: ok? use
	ModifiedDate *dayMonthYr // TODO: ok? use
	Content      []byte
	Tags         []string
}

func (post BlogPost) Bytes() []byte {
	b := strings.Builder{}
	b.WriteString("---\n")
	b.WriteString("title: " + post.Title + "\n")
	b.WriteString("draft: false\n")
	if len(post.Tags) > 0 {
		b.WriteString("tags:\n")
		for _, tag := range post.Tags {
			b.WriteString("  - " + tag + "\n")
		}
	}
	b.WriteString("---\n")
	// TODO: DATES!!!
	return []byte(b.String())
}

func (post BlogPost) Link() string {
	return fmt.Sprintf(`[%s](blog/%s)`, post.Title, withoutSpaces(post.Title))
}

const blogPostDir = "blogPosts/"

var blogPosts []BlogPost
var blogPostNames = utils.Set[string]{}

func newBlogPost(title string, creationDate dayMonthYr, modifiedDate *dayMonthYr, contentFileName string, tags ...string) {
	contentBytes, err := os.ReadFile(blogPostDir + contentFileName)
	if err != nil {
		panic("failed to read blog post " + contentFileName + ": " + err.Error())
	}
	if blogPostNames.Contains(title) {
		panic("blog post " + contentFileName + " already exists")
	}
	blogPosts = append(blogPosts, BlogPost{
		Title:        title,
		CreationDate: creationDate,
		ModifiedDate: modifiedDate,
		Content:      contentBytes,
		Tags:         tags,
	})
	blogPostNames.Add(title)
}
func initBlogPosts() {
	newBlogPost("Example Blog Post", dayMonthYr{3, 12, 2026}, nil, "exampleBlogPost.md")
}
