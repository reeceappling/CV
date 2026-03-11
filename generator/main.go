package main

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

// TODO: change top-left title that says "Quartz 4"
// TODO: remove "Aug 23, 2023 2 min read" at top of each page
// TODO: fix this link: https://cv.appli.ng/language/
// TODO: add resume
/* TODO:
---
title: Example Title
draft: false
tags:
  - example-tag
---
https://github.com/jackyzha0/quartz/blob/v4/docs/authoring%20content.md
*/

const fixmeLink = "[FIX ME](fixme.md)"

func createMainPage() { // TODO: TAGS EVERYWHERE????
	// TODO: https://quartz.jzhao.xyz/configuration PAGE TITLE
	// TODO: SPA ROUTING? // https://quartz.jzhao.xyz/configuration
	// TODO: THEMEING // https://quartz.jzhao.xyz/configuration
	// TODO: FORCE DARK MODE
	// TODO: import configs for quartz from /config
	b := strings.Builder{}
	b.WriteString(frontmatterFor("CV Home"))
	// TODO: HEADER AREA FOR LINKS TO CV, RESUME, BLOG, NOTES
	b.WriteString("# CV\n")
	b.WriteString("Welcome to my CV! It is a living document that is updated occasionally.\n")
	b.WriteString("# About\n")
	b.WriteString(fixmeLink + "\n") // TODO: SUMMARY/About

	// TODO: ADD OTHER JOBS, UNRELATED TO SOFTWARE
	b.WriteString("# Work History ([Companies](companies.md), [Positions](positions.md))\n") // TODO: RENAME
	b.WriteString("Latest Position | Company | Start Date | End Date\n")
	b.WriteString(":-- | :-- | --: | :--\n")
	for _, companyName := range companiesOrder {
		company := companies[companyName]
		endDate := company.Positions[0].End.String()
		if endDate == "current" {
			endDate = ""
		}
		b.WriteString(fmt.Sprintf("%s | %s | %s | %s\n", company.Positions[0].Name, company.Link(), company.Positions[len(company.Positions)-1].Start.String(), endDate))
	}

	b.WriteString("# Education\n")
	for _, schoolName := range []string{
		"North Carolina State University",
		"Central Academy of Technology and Arts",
	} {
		school := schools[schoolName]
		temp := make([]string, len(school.Degrees))
		for i, d := range school.Degrees {
			temp[i] = d.StringShort()
		}
		b.WriteString(fmt.Sprintf("- %s %s\n", school.Link(), strings.Join(temp, ", ")))
	}

	b.WriteString("# Certifications\n")
	b.WriteString("Cert | Certification Date | Expiration Date | link\n")
	b.WriteString(":-- | :-- | :-- | :--\n")
	for _, cert := range []*Certification{
		certAwsSaa,
		certOWASP,
		certComtrain2,
		certComtrain,
		certLifeguard,
		certCPR,
	} {
		link := ""
		if cert.Link != nil {
			link = fmt.Sprintf(`[link](%s)`, *cert.Link)
		}
		expiryDate := cert.ExpiryDate.String()
		if expiryDate == "current" {
			expiryDate = "unknown"
		}
		b.WriteString(fmt.Sprintf("%s | %s | %s | %s\n", cert.Name, cert.CertDate.String(), expiryDate, link))
	}
	// TODO: any others?

	b.WriteString("# Languages\n[Full Page](Languages.md)\n")
	b.WriteString("[see all languages LINK BROKEN](language/)\n") // TODO: delete?
	b.WriteString("## Preferred (in order)\n")
	for _, name := range []string{"Go", "Typescript", "Terraform", "Bash"} {
		b.WriteString(fmt.Sprintf("- %s\n", langs[name].Link()))
	}
	b.WriteString("## All\n")
	b.WriteString(alphabetizedLinksCompressed(langs, "language"))

	b.WriteString("# Databases\n[Full Page](dbs.md)\n\n") // TODO: why does this take 2 newlines??
	b.WriteString(alphabetizedLinksCompressed(dbs, "db"))

	b.WriteString("# Caches\n[Full Page](caches.md)\n\n") // TODO: why does this take 2 newlines??
	b.WriteString(alphabetizedLinksCompressed(caches, "cache"))

	b.WriteString("# Cloud Providers\n[Full Page](providers.md)\n")
	b.WriteString("## Preferred\n")
	for _, name := range []string{"AWS"} {
		b.WriteString(fmt.Sprintf("- %s\n", providers[name].Link()))
	}
	b.WriteString("## All\n")
	b.WriteString(alphabetizedLinksCompressed(providers, "provider"))
	// TODO: ADD GCP! Azure!

	b.WriteString("# Technologies and Libraries\n[Full Page](technologies.md)\n\n") // TODO: why does this take 2 newlines??
	b.WriteString(alphabetizedLinksCompressed(techs, "technology"))                 // TODO: dir correct?

	b.WriteString("## Containerization and Distributed Computing\n\n") // TODO: why does this take 2 newlines??
	// TODO: POPULATE THIS AREA!!!!!!!!!!!!!

	b.WriteString("# Platforms\n[Full Page](platforms.md)\n\n") // TODO: why does this take 2 newlines??
	b.WriteString(alphabetizedLinksCompressed(platforms, "platform"))

	b.WriteString("# CI/CD\n")
	cicdThings := map[string]string{
		"Github Actions": "technology",
		"Gitlab CI":      "technology",
		"Drone CI":       "technology",
	}
	temp := make([]string, len(cicdThings))
	for i, key := range slices.Collect(maps.Keys(cicdThings)) {
		temp[i] = fmt.Sprintf("[%s](%s/%s)", key, cicdThings[key], withoutSpaces(key))
	}
	b.WriteString(strings.Join(temp, ", ") + "\n")

	b.WriteString("# Observability and Monitoring\n")
	b.WriteString("## Observability\n")
	observabilityThings := map[string]string{
		"Cloudwatch": "service",  // TODO: ok?
		"Datadog":    "platform", // TODO: ok?
		"Grafana":    "platform", // TODO: ok?
		"Prometheus": "platform", // TODO: ok?
		"ServiceNow": "platform", // TODO: ok?
		// TODO: Elastic/Elasticsearch
	}
	tempObs := make([]string, len(observabilityThings))
	for i, key := range slices.Collect(maps.Keys(observabilityThings)) {
		tempObs[i] = fmt.Sprintf("[%s](%s/%s)", key, observabilityThings[key], withoutSpaces(key))
	}
	b.WriteString(strings.Join(tempObs, ", ") + "\n")
	b.WriteString("## Logs\n") // TODO: Logs
	loggingThings := map[string]string{
		"LogCentral": "platform", // TODO: ok?
	}
	tempLogging := make([]string, len(loggingThings))
	for i, key := range slices.Collect(maps.Keys(loggingThings)) {
		tempLogging[i] = fmt.Sprintf("[%s](%s/%s)", key, loggingThings[key], withoutSpaces(key))
	}
	b.WriteString(strings.Join(tempLogging, ", ") + "\n")

	b.WriteString("# Miscellaneous Skills\n")
	alphabetizedSkills := slices.Collect(maps.Keys(miscSkills))
	sort.Strings(alphabetizedSkills)
	for _, name := range alphabetizedSkills {
		b.WriteString(fmt.Sprintf("- %s\n", miscSkills[name].Link()))
	}
	b.WriteString("# Interests\n") // TODO: move????
	alphabetizedInterests := slices.Collect(maps.Keys(interests))
	sort.Strings(alphabetizedInterests)
	for _, name := range alphabetizedInterests {
		b.WriteString(fmt.Sprintf("- %s\n", linkForInterest(name)))
	}
	b.WriteString("# Subject Matters\n") // TODO: move????
	alphabetizedSms := subjectMatters.AsSlice()
	sort.Slice(alphabetizedSms, func(i, j int) bool { // TODO: ok????
		return string(alphabetizedSms[i]) < string(alphabetizedSms[j])
	})
	for _, name := range alphabetizedSms {
		b.WriteString(fmt.Sprintf("- %s\n", name.Link()))
	}

	b.WriteString("# Blog\n")
	b.WriteString("[Blog](blog/blog.md)\n")
	b.WriteString("# Notes\n")
	b.WriteString("[Notes](notes/notes.md)\n")
	if err := os.WriteFile(root+"index.md", []byte(b.String()), 777); err != nil {
		panic("failed to write main.md file: " + err.Error())
	}
}

