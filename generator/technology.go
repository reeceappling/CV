package main

var techs = map[string]*TechologyPage{}

type TechologyPage struct { // React, Github actions, etc
	Name string
	*tracked
}

func (pg *TechologyPage) Link() string {
	if pg == nil {
		return "NO_LINK"
	}
	return linkFor(pg.Name, "cv", "technology", withoutSpaces(pg.Name))
}

func NewTechnology(name string) *TechologyPage {
	out := &TechologyPage{
		Name:    name,
		tracked: newTracked(),
	}
	techs[name] = out
	return out
}

func initTechnologies() {
	for _, techName := range []string{
		"Gitlab CI",
		"Drone CI",
	} {
		NewTechnology(techName)
	}
}

// TODO: list all backlinks????
