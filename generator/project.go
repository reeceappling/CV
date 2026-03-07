package main

import (
	"appli.ng/cv/generator/utils"
	"fmt"
	"maps"
	"slices"
	"strings"
)

var projects = map[string]*Project{}

//func projectsList(projectType *projectType) []*Project {
//	out := []*Project{}
//	for _, p := range projects {
//		if projectType == nil || p.Type == *projectType {
//			out = append(out, p)
//		}
//	}
//	return out
//}

//func personalProjects() []*Project {
//	p := true
//	return projectsList(&p)
//}
//func professionalProjects() []*Project {
//	p := true
//	return projectsList(&p)
//}
//func courseworkProjects() []*Project {
//	p := true
//	return projectsList(&p)
//}

type projectType string

const (
	projectTypeSchool       projectType = "coursework"
	projectTypeProfessional projectType = "professional"
	projectTypePersonal     projectType = "personal"
)

func (pt projectType) isPersonal() bool {
	return pt == projectTypePersonal
}
func (pt projectType) isProfessional() bool {
	return pt == projectTypeProfessional
}
func (pt projectType) isCoursework() bool {
	return pt == projectTypeSchool

}

type projectTypeInfo interface {
	Type() projectType
	Finalize(*Project)
	getClient() *Client
	setClient(*Client) projectTypeInfo
	getSchool() *SchoolPage
	setSchool(page *SchoolPage) projectTypeInfo
}
type schoolProjectTypeInfo struct {
	school *SchoolPage
}

func (schoolProjectTypeInfo) Type() projectType {
	return projectTypeSchool
}
func (i schoolProjectTypeInfo) Finalize(pr *Project) {
	i.school.withProjects(pr)
}
func (i schoolProjectTypeInfo) getClient() *Client {
	return nil
}
func (i schoolProjectTypeInfo) getSchool() *SchoolPage {
	return i.school
}
func (i schoolProjectTypeInfo) setClient(cli *Client) projectTypeInfo {
	return i
}
func (i schoolProjectTypeInfo) setSchool(s *SchoolPage) projectTypeInfo {
	i.school = s
	return i
}

type professionalProjectTypeInfo struct {
	client *Client
}

func (professionalProjectTypeInfo) Type() projectType {
	return projectTypeProfessional
}
func (i professionalProjectTypeInfo) Finalize(pr *Project) {
	i.client.WithProjects(pr)
}
func (i professionalProjectTypeInfo) getClient() *Client {
	return i.client
}
func (i professionalProjectTypeInfo) getSchool() *SchoolPage {
	return nil
}
func (i professionalProjectTypeInfo) setClient(cli *Client) projectTypeInfo {
	i.client = cli
	return i
}
func (i professionalProjectTypeInfo) setSchool(s *SchoolPage) projectTypeInfo {
	return i
}

type personalProjectTypeInfo struct{}

func (i personalProjectTypeInfo) getClient() *Client {
	return nil
}
func (i personalProjectTypeInfo) getSchool() *SchoolPage {
	return nil
}
func (i personalProjectTypeInfo) setClient(cli *Client) projectTypeInfo {
	return i
}
func (i personalProjectTypeInfo) setSchool(s *SchoolPage) projectTypeInfo {
	return i
}

func (personalProjectTypeInfo) Type() projectType {
	return projectTypePersonal
}
func (personalProjectTypeInfo) Finalize(*Project) {}

type projectStatus string

const (
	statusComplete    projectStatus = "Complete"
	statusBuilding    projectStatus = "Building"
	statusMaintaining projectStatus = "Maintaining"
	statusNotOnTeam   projectStatus = "No longer part of project"
	statusShelved     projectStatus = "Shelved or not actively working on this project" // TODO: change
)

type Project struct {
	Name     string
	Summary  string
	Status   projectStatus // TODO: USE THESE!!!!
	TypeInfo projectTypeInfo

	link           *string
	Languages      map[string]Frequency
	Dbs            utils.Set[string]
	Caches         utils.Set[string]
	CloudProviders map[string]utils.Set[string]
	// TODO: CI/CD? IAC?
	Technologies utils.Set[string]
	Platforms    utils.Set[string]
	MiscSkills   utils.Set[string]
}

func (pg *Project) NameValue() string {
	return pg.Name
}

