package main

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"sort"
	"strings"
)

// TODO: LAYOUT
// TODO: remove "Aug 23, 2023 2 min read" at top of each page
// TODO: fix this link: https://reece.appli.ng/cv/language/
// TODO: add resume
// TODO: Figure out how deep the graph view should be...
// TODO: DesktopOnly.tsx???
// TODO: sort projects on projects page? (maybe chronological?) (alphabetical?)
// TODO: backlinks only on the subject matter pages????
// TODO: PUT TAGS ALL OVER PROJECTS???!!!!!!

func main() {
	for _, dir := range []string{"cv", "blog", "note"} {
		if err := os.MkdirAll("./quartz/content/"+dir, 777); err != nil {
			panic("failed to create content/" + dir + " dir: " + err.Error())
		}
	}
	for _, dir := range []string{"subjectMatter", "interest", "school", "client", "project", "company", "position", "language", "platform", "db", "cache", "provider", "service", "technology", "miscSkill"} {
		if err := os.MkdirAll("./quartz/content/cv/"+dir, 777); err != nil {
			panic("failed to create " + dir + " dir: " + err.Error())
		}
	}

	setupSubjectMatters()
	// Initialize low-level things with subject matters

	// Create data structures representing the content of the pages to write
	// initSchools() // Done outside of init
	// initClientsFirst() // Done in vars
	//initProjectsAfterClients()  // Done in vars?
	initProjectsFinal()
	initClientsAfterProjectsComplete() // Sets client on projects as well
	initPositionsAfterProjects()
	initCompaniesAfterPositions() // Must be done after positions and project setup, but before projects pages. What about clients?

	// TODO: POPULATE SUBJECT MATTERS ON TECHS, SERVICES, MISCSKILLS?

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
	createErrorPage()
	createMainPage()
	createMainCVPage()
	createBlogPages()
	createNotesPages()

}

func createMainPage() {
	b := &strings.Builder{}
	b.WriteString(frontmatterFor("Home - Reece Appling"))
	b.WriteString("# About\n")
	writeMultipleTextLines(b,
		"__Hi, I'm Reece__. Welcome to my personal website! This site is primarily a living document for [my CV](cv.md) and [my resume](static/Resume_Reece_Appling.pdf), both of which I try to keep relatively current. I also host [my personal blog](Blog.md) here, as well as an area to publish [my notes](Notes).",
		"This entire site is auto-generated from Go code and text files into markdown files \\(for [Obsidian](https://obsidian.md)\\), which are exported to html, css, and javascript via [Quartz 4](https://quartz.jzhao.xyz), hosted on AWS S3, and accessed via AWS CloudFront and Cloudflare. Feel free to check out the [source code](https://github.com/reeceappling/CV).",
		"__Want to get in contact with me?__ Some contact info should be at the footer of this page. Otherwise, [my links page](https://links.reece.appli.ng) contains many of my socials, as well as my email and phone number.",
	)
	b.WriteString("# CV\n")
	writeMultipleTextLines(b,
		// TODO: CV LAST UPDATED DATE?
		// TODO: RESUME LAST UPDATED DATE?
		"[Check out my CV](cv.md). It is a living document that is updated occasionally.",
		"Looking for a resume instead? [Download my resume here](static/Resume_Reece_Appling.pdf)",
	)
	b.WriteString("# Blog\n")
	b.WriteString("Here is [my blog](Blog.md) where I host any blog posts I make.\n")
	b.WriteString("# Public Notes\n")
	b.WriteString("Lastly, here are [my notes](Notes)\n")
	WriteFile("index.md", b.String())
}

