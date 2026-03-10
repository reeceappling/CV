package main

import (
	"appli.ng/cv/generator/utils"
	"fmt"
	"maps"
	"slices"
	"strings"
)

var clients = map[string]*Client{}
var clientsOrder = []string{}

type Client struct {
	Name        string
	Info        *string // TODO: NEW! use this when the company can't be outright named!
	Projects    []*Project
	projectsSet utils.Set[string]
	//Languages []string // Calculated later
	// Parent
	company *CompanyPage
}

func (pg *Client) NameValue() string {
	return pg.Name
}

func (pg *Client) ProjectTypeInfo() projectTypeInfo {
	return professionalProjectTypeInfo{client: pg}
}

func (pg *Client) Link() string {
	if pg == nil {
		return "NO_LINK"
	}
	return fmt.Sprintf("[%s](client/%s)", pg.Name, withoutSpaces(pg.Name))
}

func NewClient(name string, info ...string) *Client {
	out := &Client{Name: name, Projects: []*Project{}, Info: nil, projectsSet: utils.Set[string]{}}
	if _, exists := clients[name]; exists {
		panic("client already exists")
	}
	if len(info) > 0 {
		out.Info = &info[0]
	}
	clients[name] = out
	clientsOrder = append(clientsOrder, name)
	return out
}

func (pg *Client) Bytes() []byte {
	builder := strings.Builder{}
	if pg.Projects != nil && len(pg.Projects) > 0 {
		builder.WriteString("# Projects\n")
		for _, proj := range pg.Projects {
			builder.WriteString(fmt.Sprintf("- %s\n", proj.Link()))
		}
	}
	ls := languagesFor(pg.Projects)
	if len(ls) > 0 {
		builder.WriteString("# Languages\n")
		builder.WriteString("Language | Usage Frequency[[" + withoutSpaces(pg.Name) + "#^usageFrequencies\\|*]]\n")
		// TODO: remove frequency???
		builder.WriteString(":-- | :--\n")
		freqs := make(map[Frequency][]string, 6)
		for name, freq := range ls {
			if freqs[freq] == nil {
				freqs[freq] = []string{name}
			} else {
				freqs[freq] = append(freqs[freq], name)
			}
		}
		for f := 5; f >= 0; f-- {
			for _, name := range freqs[Frequency(f)] {
				builder.WriteString(fmt.Sprintf("%s | %s\n", langs[name].Link(), Frequency(f).String()))
			}
		}
	}
	// ALL LOWEST
	builder.WriteString(bytesForAll(pg, true)) // TODO: languages will exist twice???
	builder.WriteString(definitionsArea())
	builder.WriteString(usageFrequencyDefinitions())
	return []byte(builder.String())
}

func projectsFor(cs []*Client) []*Project {
	temp := map[string]*Project{}
	for _, cli := range cs {
		for _, pr := range cli.Projects {
			temp[pr.Name] = pr
		}
	}
	return slices.Collect(maps.Values(temp))
}

func (c *Client) WithProjects(projs ...*Project) *Client {
	for _, proj := range projs {
		if !c.projectsSet.Contains(proj.Name) {
			c.Projects = append(c.Projects, proj)
		}
		proj.TypeInfo = proj.TypeInfo.setClient(c) // TODO: ?????
		for lang, _ := range proj.Languages {
			langs[lang].AddClient(c)
		}
		for l := range maps.Keys(proj.Dbs) {
			dbs[l].AddClient(c)
		}
		for l := range maps.Keys(proj.Caches) {
			caches[l].AddClient(c)
		}
		for l := range maps.Keys(proj.Technologies) {
			techs[l].AddClient(c)
		}
		for l := range maps.Keys(proj.CloudProviders) {
			providers[l].AddClient(c) // TODO: cloud services?
		}
		for l := range maps.Keys(proj.Platforms) {
			platforms[l].AddClient(c)
		}
	}

	return c
}

func (pg *Client) GetAllLowest() (outCaches map[string]*CachePage, outDbs map[string]*DbPage, outLanguages map[string]Frequency, outPlatforms map[string]*PlatformPage, outProviders map[string]*CloudProviderPage, outServices map[string]map[string]*CloudServicePage, outTechnologies map[string]*TechologyPage) {
	outCaches, outDbs, outLanguages, outPlatforms, outProviders, outServices, outTechnologies = initLowest()
	for _, proj := range pg.Projects {
		tempC, tempD, tempL, tempP, tempPr, tempS, tempT := proj.GetAllLowest()
		for key, val := range tempC {
			outCaches[key] = val
		}
		for key, val := range tempD {
			outDbs[key] = val
		}
		for key, val := range tempL {
			if freq, exists := outLanguages[key]; exists {
				// only replace if freq is better!
				if val > freq {
					outLanguages[key] = val
				}
			}
		}
		for key, val := range tempP {
			outPlatforms[key] = val
		}
		for key, val := range tempPr {
			outProviders[key] = val
		}
		for key, val := range tempS {
			outServices[key] = val
		}
		for key, val := range tempT {
			outTechnologies[key] = val
		}
	}
	return
}

var (
	jdClient           = NewClient("Undisclosed Fortune 100 Agribusiness Company", "Fortune 100 Agricultural Business (Think: Green Tractors)")
	sourceAlliesClient = NewClient("Source Allies", "Source Allies internal projects")
	critColaClient     = NewClient("CritCola")
	wellAwareClient    = NewClient("Well Aware NC")
	clarkClient        = NewClient("Chapel Hill Masters Student in Public Health")
	charityClient      = NewClient("Undisclosed Charity")
	teiClient          = NewClient("TEI")
	taeClient          = NewClient("Talley Associates of Engineering") // TODO: JS photo parser
	mafcClient         = NewClient("Monroe Aquatics and Fitness Center")
)

func initClientsAfterProjectsComplete() {
	jdClient = jdClient.WithProjects(polygonBuilderProject, ogreProject, renderProject, statsProject, billingProject, explorerProject, wqdbProject, supportProject, scudsProject, ufoProject, goweProject, tileGenProject, GhaRunnersProject)
	sourceAlliesClient = sourceAlliesClient.WithProjects(saiCollegeProject) // TODO: SIMPSON COLLEGE // TODO: USE!
	critColaClient = critColaClient.WithProjects(CritColaProject)
	wellAwareClient = wellAwareClient.WithProjects(WellAwareProject)
	clarkClient = clarkClient.WithProjects(MastersDataAnalysisProject)
	charityClient = charityClient.WithProjects(CharityProject)
	//teiClient = NewClient("TEI")
	//taeClient = NewClient("Talley Associates of Engineering") // TODO: JS photo parser
	//mafcClient = NewClient("Monroe Aquatics and Fitness Center")
}

var ()