func createErrorPage() {
	b := strings.Builder{}
	b.WriteString(frontmatterFor("Not Found Page"))
	b.WriteString("Page does not exist!\n") // TODO: PUT LINK TO HOME ON EVERY PAGE
	if err := os.WriteFile(root+"error.md", []byte(b.String()), 777); err != nil {
		panic("failed to write main.md file: " + err.Error())
	}
}

// Company -> position -> client -> project ->>>>
// languages, platforms, providers(and services), dbs, caches

// TODO: PUT TAGS ALL OVER PROJECTS!!!!!!

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
			for _, svc := range slices.Collect(maps.Keys(svcs)) {
				b.WriteString(fmt.Sprintf("%s\n", cloudServices[svc].Link()))
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
}

func newTracked() *tracked {
	return &tracked{
		Companies: map[string]*CompanyPage{},
		Positions: map[string]*Position{},
		Clients:   map[string]*Client{},
		Schools:   map[string]*SchoolPage{},
		Projects:  map[string]*Project{},
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
		allKeys[i] = fmt.Sprintf("[%s](%s/%s)", db, dir, withoutSpaces(db))
	}
	return strings.Join(allKeys, ", ") + "\n"
}
func alphabetizedLinks[T any](inpMap map[string]T, dir string) string {
	allKeys := slices.Collect(maps.Keys(inpMap))
	sort.Strings(allKeys)
	for i, db := range allKeys {
		allKeys[i] = fmt.Sprintf("- [%s](%s/%s)\n", db, dir, withoutSpaces(db))
	}
	return strings.Join(allKeys, "") + "\n"
}

