package main

var cloudServices = map[string]*CloudServicePage{}

type CloudServicePage struct {
	Name     string
	provider string
	// TODO: ADD PROVIDER????
	*tracked
}

func (pg *CloudServicePage) Link() string {
	if pg == nil {
		return "NO_LINK"
	}
	return linkFor(pg.Name, "cv", "service", withoutSpaces(pg.Name))
}

func NewService(name string, provider string) *CloudServicePage {
	out := &CloudServicePage{
		Name:     name,
		provider: provider,
		tracked:  newTracked(),
	}
	cloudServices[name] = out
	return out
}