func (pg *Project) Link() string {
	if pg == nil {
		return "NO_LINK"
	}
	return fmt.Sprintf("[%s](project/%s)", pg.Name, withoutSpaces(pg.Name))
}

func (pg *Project) WithStatus(status projectStatus) *Project {
	if pg == nil {
		return pg
	}
	pg.Status = status
	return pg
}

func languagesFor(projs []*Project) map[string]Frequency {
	out := map[string]Frequency{}
	for _, proj := range projs {
		for lang, freq := range proj.Languages {
			existingFreq, exists := out[lang]
			if !exists {
				out[lang] = freq
			} else {
				if freq > existingFreq {
					out[lang] = freq
				}
			}
		}
	}
	return out
}

func (pg *Project) Company() *CompanyPage {
	cli := pg.TypeInfo.getClient()
	if cli == nil {
		return nil
	}
	return cli.company
}

func (pg *Project) Bytes() []byte {
	builder := strings.Builder{}
	// PERSONAL/Professional
	t := pg.TypeInfo.Type()
	switch t {
	case projectTypeSchool:
		builder.WriteString("Coursework-related Project\n") // TODO: ???????
	case projectTypeProfessional:
		builder.WriteString("Professional Project\n") // TODO: ???????
	case projectTypePersonal:
		builder.WriteString("Personal Project\n") // TODO: ???????
	default:
		panic("unknown project type: " + string(t))
	}
	// Client/company
	cli := pg.TypeInfo.getClient()
	if cli != nil {
		builder.WriteString(fmt.Sprintf("Client: %s\n", cli.Link()))
	}
	if comp := pg.Company(); comp != nil {
		builder.WriteString(fmt.Sprintf("Company: %s\n", comp.Link()))
	}
	// SUMMARY
	builder.WriteString("# Summary\n")
	builder.WriteString(string("Status: " + pg.Status + "\n"))
	builder.WriteString(pg.Summary + "\n")
	// LINK
	if pg.link != nil {
		builder.WriteString(fmt.Sprintf("[Source Code](%s)\n", *pg.link))
	}

	// LANGUAGES
	ls := pg.Languages
	if len(ls) > 0 {
		builder.WriteString("# Languages\n")
		builder.WriteString("Language | Usage Frequency[[" + withoutSpaces(pg.Name) + "#^usageFrequencies\\|*]]\n")
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
				builder.WriteString(fmt.Sprintf("[%s](language/%s) | %s\n", name, withoutSpaces(name), Frequency(f).String()))
			}
		}
	}
	// DBs
	if len(pg.Dbs) > 0 {
		builder.WriteString("# Databases\n")
		for db, _ := range pg.Dbs {
			builder.WriteString(fmt.Sprintf("- [%s](db/%s)\n", db, withoutSpaces(db)))
		}
	}
	// CACHES
	if len(pg.Caches) > 0 {
		builder.WriteString("# Caches\n")
		for name, _ := range pg.Caches {
			builder.WriteString(fmt.Sprintf("- [%s](cache/%s)\n", name, withoutSpaces(name)))
		}
	}
	// Cloud Providers!
	if len(pg.CloudProviders) > 0 {
		builder.WriteString("# CloudProviders\n")
		for provName, services := range pg.CloudProviders {
			svcStrings := make([]string, len(services))
			for i, svc := range slices.Collect(maps.Keys(services)) {
				svcStrings[i] = fmt.Sprintf("[%s](service/%s)", svc, svc)
			}
			builder.WriteString(fmt.Sprintf("- [%s](provider/%s)\n", provName, withoutSpaces(provName))) // TODO: FIX ME!
			builder.WriteString(strings.Join(svcStrings, ", ") + "\n")
		}
	}
	// Technologies
	if len(pg.Technologies) > 0 {
		builder.WriteString("# Technologies\n")
		for name, _ := range pg.Technologies {
			builder.WriteString(fmt.Sprintf("- [%s](technology/%s)\n", name, withoutSpaces(name))) // TODO: FIX ME!
		}
	}
	// Platforms
	if len(pg.Platforms) > 0 {
		builder.WriteString("# Platforms\n")
		for name, _ := range pg.Platforms {
			builder.WriteString(fmt.Sprintf("- [%s](platform/%s)\n", name, withoutSpaces(name))) // TODO: FIX ME!
		}
	}
	// Misc skills
	if len(pg.MiscSkills) > 0 {
		builder.WriteString("# Misc Skills\n")
		for skill, _ := range pg.MiscSkills {
			builder.WriteString(fmt.Sprintf("- [%s](miscSkill/%s)\n", skill, withoutSpaces(skill))) // TODO: FIX ME!
		}
	}
	builder.WriteString(definitionsArea())
	builder.WriteString(usageFrequencyDefinitions())
	return []byte(builder.String())
}