func main() {
	if err := os.MkdirAll("./quartz/content", 777); err != nil {
		panic("failed to create content dir: " + err.Error())
	}
	for _, dir := range []string{"school", "client", "project", "company", "position", "language", "platform", "db", "cache", "provider", "service", "technology", "miscSkill"} {
		if err := os.MkdirAll("./quartz/content/"+dir, 777); err != nil {
			panic("failed to create " + dir + " dir: " + err.Error())
		}
	}

	// Create data structures representing the content of the pages to write
	// initSchools() // Done outside of init
	// initClientsFirst() // Done in vars
	//initProjectsAfterClients()
	initProjectsFinal()
	initClientsAfterProjectsComplete() // Sets client on projects as well // TODO: may need to go after positions (worked that way before)
	initPositionsAfterProjects()
	initCompaniesAfterPositions() // Must be done after positions and project setup, but before projects pages. What about clients?

	// TODO: SEARCH BAR?

	// Start creating actual pages
	createLanguagesPages()
	createSchoolsPages()
	createPositionsPages()
	createProjectsPages()
	createClientsPages()
	createCompaniesPages()
	createPlatformsPages()
	createDbsPages()
	createCachesPages()
	createProvidersPages()
	createServicesPages()
	createTechnologiesPages()
	createMiscSkillsPages()
	createInterestsPages()
	createSubjectMatterPages()
	createMainPage()
	createErrorPage()

}

func createLanguagesPages() {
	b := strings.Builder{}
	b.WriteString(frontmatterFor("Languages"))
	b.WriteString("# Preferred (in order)\n")
	for _, name := range []string{"Go", "Typescript", "Terraform", "Bash"} {
		l := langs[name]
		b.WriteString(fmt.Sprintf("- %s\n", l.Link()))
	}
	b.WriteString("# All\n")
	b.WriteString("[see all languages](language/)\n") // TODO: delete?
	b.WriteString(alphabetizedLinks(langs, "language"))
	err := writeFileFromScratch("Languages.md", b.String())
	if err != nil {
		panic(err.Error())
	}
	if len(langs) > 0 {
		if err = os.MkdirAll(root+"language", 777); err != nil {
			panic("failed to create languages dir: " + err.Error())
		}
	}
	for name, lang := range langs {
		writePage("language", withoutSpaces(name)+".md", lang.Bytes())
	}
}

