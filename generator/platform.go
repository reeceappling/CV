package main

var platforms = map[string]*PlatformPage{}

type PlatformPage struct { // Datadog, Github, etc
	Name string
	*tracked
}

func (pg *PlatformPage) Link() string {
	if pg == nil {
		return "NO_LINK"
	}
	return linkFor(pg.Name, "cv", "platform", withoutSpaces(pg.Name))
}

func NewPlatform(name string) *PlatformPage {
	out := &PlatformPage{
		Name:    name,
		tracked: newTracked(),
	}
	platforms[name] = out
	return out
}