func (p *Project) finalize() *Project {
	p.TypeInfo.Finalize(p)
	return p
}

func (p *Project) WithLang(lang string, freq Frequency) *Project {
	l, exists := langs[lang]
	if !exists {
		l = NewLanguage(lang)
		langs[lang] = l
	}
	// Assert not already existing
	if _, exists = p.Languages[lang]; exists {
		panic(lang + " already exists on proj " + p.Name)
	}
	p.Languages[lang] = freq
	l.AddProject(p)
	return p
}
func (p *Project) WithMiscSkills(skills ...string) *Project {
	for _, skill := range skills {
		s, exists := miscSkills[skill]
		if !exists {
			s = NewSkill(skill)
			miscSkills[skill] = s
		}
		// Assert not already existing
		if _, exists = p.MiscSkills[skill]; exists {
			panic(skill + " already exists on proj " + p.Name)
		}

		s.AddProject(p)
	}
	return p
}

func (p *Project) WithCloudProvider(name string, services ...string) *Project {
	prov, exists := providers[name]
	if !exists {
		prov = NewCloudProvider(name, services...)
	}
	svcsList := prov.AddServicesByName(services...)
	// Assert not already existing
	if _, exists = p.CloudProviders[name]; !exists {
		p.CloudProviders[name] = utils.SetFrom(services...)
	}
	for _, svc := range services {
		p.CloudProviders[name].Add(svc)
	}
	// Add project on provider
	prov.AddProject(p)
	// Add project on cloud services
	for _, svc := range svcsList {
		svc.AddProject(p)
	}
	providers[name] = prov
	return p
}
func (p *Project) WithDbs(names ...string) *Project {
	for _, dbName := range names {
		if !p.Dbs.Contains(dbName) {
			p.Dbs.Add(dbName)
		}
		db, exists := dbs[dbName]
		if !exists {
			db = NewDb(dbName)
		}
		db.AddProject(p)
	}
	return p
}
func (p *Project) WithCaches(names ...string) *Project {
	for _, name := range names {
		if _, exists := p.Caches[name]; !exists {
			p.Caches.Add(name)
		}
		temp, exists := caches[name]
		if !exists {
			temp = NewCache(name)
		}
		temp.AddProject(p)
	}
	return p
}
func (p *Project) WithTechnologies(names ...string) *Project {
	for _, name := range names {
		if _, exists := p.Technologies[name]; !exists {
			p.Technologies.Add(name)
		}
		temp, exists := techs[name]
		if !exists {
			temp = NewTechnology(name)
		}
		temp.AddProject(p)
	}
	return p
}
func (p *Project) WithPlatforms(names ...string) *Project {
	for _, name := range names {
		if _, exists := p.Platforms[name]; !exists {
			p.Platforms.Add(name)
		}
		temp, exists := platforms[name]
		if !exists {
			temp = NewPlatform(name)
		}
		temp.AddProject(p)
	}
	return p
}

func (p *Project) GetAllLowest() (outCaches map[string]*CachePage, outDbs map[string]*DbPage, outLanguages map[string]Frequency, outPlatforms map[string]*PlatformPage, outProviders map[string]*CloudProviderPage, outServices map[string]map[string]*CloudServicePage, outTechnologies map[string]*TechologyPage) {
	outCaches = map[string]*CachePage{}
	for c, _ := range p.Caches {
		outCaches[c] = caches[c]
	}
	outDbs = map[string]*DbPage{}
	for v, _ := range p.Dbs {
		outDbs[v] = dbs[v]
	}
	outLanguages = map[string]Frequency{}
	for v, freq := range p.Languages {
		outLanguages[v] = freq
	}
	outPlatforms = map[string]*PlatformPage{}
	for v, _ := range p.Platforms {
		outPlatforms[v] = platforms[v]
	}
	outProviders = map[string]*CloudProviderPage{}
	for v, _ := range p.CloudProviders {
		outProviders[v] = providers[v]
	}
	outServices = map[string]map[string]*CloudServicePage{}
	for pr, svcs := range p.CloudProviders {
		servs := map[string]*CloudServicePage{}
		for _, servName := range svcs.AsSlice() {
			servs[servName] = cloudServices[servName]
		}
		outServices[pr] = servs
	}
	outTechnologies = map[string]*TechologyPage{}
	for v, _ := range p.Technologies {
		outTechnologies[v] = techs[v]
	}
	return
}