func writeFileFromScratch(filepathFromRoot, toWrite string) error {
	err := os.WriteFile(root+filepathFromRoot, []byte(toWrite), 777) // TODO; FAILING HERE
	if err != nil {
		return errors.Join(errors.New("failed to create "+filepathFromRoot), err)
	}
	return nil
}

func createProjectsPages() {
	b, bPersonal, bProfessional, bSchool := strings.Builder{}, strings.Builder{}, strings.Builder{}, strings.Builder{}

	// TODO: FIX PROJECTS PAGE!!!
	b.WriteString(frontmatterFor("Projects"))
	bProfessional.WriteString("# Professional Projects\n")
	bProfessional.WriteString("Project | Company | Client\n")
	bProfessional.WriteString(":-- | :--: | :--\n")
	bPersonal.WriteString("# Personal Projects\n")
	bSchool.WriteString("# Coursework-related Projects\n")

	for _, proj := range projects {
		switch proj.TypeInfo.Type() {
		case projectTypePersonal:
			bPersonal.WriteString(fmt.Sprintf("- %s\n", proj.Link()))
		case projectTypeProfessional:
			pr := proj.Link()
			var co, cl = "none", "none"
			client := proj.TypeInfo.getClient()
			if client != nil {
				cl = client.Link()
				if client.company != nil {
					co = client.company.Link()
				}
			}

			bProfessional.WriteString(fmt.Sprintf("%s | %s | %s\n", pr, co, cl))
		case projectTypeSchool:
			bSchool.WriteString(fmt.Sprintf("- %s\n", proj.Link())) // TODO: school link???
		default:
			panic("unknown project type")

		}
	}
	b.WriteString(bProfessional.String())
	b.WriteString(bPersonal.String())
	b.WriteString(bSchool.String())

	err := os.WriteFile(root+"projects.md", []byte(b.String()), 777)
	if err != nil {
		panic(err)
	}
	if len(langs) > 0 {
		if err = os.MkdirAll(root+"project", 777); err != nil {
			panic("failed to create projects dir: " + err.Error())
		}
	}
	for name, item := range projects {
		writePage("project", withoutSpaces(name)+".md", item.Bytes())
	}
}

func createSchoolsPages() {
	// Create schools page
	b := strings.Builder{}
	b.WriteString(frontmatterFor("Schools"))
	for name, school := range schools {
		b.WriteString(fmt.Sprintf("# %s\n", school.Link()))
		// Write all degrees
		b.WriteString(fmt.Sprintf("- Degrees:", name, withoutSpaces(name)))
		for i, deg := range school.Degrees {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(deg.StringShort())
		}
		b.WriteString("\n")
		// Write all projects
		b.WriteString(fmt.Sprintf("- Projects:\n", name, withoutSpaces(name)))
		for _, proj := range school.Projects {
			b.WriteString(fmt.Sprintf("- - %s\n", proj.Link()))
		}
		b.WriteString(fmt.Sprintf("- Extracurriculars:\n", name, withoutSpaces(name)))
		// Write all extracurriculars
		for _, ec := range school.Extracurriculars {
			b.WriteString(fmt.Sprintf("- - %s\n", ec.Name))
		}
	}
	// Write page for each school
	for name, item := range schools {
		writePage("school", withoutSpaces(name)+".md", item.Bytes())
	}
}

