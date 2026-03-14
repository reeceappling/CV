package main

import (
	"appli.ng/cv/generator/utils"
	"fmt"
	"maps"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

func frontmatterFor(title string, tags ...string) string {
	b := strings.Builder{}
	b.WriteString("---\n")
	b.WriteString("title: " + title + "\n")
	b.WriteString("draft: false\n")
	if len(tags) > 0 {
		b.WriteString("tags:\n")
		for _, tag := range tags {
			b.WriteString("  - " + tag + "\n")
		}
	}
	b.WriteString("---\n")
	return b.String()
}

const fixmeLink = "[FIX ME](error.md)"

func linkFor(text string, path ...string) string {
	return fmt.Sprintf("[%s](%s)", text, strings.Join(path, "/"))
}

func WriteFile(filenameLessRoot string, content string) {
	if err := os.WriteFile(root+filenameLessRoot, []byte(content), 777); err != nil {
		panic("failed to write " + filenameLessRoot + " file: " + err.Error())
	}
}
func WriteCVFile(filenameLessRootAndCv string, content string) {
	WriteFile("cv/"+filenameLessRootAndCv, content)
}
func MakeCVDir(dir string) {
	if err := os.MkdirAll(root+"cv/"+dir, 777); err != nil {
		panic("failed to create cv/" + dir + " dir: " + err.Error())
	}
}
func writeMultipleTextLines(b *strings.Builder, lines ...string) {
	b.WriteString(strings.Join(lines, "\n\n") + "\n")
}

type showAll interface {
	GetAllLowest() (outCaches map[string]*CachePage, outDbs map[string]*DbPage, outLanguages map[string]Frequency, outPlatforms map[string]*PlatformPage, outProviders map[string]*CloudProviderPage, outServices map[string]map[string]*CloudServicePage, outTechnologies map[string]*TechologyPage)
	NameValue() string
}

func bytesForAll(item showAll, showFrequencies bool) string {
	b := strings.Builder{}
	tempC, tempD, tempL, tempP, _, tempS, tempT := item.GetAllLowest()
	if len(tempL) > 0 {
		b.WriteString("# Languages\n")
		showFreqLink := ""
		if showFrequencies {
			showFreqLink = "[[" + withoutSpaces(item.NameValue()) + "#^usageFrequencies\\|*]]"
		}
		b.WriteString("Language | Usage Frequency" + showFreqLink + "\n")
		b.WriteString(":-- | :--\n")
		freqs := make(map[Frequency][]string, 6)
		for name, freq := range tempL {
			if freqs[freq] == nil {
				freqs[freq] = []string{name}
			} else {
				freqs[freq] = append(freqs[freq], name)
			}
		}
		for f := 5; f >= 0; f-- {
			for _, name := range freqs[Frequency(f)] {
				b.WriteString(fmt.Sprintf("%s | %s\n", langs[name].Link(), Frequency(f).String()))
			}
		}
	}
	if len(tempC) > 0 {
		b.WriteString("# Caches\n")
		for name, _ := range tempC {
			b.WriteString(fmt.Sprintf("%s\n", caches[name].Link()))
		}
	}
	if len(tempD) > 0 {
		b.WriteString("# Databases\n")
		for name, _ := range tempD {
			b.WriteString(fmt.Sprintf("%s\n", dbs[name].Link()))
		}
	}
	if len(tempP) > 0 {
		b.WriteString("# Platforms\n")
		for name, _ := range tempP {
			b.WriteString(fmt.Sprintf("%s\n", platforms[name].Link()))
		}
	}
	if len(tempT) > 0 {
		b.WriteString("# Technologies\n")
		for name, _ := range tempT {
			b.WriteString(fmt.Sprintf("%s\n", techs[name].Link()))
		}
	}
	if len(tempS) > 0 {
		b.WriteString("# Providers\n")
		for name, svcs := range tempS {
			tempSvcs := make([]string, len(svcs))
			for i, svc := range slices.Collect(maps.Keys(svcs)) {
				tempSvcs[i] = fmt.Sprintf("%s\n", cloudServices[svc].Link())
			}
			b.WriteString(fmt.Sprintf("%s: %s\n", providers[name].Link(), strings.Join(tempSvcs, ", ")))
		}
	}
	if showFrequencies {
		b.WriteString(definitionsArea())
		b.WriteString(usageFrequencyDefinitions())
	}
	return b.String()
}

type tracked struct {
	Companies map[string]*CompanyPage
	Positions map[string]*Position
	Clients   map[string]*Client
	Schools   map[string]*SchoolPage
	Projects  map[string]*Project
	Tags      utils.Set[string]
}

//type trackable interface {
//	Bytes(title string, typ string) []byte
//	AddClient(client *Client)
//	AddProject(proj *Project)
//	AddPosition(pos *Position)
//	AddCompany(comp *CompanyPage)
//	String() string
//}

func newTracked() *tracked {
	return &tracked{
		Companies: map[string]*CompanyPage{},
		Positions: map[string]*Position{},
		Clients:   map[string]*Client{},
		Schools:   map[string]*SchoolPage{},
		Projects:  map[string]*Project{},
		Tags:      utils.Set[string]{},
	}
}
func (pg *tracked) Bytes(title string, typ string) []byte {
	builder := strings.Builder{}
	if title != "" {
		builder.WriteString(frontmatterFor(title, typ))
	}

	if pg.Companies != nil && len(pg.Companies) > 0 {
		builder.WriteString("# Companies\n")
		for _, com := range pg.Companies {
			builder.WriteString(fmt.Sprintf("- %s\n", com.Link()))
		}
	}
	if pg.Positions != nil && len(pg.Positions) > 0 {
		builder.WriteString("# Positions\n")
		for _, pos := range pg.Positions {
			builder.WriteString(fmt.Sprintf("- %s\n", pos.Link()))
		}
	}
	if pg.Clients != nil && len(pg.Clients) > 0 {
		builder.WriteString("# Clients\n")
		for _, cli := range pg.Clients {
			builder.WriteString(fmt.Sprintf("- %s\n", cli.Link()))
		}
	}
	if pg.Projects != nil && len(pg.Projects) > 0 {
		builder.WriteString("# Projects\n")
		//builder.WriteString("Project | Usage Frequency | Project Type\n")
		//builder.WriteString(":-- | :--: | :--\n")
		for _, proj := range pg.Projects {
			builder.WriteString(fmt.Sprintf("- %s\n", proj.Link()))
		}
	}
	return []byte(builder.String())
}

func (pg *tracked) AddClient(client *Client) {
	pg.Clients[client.Name] = client
}
func (pg *tracked) AddSchool(school *SchoolPage) {
	pg.Schools[school.Name] = school
}
func (pg *tracked) AddProject(proj *Project) {
	pg.Projects[proj.Name] = proj
}

func (pg *tracked) AddPosition(pos *Position) {
	pg.Positions[pos.Name] = pos
}
func (pg *tracked) AddCompany(comp *CompanyPage) {
	pg.Companies[comp.Name] = comp
}

func (tr *tracked) String() string {
	if tr == nil {
		return ""
	}
	b := strings.Builder{}
	if tr.Projects != nil {
		b.WriteString("# Projects\n")
		for _, proj := range tr.Projects {
			b.WriteString(fmt.Sprintf("- %s\n", proj.Link()))
		}
	}
	if tr.Positions != nil {
		b.WriteString("# Positions\n")
		for _, pos := range tr.Positions { // TODO: ENSURE THESE ARE IN ORDER
			b.WriteString(fmt.Sprintf("- %s at %s\n", pos.Link(), pos.company.Link()))
		}
	}
	if tr.Companies != nil {
		b.WriteString("# Companies\n")
		for _, item := range tr.Companies {
			b.WriteString(fmt.Sprintf("- %s\n", item.Link()))
		}
	}
	if tr.Clients != nil {
		b.WriteString("# Clients\n")
		for _, client := range tr.Clients {
			b.WriteString(fmt.Sprintf("- %s\n", client.Link()))
		}
	}
	return b.String()
}

// Create project, add project to client, add project to position, add position to company (client will be automatically linked this way)
const root = "./quartz/content/"

func alphabetizedLinksCompressed[T any](inpMap map[string]T, dir string) string {
	allKeys := slices.Collect(maps.Keys(inpMap))
	sort.Strings(allKeys)
	for i, db := range allKeys {
		allKeys[i] = linkFor(db, "cv", dir, withoutSpaces(db))
	}
	return strings.Join(allKeys, ", ") + "\n"
}
func alphabetizedLinks[T any](inpMap map[string]T, dir string) string {
	allKeys := slices.Collect(maps.Keys(inpMap))
	sort.Strings(allKeys)
	for i, db := range allKeys {
		allKeys[i] = linkFor(db, "cv", dir, withoutSpaces(db))
	}
	return strings.Join(allKeys, "") + "\n"
}

func withoutSpaces(s string) string {
	return strings.Join(strings.Split(s, " "), "_")
}

type Frequency int

const (
	Minimal = iota
	Rarely
	Some
	Regularly
	Often
	Extensively
)

func (f Frequency) String() string {
	switch f {
	case Minimal:
		return "Minimal"
	case Rarely:
		return "Rarely"
	case Some:
		return "Some"
	case Regularly:
		return "Regularly"
	case Often:
		return "Often"
	case Extensively:
		return "Extensively"
	default:
		panic("unhandled frequency: " + strconv.Itoa(int(f)))
	}
}

func FrequencyFromString(s string) Frequency {
	switch s {
	case "Minimal":
		return Minimal
	case "Rarely":
		return Rarely
	case "Some":
		return Some
	case "Regularly":
		return Regularly
	case "Often":
		return Often
	case "Extensively":
		return Extensively
	default:
		panic("unhandled frequency: " + s)
	}
}

type Linkable interface { // TODO: use???
	EntryType() string
	Link() string
}

var linkables = map[string]Linkable{} // TODO: USE THIS
type NoLinkable struct{}

func (n NoLinkable) EntryType() string {
	panic("nilEntryType")
}

func (n NoLinkable) Link() string {
	return fixmeLink
}

func addLinkable(name string, item Linkable) {
	linkables[strings.ToLower(name)] = item
}

func lookup(name string) Linkable {
	if item, ok := linkables[strings.ToLower(name)]; ok {
		return item
	}
	panic("lookup failed for " + name) // TODO: del?
	return NoLinkable{}
}