func NewPersonalProject(name string, summary string, link *string) *Project {
	return newProject(name, summary, personalProjectTypeInfo{}, link)
}

func NewProfessionalProject(name string, summary string, client *Client, link *string) *Project {
	out := newProject(name, summary, client.ProjectTypeInfo(), link)
	client.WithProjects(out)
	return out
}

func NewSchoolProject(name string, summary string, school *SchoolPage, link *string) *Project {
	out := newProject(name, summary, school.ProjectTypeInfo(), link)
	school.withProjects(out)
	return out
}

func newProject(name string, summary string, info projectTypeInfo, link *string) *Project {
	out := &Project{
		Name:           name,
		Summary:        summary,
		link:           link,
		TypeInfo:       info,
		Languages:      map[string]Frequency{},
		Dbs:            utils.Set[string]{},
		Caches:         map[string]struct{}{},
		CloudProviders: map[string]utils.Set[string]{},
		Technologies:   utils.Set[string]{},
		Platforms:      utils.Set[string]{},
	}
	if _, exists := projects[name]; exists {
		panic("project already exists")
	}
	projects[name] = out
	return out
}

var (
	MastersDataAnalysisProject, capstoneProject, cherenkovProject, roboticsTeamProject, cncLaserCutterProject, aerospaceFinalProject, coreShufflerProject, WellAwareProject, CritColaProject, CharityProject, nfcScannerProject, measurementsProject, mushDbProject, cvProject, polygonBuilderProject, tileGenProject, ogreProject, renderProject, statsProject, billingProject, explorerProject, wqdbProject, supportProject, scudsProject, ufoProject, goweProject, saiCollegeProject, GhaRunnersProject *Project
)
var measurementsUrl = "github.com/reeceappling/measurements" // TODO: ensure ok
var mushDbUrl = "github.com/reeceappling/mushDb"             // TODO: ensure ok
var cvUrl = "github.com/reeceappling/cv"                     // TODO: ensure ok
var nfcScannerUrl = "github.com/reeceappling/nfcScanner"     // TODO: ensure ok
var coreShufflerUrl = "github.com/reeceappling/coreShuffler" // TODO: ensure ok