func createPositionsPages() {
	b := strings.Builder{}
	b.WriteString(frontmatterFor("Positions"))
	b.WriteString("Most Current\n\n")
	b.WriteString("Title | Company | StartDate | EndDate\n")
	b.WriteString(":-- | :--: | :--: | :--\n")
	psMap := positionsMap
	if len(psMap) == 0 {
		println("empty positions map")
	}
	for _, pos := range slices.SortedFunc(maps.Values(positionsMap), func(a *Position, b *Position) int {
		if a.Start.Year == b.Start.Year {
			if a.Start.Month > b.Start.Month {
				return -1
			}
			if a.Start.Month < b.Start.Month {
				return 1
			}
			return 0
		} else {
			if a.Start.Year > b.Start.Year {
				return -1
			}
			if a.Start.Year < b.Start.Year {
				return 1
			}
			return 0
		}
	}) {
		end := pos.End.String()
		if end == "current" {
			end = ""
		}
		b.WriteString(fmt.Sprintf(" %s | %s | %s | %s\n", pos.Link(), pos.company.Name, pos.Start.String(), end))
	}
	b.WriteString("Earliest")

	err := os.WriteFile(root+"positions.md", []byte(b.String()), 777)
	if err != nil {
		panic(err)
	}
	if len(langs) > 0 {
		if err = os.MkdirAll(root+"position", 777); err != nil {
			panic("failed to create position dir: " + err.Error())
		}
	}
	for name, item := range positionsMap {
		writePage("position", withoutSpaces(name)+".md", item.Bytes())
	}
}

func createClientsPages() {
	// Clients pages
	b := strings.Builder{}
	b.WriteString(frontmatterFor("Clients"))
	b.WriteString("(most recent at top)\n\nClient | Company\n")
	b.WriteString(":-- | :--\n")
	for _, clientName := range clientsOrder {
		client := clients[clientName]
		compLink := fmt.Sprintf("[%s](company/%s)", client.Name, withoutSpaces(client.Name))
		if client.company != nil {
			compLink = client.company.Link()
		}
		b.WriteString(fmt.Sprintf("%s | %s \n", client.Link(), compLink))
	}
	err := os.WriteFile(root+"clients.md", []byte(b.String()), 777)
	if err != nil {
		panic("failed to create clients page: " + err.Error())
	}
	if len(clients) > 0 {
		if err = os.MkdirAll(root+"client", 777); err != nil {
			panic("failed to create clients dir: " + err.Error())
		}
	}
	for name, item := range clients {
		writePage("client", withoutSpaces(name)+".md", item.Bytes())
	}
}

func createCompaniesPages() {
	// Companies pages
	b := strings.Builder{}
	b.WriteString(frontmatterFor("Companies"))
	b.WriteString("Latest Position | Company | Company Start Date | Company End Date\n")
	b.WriteString(":-- | :-- | --: | :--\n")
	// TODO: make this into a cool chart!
	for _, companyName := range companiesOrder {
		// TODO: Consider linking every position from here, but on separate lines
		company := companies[companyName]
		b.WriteString(fmt.Sprintf("%s | %s | %s | %s\n", company.Positions[0].Name, company.Link(), company.Start.String(), company.End.String()))
	}
	err := os.WriteFile(root+"companies.md", []byte(b.String()), 777)
	if err != nil {
		panic("failed to write companies file: " + err.Error())
	}
	if len(companies) > 0 {
		if err = os.MkdirAll(root+"company", 777); err != nil {
			panic("failed to create company dir: " + err.Error())
		}
	}
	for name, company := range companies {
		writePage("company", withoutSpaces(name)+".md", company.Bytes())
	}
}

type Linkable interface {
	Link() string
}

func createPlatformsPages() {
	// Platforms pages
	b := strings.Builder{}
	b.WriteString(frontmatterFor("Platforms"))
	alphabetizedPlatforms := slices.Collect(maps.Keys(platforms))
	sort.Strings(alphabetizedPlatforms)
	for _, platName := range alphabetizedPlatforms {
		platform := platforms[platName]
		b.WriteString(fmt.Sprintf("- %s\n", platform.Link()))
	}
	err := os.WriteFile(root+"platforms.md", []byte(b.String()), 777)
	if err != nil {
		panic("failed to write platforms file: " + err.Error())
	}
	if len(platforms) > 0 {
		if err = os.MkdirAll(root+"platform", 777); err != nil {
			panic("failed to create platforms dir: " + err.Error())
		}
	}
	for name, platform := range platforms {
		writePage("platform", withoutSpaces(name)+".md", platform.Bytes(name, "Platform"))
	}
}

