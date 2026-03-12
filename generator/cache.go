package main

var caches = map[string]*CachePage{}

type CachePage struct {
	Name string
	*tracked
}

func (pg *CachePage) Link() string {
	if pg == nil {
		return "NO_LINK"
	}
	return linkFor(pg.Name, "cv", "cache", withoutSpaces(pg.Name))
}

func NewCache(name string) *CachePage {
	out := &CachePage{
		Name:    name,
		tracked: newTracked(),
	}
	caches[name] = out
	return out
}
