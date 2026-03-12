package main

import (
	"appli.ng/cv/generator/utils"
	"fmt"
	"maps"
	"slices"
	"strings"
)

var companies = map[string]*CompanyPage{}
var companiesOrder = []string{}

type dayMonthYr struct {
	Month, Day, Year int
}

func NewPostDate(month, day, year int) *dayMonthYr {
	return &dayMonthYr{month, day, year}
}

type monthYr struct {
	Month int
	Year  int
}

func (m *monthYr) String() string {
	if m == nil {
		return "current"
	}
	return fmt.Sprintf("%d-%d", m.Month, m.Year)
}

type CompanyPage struct {
	Name      string
	Positions []*Position
	Start     monthYr
	End       *monthYr // None == current
	//Clients   []*Client // resolved later
	//Projects  []string // Resolved from clients
	//Languages []string // Calculated later
}

func (pg *CompanyPage) NameValue() string {
	return pg.Name
}

func (pg *CompanyPage) Link() string {
	if pg == nil {
		return "NO_LINK"
	}
	return linkFor(pg.Name, "cv", "company", withoutSpaces(pg.Name))
}

func NewCompany(name string, startMo, startYr int, endMo, endYr *int) *CompanyPage {
	if _, exists := companies[name]; exists {
		panic("company already exists")
	}
	out := &CompanyPage{
		Name: name,
		Start: monthYr{
			Month: startMo,
			Year:  startYr,
		},
		End:       nil,
		Positions: []*Position{},
	}
	if endMo != nil && endYr != nil {
		out.End = &monthYr{
			Month: *endMo,
			Year:  *endYr,
		}
	}
	companies[name] = out
	companiesOrder = append(companiesOrder, name)
	return out
}

func (pg *CompanyPage) Bytes() []byte {
	builder := strings.Builder{}
	builder.WriteString(frontmatterFor(pg.Name, "Company"))
	if pg.Positions != nil && len(pg.Positions) > 0 {
		builder.WriteString("# Positions\n")
		for _, pos := range pg.Positions {
			// TODO: NOT PROPERLY SORTED
			// TODO: ORDERING?
			builder.WriteString(fmt.Sprintf("- %s\n", pos.Link()))
		}
	}
	clis := pg.Clients()
	// TODO: TEI CLIENTS NOT POPULATING
	if clis != nil && len(clis) > 0 {
		builder.WriteString("# Clients\n")
		for _, client := range clis {
			builder.WriteString(fmt.Sprintf("- %s\n", client.Link()))
		}
	}
	// resolve projects
	ps := projectsFor(clis)
	if len(ps) > 0 {
		builder.WriteString("# Projects\n")
		for _, proj := range ps {
			builder.WriteString(fmt.Sprintf("- %s\n", proj.Link()))
		}
	}
	// Resolve languages
	ls := languagesFor(ps)
	if len(ls) > 0 {
		builder.WriteString("# Languages\n")
		// TODO: SORT LANGUAGES!!!!!
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
				builder.WriteString(fmt.Sprintf("- %s\n", langs[name].Link()))
			}
		}
	}
	// ALL LOWEST
	builder.WriteString(bytesForAll(pg, false)) // TODO: languages will exist twice???
	return []byte(builder.String())
}

func (c *CompanyPage) WithPositions(positions ...*Position) *CompanyPage {
	for _, position := range slices.Backward(positions) {
		// TODO: ADD COMPANY TO SUB-ITEMS!!!!
		c.Positions = append(c.Positions, position) // TODO: ensure no double-add
		position.company = c
		for _, proj := range position.Projects { // TODO: ENSURE WORKS
			proj.TypeInfo.getClient().company = c
		}
		positionsMap[position.MapName()] = position
		tempC, tempD, tempL, tempP, tempPr, tempS, tempT := position.GetAllLowest()
		for _, val := range tempC {
			val.AddCompany(c)
		}
		for _, val := range tempD {
			val.AddCompany(c)
		}
		for key, _ := range tempL {
			langs[key].AddCompany(c)
		}
		for _, val := range tempP {
			val.AddCompany(c)
		}
		for _, val := range tempPr {
			val.AddCompany(c)
		}
		for _, svcs := range tempS {
			for _, svc := range svcs {
				svc.AddCompany(c)
			}
		}
		for _, val := range tempT {
			val.AddCompany(c)
		}
	}
	// TODO: add company to languages? Dbs? caches? providers? services? etc?
	return c
}

func (c *CompanyPage) Clients() []*Client {
	out := map[string]*Client{}
	for _, pos := range c.Positions {
		for _, pr := range pos.Projects {

			if client := pr.TypeInfo.getClient(); client != nil {
				out[client.Name] = client
			}
		}
	}
	return slices.Collect(maps.Values(out))
}

func (pg *CompanyPage) GetAllLowest() (outCaches map[string]*CachePage, outDbs map[string]*DbPage, outLanguages map[string]Frequency, outPlatforms map[string]*PlatformPage, outProviders map[string]*CloudProviderPage, outServices map[string]map[string]*CloudServicePage, outTechnologies map[string]*TechologyPage) {
	outCaches, outDbs, outLanguages, outPlatforms, outProviders, outServices, outTechnologies = initLowest()
	for _, proj := range pg.Clients() {
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

func initCompaniesAfterPositions() {
	_ = NewCompany("Source Allies", 5, 2022, nil, nil).
		WithPositions(sai1, sai2, sai3)
	_ = NewCompany("Freelance", 1, 2012, utils.Pointer(6), utils.Pointer(2022)). // TODO: RENAME
											WithPositions(freelancePosition)
	_ = NewCompany("TEI", 1, 2019, utils.Pointer(3), utils.Pointer(2020)). // TODO: ENSURE DATES ARE RIGHT
										WithPositions(positionTEI)
	_ = NewCompany("Talley Associates of Engineering", 5, 2017, utils.Pointer(3), utils.Pointer(2018)). // TODO: ENSURE DATES ARE RIGHT
														WithPositions(positionTAE)
	_ = NewCompany("Monroe Aquatics and Fitness Center", 1, 2013, utils.Pointer(8), utils.Pointer(2016)). // TODO: ENSURE DATES ARE RIGHT
														WithPositions(positionLifeguard, positionSeniorLifeguard)
}