func createDbsPages() {
	// Databases pages
	b := strings.Builder{}
	b.WriteString(frontmatterFor("Databases"))
	alphabetized := slices.Collect(maps.Keys(dbs))
	sort.Strings(alphabetized)
	for _, name := range alphabetized {
		db := dbs[name]
		b.WriteString(fmt.Sprintf("- %s\n", db.Link()))
	}
	err := os.WriteFile(root+"dbs.md", []byte(b.String()), 777)
	if err != nil {
		panic("failed to write dbs file: " + err.Error())
	}
	if len(dbs) > 0 {
		if err = os.MkdirAll(root+"db", 777); err != nil {
			panic("failed to create db dir: " + err.Error())
		}
	}
	for name, db := range dbs {
		writePage("db", withoutSpaces(name)+".md", db.Bytes(name, "Database"))
	}
}
func createCachesPages() {
	// Caches pages
	b := strings.Builder{}
	b.WriteString(frontmatterFor("Caches"))
	alphabetized := slices.Collect(maps.Keys(caches))
	sort.Strings(alphabetized)
	for _, name := range alphabetized {
		cache := caches[name]
		b.WriteString(fmt.Sprintf("- %s\n", cache.Link()))
	}
	err := os.WriteFile(root+"caches.md", []byte(b.String()), 777)
	if err != nil {
		panic("failed to write caches file: " + err.Error())
	}
	if len(caches) > 0 {
		if err = os.MkdirAll(root+"cache", 777); err != nil {
			panic("failed to create cache dir: " + err.Error())
		}
	}
	for name, cache := range caches {
		writePage("cache", withoutSpaces(name)+".md", cache.Bytes(name, "Cache"))
	}
}

func createProvidersPages() {
	// Providers pages
	b := strings.Builder{}
	b.WriteString(frontmatterFor("Cloud Providers"))
	// TODO: PREFERRED VS ALL
	b.WriteString("# Preferred (in order)\n")
	for _, name := range []string{"AWS"} {
		provider := providers[name]
		b.WriteString(fmt.Sprintf("- %s\n", provider.Link()))
		// TODO: SUB-SERVICES!!!!!
	}
	b.WriteString("# All\n")
	b.WriteString(alphabetizedLinksCompressed(providers, "provider")) // TODO: SUB-SERVICES!!!!
	err := os.WriteFile(root+"providers.md", []byte(b.String()), 777)
	if err != nil {
		panic("failed to write providers file: " + err.Error())
	}
	if len(providers) > 0 {
		if err = os.MkdirAll(root+"provider", 777); err != nil {
			panic("failed to create provider dir: " + err.Error())
		}
	}
	for name, provider := range providers {
		writePage("provider", withoutSpaces(name)+".md", provider.Bytes())
	}
}

func createServicesPages() {
	// Services pages
	b := strings.Builder{}
	b.WriteString(frontmatterFor("Cloud Services"))
	b.WriteString("Service | Provider\n")
	b.WriteString(":-- | :--\n")
	for _, service := range cloudServices { // TODO: alphabetize?!!!!!
		b.WriteString(fmt.Sprintf("%s | %s\n", service.Link(), service.provider))
	}
	err := os.WriteFile(root+"services.md", []byte(b.String()), 777)
	if err != nil {
		panic("failed to write services file: " + err.Error())
	}
	if len(cloudServices) > 0 {
		if err = os.MkdirAll(root+"service", 777); err != nil {
			panic("failed to create service dir: " + err.Error())
		}
	}
	for name, service := range cloudServices {
		writePage("service", withoutSpaces(name)+".md", service.Bytes(name, "Cloud Service"))
	}
}

