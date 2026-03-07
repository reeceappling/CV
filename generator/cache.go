package main

import "fmt"

var caches = map[string]*CachePage{}

type CachePage struct {
	Name string
	*tracked
}

func (pg *CachePage) Link() string {
	if pg == nil {
		return "NO_LINK"
	}
	return fmt.Sprintf("[%s](cache/%s)", pg.Name, withoutSpaces(pg.Name))
}

func NewCache(name string) *CachePage {
	out := &CachePage{
		Name:    name,
		tracked: newTracked(),
	}
	caches[name] = out
	return out
}
