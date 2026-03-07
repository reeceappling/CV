package main

import (
	"fmt"
)

var platforms = map[string]*PlatformPage{}

type PlatformPage struct { // Datadog, Github, etc
	Name string
	*tracked
}

func (pg *PlatformPage) Link() string {
	if pg == nil {
		return "NO_LINK"
	}
	return fmt.Sprintf("[%s](platform/%s)", pg.Name, withoutSpaces(pg.Name))
}

func NewPlatform(name string) *PlatformPage {
	out := &PlatformPage{
		Name:    name,
		tracked: newTracked(),
	}
	platforms[name] = out
	return out
}