// initProjectsAfterClients must be called after clients and schools have been initially created
func initProjectsAfterClients() {
	polygonBuilderProject = NewProfessionalProject("Polygon Builder", fixmeLink, jdClient, nil)
	tileGenProject = NewProfessionalProject("Tile Generator", fixmeLink, jdClient, nil)
	GhaRunnersProject = NewProfessionalProject("Github Actions GPU Runners", "FIX SUMMARY", jdClient, nil)
	ogreProject = NewProfessionalProject("Organizational Geospatial Rollup Engine", fixmeLink, jdClient, nil)
	renderProject = NewProfessionalProject("Render", fixmeLink, jdClient, nil) // TODO: MORE!
	statsProject = NewProfessionalProject("Statistics", fixmeLink, jdClient, nil)
	billingProject = NewProfessionalProject("Billing", fixmeLink, jdClient, nil)
	explorerProject = NewProfessionalProject("Transform Explorer", fixmeLink, jdClient, nil)
	wqdbProject = NewProfessionalProject("Work Queue Database", fixmeLink, jdClient, nil)
	supportProject = NewProfessionalProject("Support Api", fixmeLink, jdClient, nil)
	scudsProject = NewProfessionalProject("Scuds Api", fixmeLink, jdClient, nil)
	ufoProject = NewProfessionalProject("UFO API", fixmeLink, jdClient, nil)
	goweProject = NewProfessionalProject("Gowe Builder", fixmeLink, jdClient, nil)
	saiCollegeProject = NewProfessionalProject("SimpsonUniv", fixmeLink, sourceAlliesClient, nil) // TODO: MORE!
	mushDbProject = NewPersonalProject("MushDb", fixmeLink, &mushDbUrl)
	cvProject = NewPersonalProject("CV", "This project! A generator which creates markdown files that can be viewed via Obsidian, or published to the web.", &cvUrl)
	measurementsProject = NewPersonalProject("Measurements", fixmeLink, &measurementsUrl)
	nfcScannerProject = NewPersonalProject("Nfc Scanner", fixmeLink, &nfcScannerUrl)
	coreShufflerProject = NewPersonalProject("Simulate Core Shuffler", fixmeLink, &coreShufflerUrl)
	CharityProject = NewProfessionalProject("Charity Site", fixmeLink, charityClient, nil)
	WellAwareProject = NewProfessionalProject("Well Aware NC", fixmeLink, clarkClient, nil) // TODO: change client to the lab???
	CritColaProject = NewProfessionalProject("CritCola", fixmeLink, critColaClient, nil)
	MastersDataAnalysisProject = NewProfessionalProject("Masters Data Analysis", fixmeLink, clarkClient, nil)
	capstoneProject = NewSchoolProject("Capstone Project. Uranium Silicide Accident-Tolerant Fuel cycle design for Duke Energy Catawba Nuclear Plant", "FIX M_E", schoolNCSU, &coreShufflerUrl)
	cherenkovProject = NewSchoolProject("Cherenkov radiation detector", fixmeLink, schoolNCSU, nil)
	roboticsTeamProject = NewSchoolProject("Robotics team 3720", fixmeLink, schoolCata, nil)
	cncLaserCutterProject = NewSchoolProject("CNC Laser Cutter", fixmeLink, schoolCata, nil)
	aerospaceFinalProject = NewSchoolProject("Aerospace senior design course", "Designed, created, and tested a rocket from scratch", schoolCata, nil)
}