func createMainCVPage() { // TODO: TAGS EVERYWHERE????
	// TODO: https://quartz.jzhao.xyz/configuration PAGE TITLE
	// TODO: THEMEING // https://quartz.jzhao.xyz/configuration
	// TODO: FORCE DARK MODE
	b := strings.Builder{}
	b.WriteString(frontmatterFor("CV"))
	// TODO: HEADER AREA FOR LINKS TO CV, RESUME, BLOG, NOTES?

	b.WriteString("Welcome to my CV! It is a living document that is updated occasionally.\n\n")                     // Why does this need 2 newlines? // TODO: move down?
	b.WriteString("Looking for a resume instead? [Download it here](Resume.pdf)ENSURE WORKING\n\n")                  // TODO: ENSURE OK
	b.WriteString("Feel free to check out the [source code](https://github.com/reeceappling/CV) for this website\n") // TODO: ENSURE OK

	b.WriteString("# About\n")

	b.WriteString(fixmeLink + "\n") // TODO: SUMMARY/About

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
	for _, cert := range certs {
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

	b.WriteString("# [Languages](Languages.md)\n")
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

	b.WriteString("## Containerization\n\n")      // TODO: why does this take 2 newlines??
	b.WriteString("## Distributed Computing\n\n") // TODO: why does this take 2 newlines??

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
		temp[i] = linkFor(key, "cv", cicdThings[key], withoutSpaces(key))
	}
	b.WriteString(strings.Join(temp, ", ") + "\n")

	b.WriteString("# Observability and Monitoring\n")
	b.WriteString("## Observability\n")
	// TODO: loop through all to find observability instead of doing it this way???
	observabilityThings := map[string]string{
		"Cloudwatch": "service",  // TODO: ok?
		"Datadog":    "platform", // TODO: ok?
		"Grafana":    "platform", // TODO: ok?
		"Prometheus": "platform", // TODO: ok?
		"ServiceNow": "platform", // TODO: ok?
		// TODO: Elastic/Elasticsearch somewhere
	}
	tempObs := make([]string, len(observabilityThings))
	for i, key := range slices.Collect(maps.Keys(observabilityThings)) {
		tempObs[i] = linkFor(key, "cv", observabilityThings[key], withoutSpaces(key))
	}
	b.WriteString(strings.Join(tempObs, ", ") + "\n")
	b.WriteString("## Logs\n") // TODO: Logs
	loggingThings := map[string]string{
		"LogCentral": "platform", // TODO: ok?
	}
	tempLogging := make([]string, len(loggingThings))
	for i, key := range slices.Collect(maps.Keys(loggingThings)) {
		tempLogging[i] = linkFor(key, "cv", loggingThings[key], withoutSpaces(key))
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
	b.WriteString("# Subject Matters\n")
	alphabetizedSms := slices.Collect(maps.Keys(subjectMatters))
	sort.Slice(alphabetizedSms, func(i, j int) bool { // TODO: ok????
		return string(alphabetizedSms[i]) < string(alphabetizedSms[j])
	})
	for _, name := range alphabetizedSms {
		b.WriteString(fmt.Sprintf("- %s\n", name.Link()))
	}
	WriteFile("cv.md", b.String())
}

func createErrorPage() {
	b := strings.Builder{}
	b.WriteString(frontmatterFor("Not Found Page"))
	b.WriteString("Page does not exist!\n") // TODO: PUT LINK TO HOME ON EVERY PAGE
	WriteFile("error.md", b.String())
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
	for _, lang := range langs {
		if len(lang.SubjectMatters) == 0 {
			panic("lang " + lang.Name + "has no subject matters")
		}
	}
	b.WriteString(alphabetizedLinks(langs, "language"))
	WriteCVFile("Languages.md", b.String())
	if len(langs) > 0 {
		MakeCVDir("language")
	}
	for name, lang := range langs {
		if len(lang.SubjectMatters) == 0 {
			panic("no subject matter on language " + name)
		}
		WriteCVFile("language/"+withoutSpaces(name)+".md", string(lang.Bytes()))
	}
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

	WriteCVFile("projects.md", b.String())
	if len(langs) > 0 {
		MakeCVDir("project")
	}
	for name, item := range projects {
		if len(item.SubjectMatters) == 0 {
			panic("no subject matter on project " + name)
		}
		WriteCVFile("project/"+withoutSpaces(name)+".md", string(item.Bytes()))
	}
}

func createSchoolsPages() {
	// Create schools page
	b := strings.Builder{}
	b.WriteString(frontmatterFor("Schools"))
	for _, school := range schools {
		b.WriteString(fmt.Sprintf("# %s\n", school.Link()))
		// Write all degrees
		b.WriteString("- Degrees:")
		for i, deg := range school.Degrees {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(deg.StringShort())
		}
		b.WriteString("\n")
		// Write all projects
		b.WriteString("- Projects:\n")
		for _, proj := range school.Projects {
			b.WriteString(fmt.Sprintf("- - %s\n", proj.Link()))
		}
		b.WriteString("- Extracurriculars:\n")
		// Write all extracurriculars
		for _, ec := range school.Extracurriculars {
			b.WriteString(fmt.Sprintf("- - %s\n", ec.Name))
		}
	}
	// Write page for each school
	for name, item := range schools {
		if len(item.SubjectMatters) == 0 {
			panic("no subject matter on school " + name)
		}
		WriteCVFile("school/"+withoutSpaces(name)+".md", string(item.Bytes()))
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

	WriteFile("index.md", b.String())
	WriteCVFile("positions.md", b.String())
	if len(langs) > 0 {
		MakeCVDir("position")
	}
	for name, item := range positionsMap {
		if len(item.SubjectMatters) == 0 {
			panic("no subject matter on position " + name)
		}
		WriteCVFile("position/"+withoutSpaces(name)+".md", string(item.Bytes()))
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
		compLink := linkFor(client.Name, "cv", "company", withoutSpaces(client.Name))
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
		WriteCVFile("client/"+withoutSpaces(name)+".md", string(item.Bytes()))
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
		WriteCVFile("company/"+withoutSpaces(name)+".md", string(company.Bytes()))
	}
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
		if len(platform.SubjectMatters) == 0 {
			panic("no subject matter on platform " + name)
		}
		WriteCVFile("platform/"+withoutSpaces(name)+".md", string(platform.Bytes(name, "Platform")))
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
		WriteCVFile("db/"+withoutSpaces(name)+".md", string(db.Bytes(name, "Database")))
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
		WriteCVFile("cache/"+withoutSpaces(name)+".md", string(cache.Bytes(name, "Cache")))
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
		WriteCVFile("provider/"+withoutSpaces(name)+".md", string(provider.Bytes()))
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
		WriteCVFile("service/"+withoutSpaces(name)+".md", string(service.Bytes(name, "Cloud Service")))
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
		if len(tech.SubjectMatters) == 0 {
			panic("no subject matter on technology " + name)
		}
		WriteCVFile("technology/"+withoutSpaces(name)+".md", string(tech.Bytes(name, "Technology")))
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
		WriteCVFile("miscSkill/"+withoutSpaces(name)+".md", string(skill.Bytes(name, "Misc Skill")))
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
		WriteCVFile("interest/"+withoutSpaces(name)+".md", string(interest.Bytes(name, "Interest")))
	}
}
func createSubjectMatterPages() {
	// TODO; THIS!!!!
	b := strings.Builder{}
	b.WriteString(frontmatterFor("Subject Matters"))
	alphabetizedSms := slices.Collect(maps.Keys(subjectMatters))
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
		WriteCVFile("subjectMatter/"+withoutSpaces(string(sm))+".md", string(sm.Bytes()))
	}
}