func createTechnologiesPages() {
	// Technologies pages
	b := strings.Builder{}
	b.WriteString(frontmatterFor("Technologies"))
	alphabetized := slices.Collect(maps.Keys(techs))
	sort.Strings(alphabetized)
	for _, name := range alphabetized {
		tech := techs[name]
		b.WriteString(fmt.Sprintf("- %s\n", tech.Link()))
	}
	err := os.WriteFile(root+"technologies.md", []byte(b.String()), 777)
	if err != nil {
		panic("failed to write technologies file: " + err.Error())
	}
	if len(techs) > 0 {
		if err = os.MkdirAll(root+"technology", 777); err != nil {
			panic("failed to create technology dir: " + err.Error())
		}
	}
	for name, tech := range techs {
		writePage("technology", withoutSpaces(name)+".md", tech.Bytes(name, "Technology"))
	}
}

func createMiscSkillsPages() {
	// Misc skills pages
	b := strings.Builder{}
	b.WriteString(frontmatterFor("Miscellaneous Skills"))
	alphabetized := slices.Collect(maps.Keys(miscSkills))
	sort.Strings(alphabetized)
	for _, name := range alphabetized {
		skill := miscSkills[name]
		b.WriteString(fmt.Sprintf("- %s\n", skill.Link()))
	}
	err := os.WriteFile(root+"miscSkills.md", []byte(b.String()), 777)
	if err != nil {
		panic("failed to write miscSkills file: " + err.Error())
	}
	if len(miscSkills) > 0 {
		if err = os.MkdirAll(root+"miscSkill", 777); err != nil {
			panic("failed to create miscSkill dir: " + err.Error())
		}
	}
	for name, skill := range miscSkills {
		writePage("miscSkill", withoutSpaces(name)+".md", skill.Bytes(name, "Misc Skill"))
	}
}
func createInterestsPages() {
	b := strings.Builder{}
	b.WriteString(frontmatterFor("Interests"))
	alphabetizedInterests := slices.Collect(maps.Keys(interests))
	sort.Strings(alphabetizedInterests)
	for _, name := range alphabetizedInterests {
		b.WriteString(fmt.Sprintf("- %s\n", linkForInterest(name)))
	}
	err := os.WriteFile(root+"interests.md", []byte(b.String()), 777)
	if err != nil {
		panic("failed to write interests file: " + err.Error())
	}
	if len(interests) > 0 {
		if err = os.MkdirAll(root+"interest", 777); err != nil {
			panic("failed to create interest dir: " + err.Error())
		}
	}
	for name, interest := range interests {
		writePage("interest", withoutSpaces(name)+".md", interest.Bytes(name, "Interest"))
	}
}
func createSubjectMatterPages() {
	// TODO; THIS!!!!
	b := strings.Builder{}
	b.WriteString(frontmatterFor("Subject Matters"))
	alphabetizedSms := subjectMatters.AsSlice()
	sort.Slice(alphabetizedSms, func(i, j int) bool { // TODO: ok????
		return string(alphabetizedSms[i]) < string(alphabetizedSms[j])
	})
	for _, name := range alphabetizedSms {
		b.WriteString(fmt.Sprintf("- %s\n", name.Link()))
	}
	err := os.WriteFile(root+"subjectMatters.md", []byte(b.String()), 777)
	if err != nil {
		panic("failed to write subjectMatters file: " + err.Error())
	}
	if len(subjectMatters) > 0 {
		if err = os.MkdirAll(root+"subjectMatter", 777); err != nil {
			panic("failed to create subjectMatter dir: " + err.Error())
		}
	}
	for sm, _ := range subjectMatters {
		// TODO: LIKELY USE A TAG SYSTEM INSTEAD!!!!!
		writePage("subjectMatter", string(sm)+".md", []byte("SUBJECT MATTER PAGE NOT IMPLEMENTED") /*skill.Bytes()*/) // TODO: FIX!
	}
}

func writePage(dirName, filename string, bs []byte) {
	path := dirName + "/" + filename
	err := os.WriteFile(root+path, bs, 777)
	if err != nil {
		panic(fmt.Sprintf("failed to write page %s: %s", path, err.Error()))
	}
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