func initProjectsFinal() {
	polygonBuilderProject = polygonBuilderProject.WithStatus(statusMaintaining).
		WithLang("Go", Extensively).
		WithLang("Scala", Often).
		WithLang("Terraform", Regularly).
		WithLang("Docker", Regularly).
		WithLang("Java", Some).
		WithLang("Bash", Some).
		WithLang("C", Some).
		WithLang("Cpp", Some).
		WithLang("SQL", Some).
		WithLang("Javascript", Rarely).
		WithLang("Python", Rarely).
		WithLang("Rust", Minimal).
		WithCloudProvider("AWS",
			"ECS", "S3", "EC2", "IAM", "SecretsManager", "DynamoDB", "DAX", "Cloudwatch", "Lambda", "ECR", "Route53", "Kinesis", "SQS", "SNS", "ApiGateway", // TODO: FARGATE NOT BEING ON THIS LIST IS CAUSING PROBLEMS
			// TODO: complete list (all of them)
		).WithDbs("Aurora", "DynamoDB", "Postgres").
		WithCaches("Redis", "memcached").
		WithTechnologies("CUDA", "Github Actions").                     // TODO: NO CUDA ON POLYGONS
		WithPlatforms("Datadog", "Logcentral", "Github", "ServiceNow"). // TODO: more?
		finalize()
	tileGenProject = tileGenProject.WithStatus(statusMaintaining).
		WithLang("Go", Extensively).
		WithLang("Scala", Often).
		WithLang("Terraform", Regularly).
		WithLang("Docker", Regularly).
		WithLang("Bash", Some).
		WithCloudProvider("AWS",
			"ECS", "S3", "EC2", "IAM", "SecretsManager", "DynamoDB", "DAX", "Cloudwatch", "Lambda", "ECR", "Route53", "Kinesis", "SQS", "SNS", "ApiGateway",
			// TODO: complete list (all of them)
		).WithDbs("DynamoDB", "Postgres").
		WithCaches("Redis", "memcached").
		WithTechnologies("CUDA", "Github Actions").                     // TODO: NO CUDA ON POLYGONS
		WithPlatforms("Datadog", "Logcentral", "Github", "ServiceNow"). // TODO: more?
		finalize()
	GhaRunnersProject = GhaRunnersProject.WithStatus(statusMaintaining).
		WithLang("Terraform", Regularly).
		WithLang("Docker", Regularly).
		WithLang("Bash", Some).
		WithCloudProvider("AWS",
			"ECS", "S3", "EC2", "IAM", "SecretsManager", "Cloudwatch", "Lambda", "ECR", "Route53",
		).
		WithTechnologies("Github Actions").
		finalize()
	ogreProject = ogreProject.WithStatus(statusBuilding).
		WithLang("Go", Extensively).
		WithLang("Terraform", Regularly).
		WithLang("Docker", Regularly).
		WithLang("Bash", Some).
		WithLang("SQL", Some).
		WithCloudProvider("AWS",
			"ECS", "S3", "EC2", "IAM", "SecretsManager", "DynamoDB", "DAX", "Cloudwatch", "Lambda", "ECR", "Route53", "Kinesis", "SQS", "SNS",
			// TODO: complete list (all of them)
		).WithDbs("Aurora", "DynamoDB", "Postgres").
		WithCaches("Redis", "memcached").
		WithTechnologies("Github Actions"). // TODO: ?????
		WithPlatforms("Datadog", "Logcentral", "Github", "ServiceNow").
		finalize()
	renderProject = renderProject.WithStatus(statusMaintaining).
		WithLang("Go", Extensively).
		WithLang("Terraform", Often).  // TODO: ADD IAC TAG
		WithLang("Docker", Regularly). // TODO: add containerization tag
		WithLang("Bash", Some).
		WithLang("SQL", Some).
		WithTechnologies("Github Actions").
		finalize()
	statsProject = statsProject.WithStatus(statusMaintaining).
		WithLang("Go", Extensively).
		WithLang("Terraform", Often).
		WithLang("Docker", Regularly).
		WithLang("Bash", Some).
		WithLang("C", Often).
		WithLang("Cpp", Often).
		WithLang("SQL", Some).
		WithLang("Javascript", Rarely).
		WithLang("Python", Rarely).
		WithLang("Rust", Minimal).
		WithCloudProvider("AWS",
			"ECS", "S3", "EC2", "IAM", "SecretsManager", "DynamoDB", "DAX", "Cloudwatch", "Lambda", "ECR", "Route53", "Kinesis", "SQS", "SNS",
			// TODO: complete list (all of them)
		).WithDbs("Aurora", "DynamoDB", "Postgres").
		WithCaches("Redis", "memcached").
		WithTechnologies("CUDA", "Github Actions").
		WithPlatforms("Datadog", "Logcentral", "Github", "ServiceNow"). // TODO: MORE!
		finalize()
	billingProject = billingProject.WithStatus(statusBuilding).
		WithLang("Go", Extensively).
		WithLang("Terraform", Often).
		WithLang("Docker", Regularly).
		WithLang("Bash", Some).
		WithLang("SQL", Some).
		WithLang("Javascript", Rarely).
		WithCloudProvider("AWS",
			"ECS", "S3", "EC2", "IAM", "SecretsManager", "DynamoDB", "DAX", "Cloudwatch", "Lambda", "ECR", "Route53", "Kinesis", "SQS", "SNS",
			// TODO: complete list (all of them)
		).WithDbs("Aurora", "DynamoDB", "Postgres").
		WithCaches("Redis", "memcached").
		WithTechnologies("Github Actions").
		WithPlatforms("Datadog", "Logcentral", "Github", "ServiceNow"). // TODO: MORE!
		finalize()                                                      // TODO: MORE!
	explorerProject = explorerProject.WithStatus(statusMaintaining).
		WithLang("Scala", Often).
		WithLang("Javascript", Often).
		WithLang("Terraform", Often).
		WithLang("Docker", Regularly).
		WithLang("Bash", Some).
		WithCloudProvider("AWS",
			"ECS", "Fargate", "S3", "EC2", "IAM", "SecretsManager", "DynamoDB", "DAX", "Cloudwatch", "Lambda", "ECR", "Route53", "Kinesis", "SQS", "SNS",
			// TODO: complete list (all of them)
		).WithDbs("Aurora", "DynamoDB", "Postgres").
		WithCaches("Redis", "memcached").
		WithTechnologies("Github Actions").
		WithPlatforms("Datadog", "Logcentral", "Github", "ServiceNow"). // TODO: MORE!
		finalize()
	wqdbProject = wqdbProject.WithStatus(statusComplete).
		WithLang("SQL", Often).
		WithLang("Go", Often).
		WithLang("Scala", Often).
		WithLang("Terraform", Often).
		WithLang("Docker", Regularly).
		WithLang("Bash", Some).
		WithCloudProvider("AWS",
			"ECS", "Fargate", "S3", "EC2", "IAM", "SecretsManager", "DynamoDB", "DAX", "Cloudwatch", "Lambda", "ECR", "Route53", "Kinesis", "SQS", "SNS",
			// TODO: complete list (all of them)
		).WithDbs("Aurora", "Postgres").
		WithCaches().
		WithTechnologies("Github Actions").
		WithPlatforms("Datadog", "Logcentral", "Github", "ServiceNow"). // TODO: MORE!
		finalize()                                                      // TODO: MORE!
	supportProject = supportProject.WithStatus(statusComplete).
		WithLang("SQL", Regularly).
		WithLang("Go", Often).
		WithLang("Terraform", Often).
		WithLang("Docker", Regularly).
		WithLang("Bash", Some).
		WithCloudProvider("AWS",
			"ECS", "Fargate", "S3", "EC2", "IAM", "SecretsManager", "DynamoDB", "DAX", "Cloudwatch", "Lambda", "ECR", "Route53", "Kinesis", "SQS", "SNS",
			// TODO: complete list (all of them)
		).WithDbs().
		WithCaches().
		WithTechnologies("Github Actions").
		WithPlatforms("Datadog", "Logcentral", "Github", "ServiceNow"). // TODO: MORE!
		finalize()
	scudsProject = scudsProject.WithStatus(statusComplete).
		WithLang("Scala", Extensively).
		WithLang("Go", Often).
		WithLang("Terraform", Often).
		WithLang("Docker", Regularly).
		WithLang("Bash", Some).
		// TODO: what others?
		WithCloudProvider("AWS",
			"ECS", "Fargate", "S3", "EC2", "IAM", "SecretsManager", "DynamoDB", "DAX", "Cloudwatch", "Lambda", "ECR", "Route53", "Kinesis", "SQS", "SNS",
			// TODO: complete list (all of them)
		).WithDbs().
		WithCaches("Redis").
		WithTechnologies("ElasticSearch", "Github Actions").
		WithPlatforms("Datadog", "Logcentral", "Github", "ServiceNow"). // TODO: MORE!
		finalize()
	ufoProject = ufoProject.WithStatus(statusMaintaining).
		WithLang("Scala", Extensively).
		WithLang("Go", Regularly).
		WithLang("Terraform", Regularly).
		WithLang("Docker", Regularly).
		WithLang("Bash", Some).
		// TODO: what others?
		WithCloudProvider("AWS",
			"ECS", "Fargate", "S3", "EC2", "IAM", "SecretsManager", "DynamoDB", "DAX", "Cloudwatch", "Lambda", "ECR", "Route53", "Kinesis", "SQS", "SNS",
			// TODO: complete list (all of them)
		).WithDbs().
		WithCaches("Redis").
		WithTechnologies("Github Actions").
		WithPlatforms("Datadog", "Logcentral", "Github", "ServiceNow"). // TODO: MORE!
		finalize()
	goweProject = goweProject.WithStatus(statusMaintaining).
		WithLang("Go", Extensively).
		WithLang("Terraform", Regularly).
		WithLang("Bash", Some).
		WithLang("Docker", Regularly).
		WithCloudProvider("AWS",
			"ECS", "Fargate", "S3", "EC2", "IAM", "SecretsManager", "DynamoDB", "DAX", "Cloudwatch", "Lambda", "ECR", "Route53", "Kinesis", "SQS", "SNS",
			// TODO: complete list (all of them)
		).WithDbs().
		WithCaches("Redis").
		WithTechnologies("Github Actions").
		WithPlatforms("Datadog", "Logcentral", "Github", "ServiceNow"). // TODO: MORE!
		finalize()
	saiCollegeProject = saiCollegeProject.WithStatus(statusComplete). // TODO: MORE!
										finalize()
	// TODO: sai project for JAMF
	mushDbProject = mushDbProject.WithStatus(statusBuilding).
		WithLang("Go", Extensively).
		WithLang("Typescript", Extensively).
		WithLang("Javascript", Extensively).
		WithLang("Docker", Often).
		WithLang("Docker Compose", Often).
		WithLang("Kubernetes", Some).
		WithDbs("mongodb").
		WithCaches().
		WithTechnologies("Github Actions", "React", "NextJs", "RFID", "NFC", "I2C", "Cloudflare", "Grafana", "Github Actions"). // TODO: MORE?
		WithPlatforms().
		WithCloudProvider("GCP", "Google Cloud DNS", "Identity Platform"). // TODO: MORE!
		finalize()
	cvProject = cvProject.WithStatus(statusBuilding).
		WithLang("Go", Extensively).
		WithLang("Javascript", Some).
		WithDbs().
		WithCaches().
		WithTechnologies("Github Actions", "Obsidian", "Markdown", "Quartz"). // TODO: MORE?
		WithPlatforms().                                                      // TODO: MORE!
		finalize()
	measurementsProject = measurementsProject.WithStatus(statusShelved).
		WithLang("Go", Extensively).
		WithDbs().
		WithCaches().
		WithTechnologies("Github Actions"). // TODO: MORE?
		WithPlatforms().                    // TODO: MORE!
		finalize()
	nfcScannerProject = nfcScannerProject.WithStatus(statusBuilding).
		WithLang("Go", Extensively).
		WithDbs().
		WithCaches().
		WithTechnologies("Github Actions"). // TODO: MORE?
		WithPlatforms().                    // TODO: MORE!
		finalize()
	coreShufflerProject = coreShufflerProject.WithStatus(statusComplete).
		WithLang("Javascript", Extensively).
		WithDbs().
		WithCaches().
		WithTechnologies("JQuery"). // TODO: MORE?
		WithPlatforms().            // TODO: MORE!
		finalize()
	CharityProject = CharityProject.WithStatus(statusComplete). // TODO: MORE!
									finalize()
	WellAwareProject = WellAwareProject.WithStatus(statusComplete).
		WithCloudProvider("AWS",
						"S3", "IAM", "Route53").
		WithTechnologies("Gitlab CI"). // TODO: MORE!
		finalize()
	CritColaProject = CritColaProject.WithStatus(statusComplete).
		WithCloudProvider("AWS",
			"S3", "EC2", "IAM", "SecretsManager", "ECR", "Route53").
		WithLang("Docker", Regularly).
		WithLang("Terraform", Often).
		WithTechnologies("Gitlab CI"). // TODO: MORE!
		finalize()
	MastersDataAnalysisProject = MastersDataAnalysisProject.WithStatus(statusComplete).
		WithLang("Python", Often).
		WithTechnologies("Gitlab CI"). // TODO: MORE!
		finalize()
	capstoneProject = capstoneProject.WithStatus(statusComplete).
		WithLang("Javascript", Some). // TODO: any others? CASMO, SIMULATE? FORTRAN
		WithLang("Fortran", Often).
		WithDbs().
		WithCaches().
		WithTechnologies("JQuery", "Git"). // TODO: MORE?
		WithPlatforms("Github").           // TODO: MORE!
		finalize()
	cherenkovProject = cherenkovProject.WithStatus(statusComplete).
		WithLang("Python", Often). // TODO: any others?
		WithDbs().
		WithCaches().
		WithTechnologies("Arduino"). // TODO: MORE?
		WithPlatforms("Github").     // TODO: MORE!
		finalize()
	roboticsTeamProject = roboticsTeamProject.WithStatus(statusComplete).
		WithLang("Python", Often). // TODO: any others?
		WithDbs().
		WithCaches().
		WithTechnologies("Arduino"). // TODO: MORE?
		WithPlatforms().             // TODO: MORE!
		finalize()
	cncLaserCutterProject = cncLaserCutterProject.WithStatus(statusComplete).
		WithLang("Java", Often). // TODO: any others?
		WithLang("Javascript", Regularly).
		WithLang("HTML", Regularly).
		WithDbs().
		WithCaches().
		WithTechnologies("Arduino"). // TODO: MORE?
		WithPlatforms().             // TODO: MORE!
		finalize()
	aerospaceFinalProject = aerospaceFinalProject.WithStatus(statusComplete).
		WithLang("Java", Minimal). // TODO: any others?
		WithLang("Javascript", Minimal).
		WithDbs().
		WithCaches().
		WithTechnologies("Arduino"). // TODO: MORE?
		WithPlatforms().             // TODO: MORE!
		finalize()
}
