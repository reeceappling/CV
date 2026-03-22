package main

import "strings"

var caches = map[string]*CachePage{}

type CachePage struct {
	Name string
	*tracked
}

func Link(linkable Linkable) string {
	return linkFor(linkable.Title(), linkable.Dst())
}

func (pg *CachePage) Dst() string {
	if pg == nil {
		return "NO_LINK"
	}
	return dstFor("cv", "cache", withoutSpaces(pg.Name))
}

func (pg *CachePage) Title() string {
	if pg == nil {
		return "NO_TITLE"
	}
	return pg.Name
}
func dstFor(elems ...string) string {
	return strings.Join(elems, "/")
}
func (pg *CachePage) EntryType() string {
	return "Cache"
}

func NewCache(name string) *CachePage {
	out := &CachePage{
		Name:    name,
		tracked: newTracked(),
	}
	caches[name] = out
	addLinkable(strings.ToLower(name), out)
	return out
}
