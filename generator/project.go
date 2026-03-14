package main

import (
	"appli.ng/cv/generator/utils"
	"fmt"
	"maps"
	"slices"
	"strings"
)

var projects = map[string]*Project{}

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
	Technologies   utils.Set[string]
	Platforms      utils.Set[string]
	MiscSkills     utils.Set[string]
	// TODO: MORE!
	RelatedInterests utils.Set[string]
	SubjectMatters   utils.Set[SubjectMatter]
	Tags             utils.Set[Tag]
}

func (pg *Project) NameValue() string {
	return pg.Name
}

func (pg *Project) Link() string {
	if pg == nil {
		return "NO_LINK"
	}
	return linkFor(pg.Name, "cv", "project", withoutSpaces(pg.Name))
}

func (pg *Project) WithStatus(status projectStatus) *Project {
	if pg == nil {
		return pg
	}
	pg.Status = status
	return pg
}

func (pg *Project) WithSubjectMatters(sms ...SubjectMatter) *Project {
	if pg == nil {
		return pg
	}
	for _, sm := range sms {
		if sm == smFullStack {
			pg.SubjectMatters.Add(smFrontend, smBackend)
		}
		pg.SubjectMatters.Add(sm)
	}

	return pg
}

func (pg *Project) WithSummary(summary string) *Project {
	if pg == nil {
		return pg
	}
	pg.Summary = summary
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

func (pg *Project) EntryType() string {
	return "Project"
}

func (pg *Project) Bytes() []byte {
	builder := strings.Builder{}
	builder.WriteString(frontmatterFor(pg.Name, "Project"))
	// PERSONAL/Professional
	t := pg.TypeInfo.Type()
	switch t {
	case projectTypeSchool:
		builder.WriteString("Coursework-related Project\n")
	case projectTypeProfessional:
		builder.WriteString("Professional Project\n")
	case projectTypePersonal:
		builder.WriteString("Personal Project\n")
	default:
		panic("unknown project type: " + string(t))
	}
	// Client/company/school

	if cli := pg.TypeInfo.getClient(); cli != nil {
		builder.WriteString(fmt.Sprintf("Client: %s\n", cli.Link()))
		if comp := cli.company; comp != nil {
			builder.WriteString(fmt.Sprintf("Company: %s\n", comp.Link()))
		}
	} else if sch := pg.TypeInfo.getSchool(); sch != nil {
		builder.WriteString(fmt.Sprintf("School: %s\n", sch.Link()))
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
				builder.WriteString(fmt.Sprintf("%s | %s\n", langs[name].Link(), Frequency(f).String()))
			}
		}
	}
	// DBs
	if len(pg.Dbs) > 0 {
		builder.WriteString("# Databases\n")
		for db, _ := range pg.Dbs {
			builder.WriteString(fmt.Sprintf("- %s\n", dbs[db].Link()))
		}
	}
	// CACHES
	if len(pg.Caches) > 0 {
		builder.WriteString("# Caches\n")
		for name, _ := range pg.Caches {
			builder.WriteString(fmt.Sprintf("- %s\n", caches[name].Link()))
		}
	}
	// Cloud Providers!
	if len(pg.CloudProviders) > 0 {
		builder.WriteString("# CloudProviders\n")
		for provName, services := range pg.CloudProviders {
			svcStrings := make([]string, len(services))
			for i, svc := range slices.Collect(maps.Keys(services)) {
				svcStrings[i] = cloudServices[svc].Link()
			}
			builder.WriteString(fmt.Sprintf("- %s\n", providers[provName].Link()))
			builder.WriteString(strings.Join(svcStrings, ", ") + "\n")
		}
	}
	// Technologies
	if len(pg.Technologies) > 0 {
		builder.WriteString("# Technologies\n")
		for name, _ := range pg.Technologies {
			builder.WriteString(fmt.Sprintf("- %s\n", techs[name].Link()))
		}
	}
	// Platforms
	if len(pg.Platforms) > 0 {
		builder.WriteString("# Platforms\n")
		for name, _ := range pg.Platforms {
			builder.WriteString(fmt.Sprintf("- %s\n", platforms[name].Link()))
		}
	}
	// Misc skills
	if len(pg.MiscSkills) > 0 {
		builder.WriteString("# Misc Skills\n")
		for skill, _ := range pg.MiscSkills {
			builder.WriteString(fmt.Sprintf("- %s\n", miscSkills[skill].Link()))
		}
	}
	// SubjectMatters
	if len(pg.SubjectMatters) > 0 {
		builder.WriteString("# Related Subject Matters\n")
		for sm, _ := range pg.SubjectMatters {
			builder.WriteString(fmt.Sprintf("- %s\n", sm.Link()))
		}
	}
	// Related Interests
	if len(pg.RelatedInterests) > 0 {
		builder.WriteString("# Related Interests\n")
		for interest, _ := range pg.RelatedInterests {
			builder.WriteString(fmt.Sprintf("- %s\n", linkForInterest(interest)))
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

func (p *Project) WithLang(lang string, freq Frequency, sms ...SubjectMatter) *Project {
	l, exists := langs[lang]
	if !exists {
		l = NewLanguage(lang)
		langs[lang] = l
	}
	l.WithSubjectMatters(sms...)
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
func (p *Project) WithInterests(intrsts ...string) *Project {
	for _, interest := range intrsts {
		s, exists := interests[interest]
		if !exists {
			s = NewInterest(interest)
			interests[interest] = s
		}
		// Assert not already existing
		if _, exists = p.RelatedInterests[interest]; exists {
			panic(interest + " already exists on proj " + p.Name)
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
	// TODO: Aurora, DynamoDB, DocumentDB, RDS should all point to the AWS service!!!!??? (probably not)
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
func (p *Project) WithTags(tags ...Tag) *Project {
	p.Tags.Add(tags...)
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
	if strings.Contains(name, ".") || strings.Contains(name, "%") {
		panic("invalid character in project " + name)
	}
	out := &Project{
		Name:             name,
		Summary:          summary,
		link:             link,
		TypeInfo:         info,
		Languages:        map[string]Frequency{},
		Dbs:              utils.Set[string]{},
		Caches:           map[string]struct{}{},
		CloudProviders:   map[string]utils.Set[string]{},
		Technologies:     utils.Set[string]{},
		Platforms:        utils.Set[string]{},
		SubjectMatters:   map[SubjectMatter]struct{}{},
		RelatedInterests: map[string]struct{}{},
		Tags:             utils.Set[Tag]{},
	}
	if _, exists := projects[name]; exists {
		panic("project already exists")
	}
	projects[name] = out
	addLinkable(strings.ToLower(name), out)
	return out
}

// var (
//
//	projectLinAlgCryptography, MastersDataAnalysisProject, capstoneProject, cherenkovProject, roboticsTeamProject, cncLaserCutterProject, aerospaceFinalProject, coreShufflerProject, WellAwareProject, CritColaProject, CharityProject, nfcScannerProject, measurementsProject, mushDbProject, cvProject, polygonBuilderProject, tileGenProject, ogreProject, renderProject, statsProject, billingProject, explorerProject, wqdbProject, supportProject, scudsProject, ufoProject, goweProject, simpsonUnivProject, GhaRunnersProject, projectAgentSwarm *Project
//
// )
var measurementsUrl = "github.com/reeceappling/measurements" // TODO: ensure ok
var mushDbUrl = "github.com/reeceappling/mushDb"             // TODO: ensure ok
var cvUrl = "github.com/reeceappling/CV"                     // TODO: ensure ok
var nfcScannerUrl = "github.com/reeceappling/nfcScanner"     // TODO: ensure ok
var coreShufflerUrl = "github.com/reeceappling/coreShuffler" // TODO: ensure ok

var (
	projectAgentSwarm              = NewPersonalProject("AI Agent Swarm", fixmeLink, nil)
	polygonBuilderProject          = NewProfessionalProject("Polygon Builder", fixmeLink, jdClient, nil)
	tileGenProject                 = NewProfessionalProject("Tile Generator", fixmeLink, jdClient, nil)
	GhaRunnersProject              = NewProfessionalProject("Github Actions GPU Runners", "FIX SUMMARY", jdClient, nil)
	ogreProject                    = NewProfessionalProject("Organizational Geospatial Rollup Engine", fixmeLink, jdClient, nil)
	renderProject                  = NewProfessionalProject("Render Cluster", fixmeLink, jdClient, nil) // TODO: MORE!
	statsProject                   = NewProfessionalProject("Statistics Cluster", fixmeLink, jdClient, nil)
	billingProject                 = NewProfessionalProject("Billing Cluster", fixmeLink, jdClient, nil)
	explorerProject                = NewProfessionalProject("Transform Explorer", fixmeLink, jdClient, nil)
	wqdbProject                    = NewProfessionalProject("Work Queue Database", fixmeLink, jdClient, nil)
	supportProject                 = NewProfessionalProject("Support Cluster", fixmeLink, jdClient, nil)
	scudsProject                   = NewProfessionalProject("Scuds API", fixmeLink, jdClient, nil)
	ufoProject                     = NewProfessionalProject("UFO API", fixmeLink, jdClient, nil)
	goweProject                    = NewProfessionalProject("Gowe Builder", fixmeLink, jdClient, nil)
	rasterRenderProject            = NewProfessionalProject("Raster Render", fixmeLink, jdClient, nil)
	simpsonUnivProject             = NewProfessionalProject("Simpson University", fixmeLink, simpsonUniversityClient, nil) // TODO: MORE!
	jamfProject                    = NewProfessionalProject("JAMF companywide setup", fixmeLink, sourceAlliesClient, nil)
	smallImprovementsProject       = NewProfessionalProject("Small Improvements Bot", fixmeLink, sourceAlliesClient, nil)
	internalResumeGeneratorProject = NewProfessionalProject("Source Allies Internal Consultant Resume Generator", fixmeLink, sourceAlliesClient, nil)
	mushDbProject                  = NewPersonalProject("MushDb", fixmeLink, &mushDbUrl)
	cvProject                      = NewPersonalProject("Personal Site and CV", "This project! A generator which creates markdown files that can be viewed via Obsidian, or published to the web.", &cvUrl) // TODO: make multiple strings an available option for summary
	linksPage                      = NewPersonalProject("Personal Links Page", "A page to put all my links", nil)                                                                                           // TODO: LINKS PAGE
	measurementsProject            = NewPersonalProject("Measurements", fixmeLink, &measurementsUrl)
	nfcScannerProject              = NewPersonalProject("Nfc Scanner", fixmeLink, &nfcScannerUrl)
	teiProjects                    = NewProfessionalProject("Tei Projects", fixmeLink, teiClient, nil) // TODO: maybe add an actual project
	taePhotoImporter               = NewProfessionalProject("Tae Field Photograph Importer", fixmeLink, taeClient, nil)
	mafcProjects                   = NewProfessionalProject("Mafc Projects", fixmeLink, mafcClient, nil) // TODO: maybe add actual projects?
	coreShufflerProject            = NewPersonalProject("Simulate Core Shuffler", fixmeLink, &coreShufflerUrl)
	CharityProject                 = NewProfessionalProject("Charity Site", fixmeLink, charityClient, nil)
	WellAwareProject               = NewProfessionalProject("Well Aware NC", fixmeLink, clarkClient, nil) // TODO: change client to the lab???
	CritColaProject                = NewProfessionalProject("CritCola", fixmeLink, critColaClient, nil)   // TODO: ADD OTHER CRITCOLA PROJECTS
	ArrowNailProject               = NewProfessionalProject("Hail History Tracker", fixmeLink, arrowNailClient, nil)
	WildlifeRProject               = NewProfessionalProject("Wildlife R Data Analysis", fixmeLink, wildlifeRClient, nil)
	MastersDataAnalysisProject     = NewProfessionalProject("Masters Data Analysis", fixmeLink, clarkClient, nil)
	capstoneProject                = NewSchoolProject("Capstone Project-Uranium Silicide Accident Tolerant Fuel cycle design for Duke Energy Catawba Nuclear Plant", "FIX M_E", schoolNCSU, &coreShufflerUrl)
	reactorAnalysisFinal           = NewSchoolProject("Reactor Analysis Exam", "Reactor analysis final exam code", schoolNCSU, utils.Pointer(fixmeLink)) // TODO: ADD REAL LINK!
	monteCarloProject              = NewSchoolProject("MonteCarlo scattering and decay project", "FIX_ME", schoolNCSU, utils.Pointer(fixmeLink))         // TODO: ADD REAL LINK!
	// TODO; FORTRAN SCHOOL PROJECT
	projectLinAlgCryptography = NewSchoolProject("Linear algebra cryptography algorithm", "FIX M_E", schoolNCSU, nil)
	cherenkovProject          = NewSchoolProject("Cherenkov radiation sensor", fixmeLink, schoolNCSU, nil)
	roboticsTeamProject       = NewSchoolProject("Robotics team 3720", fixmeLink, schoolCata, nil)
	cncLaserCutterProject     = NewSchoolProject("CNC Laser Cutter", fixmeLink, schoolCata, nil)
	aerospaceFinalProject     = NewSchoolProject("Aerospace senior design course", "Designed, created, and tested a rocket from scratch", schoolCata, nil)
	miscSmallPersonalProjects = NewPersonalProject("Misc small personal projects", "A conglomeration of personal projects which did not each deserve their own entry", nil)
)

func initProjectsFinal() {
	projectAgentSwarm = projectAgentSwarm.WithStatus(statusMaintaining).
		WithSummary("SUMMARY HERE"). // TODO: MORE!
		WithLang("Go", Extensively, smBackend).
		WithInterests("Agentic AI", "Agent Swarms"). // TODO: any more?
		WithPlatforms("Poolside AI").
		WithTechnologies("LangGraph", "LangChain", "LLM", "AI", "AI Agents", "OpenAI API spec").
		WithSubjectMatters(smAI, smBackend).
		//WithTags(tagAI, tagBackend).
		finalize()
	polygonBuilderProject = polygonBuilderProject.WithStatus(statusMaintaining).
		WithSummary("SUMMARY HERE"). // TODO: MORE!
		WithLang("Go", Extensively, smBackend, smScripting).
		WithLang("Scala", Often, smBackend).
		WithLang("Terraform", Regularly, smIAC).
		WithLang("Docker", Regularly, smContainerization).
		WithLang("Java", Some, smBackend).
		WithLang("Bash", Some, smScripting).
		WithLang("C", Some, smBackend).
		WithLang("Cpp", Some, smBackend).
		WithLang("SQL", Some, smBackend).
		WithLang("Javascript", Rarely, smFullStack).
		WithLang("Python", Rarely, smScripting).
		WithLang("Rust", Minimal, smScripting).
		WithCloudProvider("AWS",
			"ECS", "S3", "EC2", "IAM", "SecretsManager", "DynamoDB", "DAX", "Cloudwatch", "Lambda", "ECR", "Route53", "Kinesis", "SQS", "SNS", "ApiGateway", "EBS", "ELB", "ALB",
		).
		WithDbs("Aurora", "DynamoDB", "Postgres").
		WithCaches("Redis", "memcached").
		WithTechnologies("Github Actions", "Parquet", "Avro", "REST API").
		WithPlatforms("Datadog", "Logcentral", "Github", "ServiceNow", "Azure DevOps", "Jira", "Confluence").
		WithSubjectMatters(smApi, smGeospatial, smBackend, smTopology, smLinearAlgebra, smClusterComputing, smDistributedComputing, smContainerization, smIAC, smCiCd).
		//WithTags(tagBackend, tagTopology, tagLinearAlgebra, tagClusterComputing, tagDistributedComputing, tagContainerization, tagIAC, tagCiCd)
		finalize()
	tileGenProject = tileGenProject.WithStatus(statusMaintaining).
		WithSummary("SUMMARY HERE"). // TODO: MORE!
		WithLang("Go", Extensively).
		WithLang("Scala", Often).
		WithLang("Terraform", Regularly).
		WithLang("Docker", Regularly).
		WithLang("Bash", Some).
		WithCloudProvider("AWS",
			"ECS", "S3", "EC2", "IAM", "SecretsManager", "DynamoDB", "DAX", "Cloudwatch", "Lambda", "ECR", "Route53", "Kinesis", "SQS", "SNS", "ApiGateway",
		).
		WithDbs("DynamoDB", "Postgres").
		WithCaches("Redis", "memcached", "DAX").
		WithTechnologies("CUDA", "Github Actions", "Parquet", "Avro", "LocalStack", "REST API").
		WithPlatforms("Datadog", "Logcentral", "Github", "ServiceNow", "Azure DevOps", "Jira", "Confluence").
		WithSubjectMatters(smApi, smGeospatial, smBackend, smLinearAlgebra, smClusterComputing, smDistributedComputing, smContainerization, smIAC, smCiCd).
		finalize()
	GhaRunnersProject = GhaRunnersProject.WithStatus(statusMaintaining).
		WithSummary("SUMMARY HERE"). // TODO: MORE!
		WithLang("Terraform", Regularly).
		WithLang("Docker", Regularly).
		WithLang("Bash", Some).
		WithCloudProvider("AWS",
			"ECS", "S3", "EC2", "IAM", "SecretsManager", "Cloudwatch", "Lambda", "ECR", "Route53", // Route53 add networking sm
		).
		WithTechnologies("Github Actions").
		WithPlatforms("Datadog", "Logcentral", "Github", "Azure DevOps", "Jira", "Confluence", "Teams").
		WithSubjectMatters(smBackend, smCiCd, smDevOps, smContainerization, smIAC, smCiCd).
		finalize()
	ogreProject = ogreProject.WithStatus(statusBuilding).
		WithSummary("SUMMARY HERE"). // TODO: MORE!
		WithLang("Go", Extensively).
		WithLang("Terraform", Regularly).
		WithLang("Docker", Regularly).
		WithLang("Bash", Some).
		WithLang("SQL", Some).
		WithDbs("Aurora", "DynamoDB", "Postgres").
		WithCaches("Redis", "memcached").
		WithTechnologies("Github Actions", "Parquet", "Avro", "REST API").
		WithPlatforms("Datadog", "Logcentral", "Github", "ServiceNow", "Azure DevOps", "Jira", "Confluence").
		WithCloudProvider("AWS",
			"ECS", "S3", "EC2", "IAM", "SecretsManager", "DynamoDB", "DAX", "Cloudwatch", "Lambda", "ECR", "Route53", "Kinesis", "SQS", "SNS",
		).
		WithSubjectMatters(smApi, smGeospatial, smBackend, smCiCd, smClusterComputing, smDistributedComputing, smContainerization, smIAC).
		finalize()
	renderProject = renderProject.WithStatus(statusMaintaining).
		WithSummary("SUMMARY HERE"). // TODO: MORE!
		WithLang("Go", Extensively).
		WithLang("Terraform", Often).
		WithLang("Docker", Regularly).
		WithLang("Bash", Some).
		WithLang("SQL", Some).
		WithTechnologies("Github Actions", "Parquet", "Avro", "REST API", "GraphQL").
		WithDbs("Postgres").
		WithCaches("Redis", "memcached").
		WithPlatforms("Datadog", "Logcentral", "Github", "ServiceNow", "GraphQL Apollo", "Azure DevOps", "Jira", "Confluence").
		WithCloudProvider("AWS",
			"ECS", "S3", "EC2", "IAM", "SecretsManager", "DynamoDB", "DAX", "Cloudwatch", "Lambda", "ECR", "Route53", "Kinesis", "SQS", "SNS",
		).
		WithSubjectMatters(smApi, smGeospatial, smBackend, smCiCd, smClusterComputing, smDistributedComputing, smContainerization, smIAC, smGraphics).
		finalize()
	statsProject = statsProject.WithStatus(statusMaintaining).
		WithSummary("SUMMARY HERE"). // TODO: MORE!
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
		WithDbs("Aurora", "Postgres").
		WithCaches("Redis", "memcached").
		WithTechnologies("CGo", "CUDA", "Github Actions", "Parquet", "Avro", "GraphQL", "SIMD", "REST API").
		WithPlatforms("Datadog", "Logcentral", "Github", "ServiceNow", "GraphQL Apollo", "Azure DevOps", "Jira", "Confluence").
		WithCloudProvider("AWS",
			"ECS", "S3", "EC2", "IAM", "SecretsManager", "DynamoDB", "DAX", "Cloudwatch", "Lambda", "ECR", "Route53", "Kinesis", "SQS", "SNS",
		).
		WithSubjectMatters(smApi, smGeospatial, smBackend, smCiCd, smClusterComputing, smDistributedComputing, smContainerization, smIAC, smGPU, smStatistics).
		finalize()
	billingProject = billingProject.WithStatus(statusBuilding).
		WithSummary("SUMMARY HERE"). // TODO: MORE!
		WithLang("Go", Extensively).
		WithLang("Terraform", Often).
		WithLang("Docker", Regularly).
		WithLang("Bash", Some).
		WithLang("SQL", Some).
		WithLang("Javascript", Rarely).
		WithCloudProvider("AWS",
			"ECS", "S3", "EC2", "IAM", "SecretsManager", "DynamoDB", "DAX", "Cloudwatch", "Lambda", "ECR", "Route53", "Kinesis", "SQS", "SNS",
		).
		WithDbs("DynamoDB", "Postgres").
		WithCaches("Redis", "memcached").
		WithTechnologies("Github Actions", "Parquet", "Avro", "REST API").
		WithPlatforms("Datadog", "Logcentral", "Github", "ServiceNow", "Azure DevOps", "Jira", "Confluence").
		WithSubjectMatters(smApi, smGeospatial, smCiCd, smClusterComputing, smDistributedComputing, smContainerization, smIAC, smStatistics, smBackend).
		finalize()
	explorerProject = explorerProject.WithStatus(statusMaintaining).
		WithSummary("SUMMARY HERE"). // TODO: MORE!
		WithLang("Scala", Often).
		WithLang("Javascript", Often).
		WithLang("Terraform", Often).
		WithLang("Docker", Regularly).
		WithLang("Bash", Some).
		WithCloudProvider("AWS",
			"ECS", "Fargate", "S3", "EC2", "IAM", "SecretsManager", "DynamoDB", "DAX", "Cloudwatch", "Lambda", "ECR", "Route53", "Kinesis", "SQS", "SNS",
		).
		WithDbs("Aurora", "DynamoDB", "Postgres").
		WithCaches("Redis", "memcached").
		WithTechnologies("Github Actions", "REST API").
		WithPlatforms("Datadog", "Logcentral", "Github", "ServiceNow", "Azure DevOps", "Jira", "Confluence").
		WithSubjectMatters(smApi, smGeospatial, smCiCd, smContainerization, smIAC, smFullStack).
		finalize()
	wqdbProject = wqdbProject.WithStatus(statusComplete).
		WithSummary("Work queue database for the various builders"). // TODO: list builders
		WithLang("SQL", Often).
		WithLang("Go", Often).
		WithLang("Scala", Often).
		WithLang("Terraform", Often).
		WithLang("Docker", Regularly).
		WithLang("Bash", Some).
		WithCloudProvider("AWS",
			"ECS", "Fargate", "S3", "EC2", "IAM", "SecretsManager", "DynamoDB", "DAX", "Cloudwatch", "Lambda", "ECR", "Route53", "Kinesis", "SQS", "SNS",
		).WithDbs("Aurora", "Postgres").
		WithTechnologies("Github Actions").
		WithPlatforms("Datadog", "Logcentral", "Github", "ServiceNow", "Azure DevOps", "Jira", "Confluence").
		WithSubjectMatters(smCiCd, smIAC, smBackend).
		finalize()
	supportProject = supportProject.WithStatus(statusComplete).
		WithSummary("Support API cluster for troubleshooting").
		WithLang("SQL", Regularly).
		WithLang("Go", Often).
		WithLang("Terraform", Often).
		WithLang("Docker", Regularly).
		WithLang("Bash", Some).
		WithCloudProvider("AWS",
			"ECS", "Fargate", "S3", "EC2", "IAM", "SecretsManager", "DynamoDB", "DAX", "Cloudwatch", "Lambda", "ECR", "Route53", "Kinesis", "SQS", "SNS",
		).
		WithTechnologies("Github Actions").
		WithPlatforms("Datadog", "Logcentral", "Github", "ServiceNow", "Azure DevOps", "Jira", "Confluence").
		WithSubjectMatters(smApi, smCiCd, smContainerization, smIAC, smBackend).
		finalize()
	scudsProject = scudsProject.WithStatus(statusComplete).
		WithSummary("SUMMARY HERE"). // TODO: MORE!
		WithLang("Scala", Extensively).
		WithLang("Go", Often).
		WithLang("Terraform", Often).
		WithLang("Docker", Regularly).
		WithLang("Bash", Some).
		WithCloudProvider("AWS",
			"ECS", "Fargate", "S3", "EC2", "IAM", "SecretsManager", "DynamoDB", "DAX", "Cloudwatch", "Lambda", "ECR", "Route53", "Kinesis", "SQS", "SNS",
		).
		WithDbs("ElasticSearch").
		WithCaches("Redis").
		WithTechnologies("ElasticSearch", "Github Actions", "Parquet", "Avro", "REST API").
		WithPlatforms("Datadog", "Logcentral", "Github", "ServiceNow", "Azure DevOps", "Jira", "Confluence").
		WithSubjectMatters(smApi, smGeospatial, smCiCd, smContainerization, smIAC, smStatistics, smBackend).
		finalize()
	ufoProject = ufoProject.WithStatus(statusMaintaining).
		WithSummary("Unified Field Operations API Cluster").
		WithLang("Scala", Extensively).
		WithLang("Go", Regularly).
		WithLang("Terraform", Regularly).
		WithLang("Docker", Regularly).
		WithLang("Bash", Some).
		WithCloudProvider("AWS",
			"ECS", "Fargate", "S3", "EC2", "IAM", "SecretsManager", "DynamoDB", "DAX", "Cloudwatch", "Lambda", "ECR", "Route53", "Kinesis", "SQS", "SNS",
		).
		WithCaches("Redis").
		WithTechnologies("Github Actions", "REST API").
		WithPlatforms("Datadog", "Logcentral", "Github", "ServiceNow", "Azure DevOps", "Jira", "Confluence").
		WithSubjectMatters(smApi, smGeospatial, smCiCd, smContainerization, smIAC, smStatistics, smBackend).
		finalize()
	goweProject = goweProject.WithStatus(statusMaintaining).
		WithSummary("SUMMARY HERE"). // TODO: MORE!
		WithLang("Go", Extensively).
		WithLang("Terraform", Regularly).
		WithLang("Bash", Some).
		WithLang("Docker", Regularly).
		WithCloudProvider("AWS",
			"ECS", "Fargate", "S3", "EC2", "IAM", "SecretsManager", "DynamoDB", "DAX", "Cloudwatch", "Lambda", "ECR", "Route53", "Kinesis", "SQS", "SNS",
		).
		WithDbs("DuckDB", "Aurora", "Postgres").
		WithCaches("Redis").
		WithTechnologies("Github Actions", "Parquet", "Avro").
		WithPlatforms("Datadog", "Logcentral", "Github", "ServiceNow", "Azure DevOps", "Jira", "Confluence").
		WithSubjectMatters(smGeospatial, smCiCd, smContainerization, smIAC, smClusterComputing, smDistributedComputing, smLinearAlgebra, smBackend).
		finalize()
	rasterRenderProject = rasterRenderProject. // TODO: this whole thing!
							WithSummary("Tile renderer API service cluster. Later replaced by the Render Cluster"). // TODO: link to render cluster
							WithLang("Scala", Extensively).
							WithLang("Terraform", Regularly).
							WithLang("Bash", Some).
							WithStatus(statusComplete).
							WithPlatforms("DroneCI", "Rally").
							WithSubjectMatters(smBackend).
							WithCloudProvider("AWS", "ECS", "EC2", "API Gateway", "IAM", "ELB", "ALB").
							finalize()
	simpsonUnivProject = simpsonUnivProject.WithStatus(statusComplete). // TODO: MORE!
										WithSummary("SUMMARY HERE"). // TODO: MORE!
										WithLang("Html", Regularly).
										WithLang("CSS", Regularly).
										WithLang("Javascript", Some).
										WithSubjectMatters(smFrontend).
										WithTechnologies("Drupal").
										WithPlatforms("Confluence").
										finalize()
	mushDbProject = mushDbProject.WithStatus(statusBuilding).
		WithSummary("SUMMARY HERE"). // TODO: MORE!
		WithLang("Go", Extensively).
		WithLang("Typescript", Extensively).
		// TODO: terraform?
		WithLang("Javascript", Extensively).
		WithLang("Docker", Often).
		WithLang("Docker Compose", Often).
		//WithLang("Kubernetes", Some). // TODO: lang or tech?
		WithDbs("mongodb").
		WithTechnologies("Github Actions", "React", "NextJs", "RFID", "NFC", "I2C", "SPI", "Cloudflare Tunnels", "Kubernetes", "REST API").
		WithPlatforms("Github", "Cloudflare", "Grafana").
		WithCloudProvider("GCP", "Google Cloud DNS", "Identity Platform").
		WithSubjectMatters(smApi, smContainerization, smDistributedComputing, smMycology, smFullStack).
		WithInterests(string(smMycology)).
		finalize()
	cvProject = cvProject.WithStatus(statusBuilding).
		WithSummary("This project! A generator which creates markdown files that can be viewed via Obsidian, or published to the web.").
		WithLang("Go", Extensively).
		WithLang("Terraform", Often).
		WithLang("Bash", Some).
		WithLang("Typescript", Some).
		WithLang("Javascript", Rarely).
		WithLang("Html", Rarely).
		WithLang("CSS", Rarely).
		WithLang("SCSS", Rarely, smFrontend).
		WithTechnologies("Github Actions", "Obsidian", "Markdown", "Quartz", "Quartz 4").
		WithPlatforms("Github").
		WithCloudProvider("AWS", "S3", "IAM", "Cloudfront", "Cloudfront Functions", "ACM").
		WithSubjectMatters(smCiCd, smFrontend, smDevOps, smServerless).
		finalize()
	linksPage = linksPage.
		WithStatus(statusComplete).
		WithLang("Html", Regularly).
		WithLang("CSS", Some).
		WithLang("Javascript", Rarely).
		WithSubjectMatters(smFrontend)
	measurementsProject = measurementsProject.WithStatus(statusShelved).
		WithSummary("A library to do arbitrary conversion rates, which can be time-dependent, multivariate, nested, and compound").
		WithLang("Go", Extensively).
		WithTechnologies("Github Actions").
		WithPlatforms("Github").
		WithSubjectMatters(smThermodynamics).
		finalize()
	nfcScannerProject = nfcScannerProject.WithStatus(statusBuilding).
		WithSummary("SUMMARY HERE"). // TODO: MORE!
		WithLang("Go", Extensively).
		WithTechnologies("Github Actions", "Webhooks", "Websockets", "Server-Sent Events", "Pub-Sub", "NFC", "I2C", "SPI", "RFID", "REST API").
		WithPlatforms("Github").
		WithSubjectMatters(smApi, smContainerization, smDistributedComputing).
		WithInterests(string(smMycology)).
		finalize()
	coreShufflerProject = coreShufflerProject.WithStatus(statusComplete).
		// TODO: MORE! LINKS!
		WithSummary("Related to capstone project. Takes fuel assembly layout input, allows you to shuffle them, and outputs the CASMO and SIMULATE files necessary to run a cycle with the new layout.").
		// TODO: ADD A PICTURE!!!!!
		WithLang("Javascript", Extensively).
		WithLang("Html", Extensively).
		WithLang("CSS", Extensively).
		WithTechnologies("JQuery").
		WithPlatforms("Github", "Gitlab").
		WithSubjectMatters(smNuclearEngineering, smParticlePhysics, smFluidMechanics).
		WithInterests(string(smNuclearEngineering)).
		finalize()
	CharityProject = CharityProject.WithStatus(statusComplete).
		WithSummary("Designed, created, and hosted website for a local charity."). // TODO: MORE!
		WithLang("Html", Often).
		WithLang("CSS", Often).
		WithLang("Javascript", Often).
		WithSubjectMatters(smFrontend).
		finalize()
	WellAwareProject = WellAwareProject.WithStatus(statusComplete).
		WithSummary("Designed, created, and hosted website for WellAware NC. A University of North Carolina, Chapel Hill lab-affiliated organization."). // TODO: this
		WithLang("Html", Often).
		WithLang("CSS", Often).
		WithLang("Javascript", Often).
		WithCloudProvider("AWS",
			"S3", "IAM", "Route53").
		WithTechnologies("Gitlab CI", "React", "NodeJS").
		WithSubjectMatters(smCiCd, smDevOps, smFrontend).
		finalize()
	CritColaProject = CritColaProject.WithStatus(statusComplete).
		WithSummary("SUMMARY HERE"). // TODO: this
		WithLang("Docker", Regularly).
		WithLang("Terraform", Often).
		WithLang("Java", Rarely).
		WithLang("Bash", Some).
		WithLang("Javascript", Rarely).
		WithTechnologies("Gitlab CI").
		WithPlatforms("Gitlab", "Discord").
		WithCloudProvider("AWS",
			"S3", "EC2", "IAM", "SecretsManager", "ECR", "Route53").
		WithSubjectMatters(smBackend, smCiCd, smDevOps, smScripting).
		WithInterests("Gaming").
		finalize()
	MastersDataAnalysisProject = MastersDataAnalysisProject.WithStatus(statusComplete).
		WithSummary("Created scripts to automate the data analysis for an individual's project for his master's degree in Public Health at UNC Chapel Hill."). // TODO: this
		WithLang("Python", Often).
		WithTechnologies("Gitlab CI").
		WithPlatforms("Gitlab").
		WithSubjectMatters(smStatistics, smLinearAlgebra).
		finalize()
	capstoneProject = capstoneProject.WithStatus(statusComplete).
		WithSummary("Capstone project for my undergraduate degree at NC State. Designed a 3 year fuel cycle utilizing Uranium Silicide Accident-Tolerant fuel rods. Sponsored by the Catawba Nuclear plant and Duke Energy."). // TODO: this
		WithLang("Javascript", Some).
		WithLang("Fortran", Often).
		WithLang("Bash", Some).
		WithTechnologies("JQuery", "Git", "CASMO4e", "SIMULATE3").
		WithPlatforms("Github").
		WithSubjectMatters(smNuclearEngineering, smLinearAlgebra, smFluidMechanics, smThermodynamics, smParticlePhysics).
		WithInterests(string(smNuclearEngineering)).
		finalize()
	reactorAnalysisFinal = reactorAnalysisFinal.
		WithLang("Javascript", Extensively).
		WithStatus(statusComplete).
		WithSummary("Final exam for the final course of my Nuclear Engineering Bachelors degree"). // TODO: this! LInk!
		WithSubjectMatters(smNuclearEngineering, smThermodynamics, smFluidMechanics).
		finalize()
	monteCarloProject = monteCarloProject.
		WithSummary(fixmeLink).
		WithStatus(statusComplete).
		WithInterests(string(smParticlePhysics)). // TODO: ok?
		WithLang("Java", Often).
		WithSubjectMatters(smParticlePhysics, smNuclearEngineering).
		finalize()
	projectLinAlgCryptography = projectLinAlgCryptography.
		WithSummary("Created and analyzed the efficacy of a cryptographic algorithm utilizing basic matrix mathematics.").
		WithMiscSkills("Linear Algebra").
		WithInterests("Cryptography").
		WithSubjectMatters(smCybersecurity, smCryptography).
		WithInterests(string(smCryptography), string(smCybersecurity)).
		finalize()
	cherenkovProject = cherenkovProject.WithStatus(statusComplete).
		WithSummary("Designed and built a sensor for Cherenkov radiation, which was tested by lowering the sensor into the core of the PULSTAR reactor at "+schoolNCSU.Link()+". The sensor utilized and Arduino for signal processing.").
		WithLang("Python", Often).
		WithLang("Fortran", Often).
		WithTechnologies("Arduino").
		WithPlatforms("Github").
		WithSubjectMatters(smNuclearEngineering, smLinearAlgebra, smParticlePhysics, smStatistics, smEmbeddedSystems, smElectronics).
		WithInterests(string(smNuclearEngineering)).
		finalize()
	roboticsTeamProject = roboticsTeamProject.WithStatus(statusComplete).
		WithSummary("FIRST team 3720 robots for competition each year. Design, construction, programming, and competing.").
		WithLang("Java", Extensively).
		WithLang("Html", Some).
		WithLang("CSS", Some).
		WithLang("Javascript", Some).
		WithTechnologies("PWM").
		WithSubjectMatters(smRobotics, smLinearAlgebra, smEmbeddedSystems, smElectronics).
		WithInterests(string(smRobotics)).
		finalize()
	cncLaserCutterProject = cncLaserCutterProject.WithStatus(statusComplete).
		WithSummary("Designed and build a 2-axis CNC laser cutter. Worked with standard G&M codes (except the ones used to typically change drill bit speed).").
		WithLang("Java", Often).
		WithLang("Javascript", Regularly).
		WithLang("Html", Regularly).
		WithTechnologies("PWM", "G and M codes").
		WithSubjectMatters(smRobotics, smLinearAlgebra, smEmbeddedSystems, smElectronics).
		finalize()
	aerospaceFinalProject = aerospaceFinalProject.WithStatus(statusComplete).
		WithSummary("Designed and built a rocket from scratch, including the fuel (Used Potassium Nitrate and Glucose/Sucrose).").
		WithLang("Java", Minimal).
		WithLang("Javascript", Minimal).
		WithSubjectMatters(smFluidMechanics, smElectronics, smChemistry).
		finalize()
	miscSmallPersonalProjects = miscSmallPersonalProjects.
		WithSummary("A catch-all for my miscellaneous personal projects I didn't want to add entire pages for").
		WithStatus(statusShelved).
		WithLang("Python", Regularly, smScripting).
		WithLang("Fortran", Some, smScripting, smThermodynamics, smFluidMechanics, smParticlePhysics, smNuclearEngineering, smEducation, smEmbeddedSystems, smStatistics, smLinearAlgebra, smStatics).
		WithLang("Solidity", Some, smCryptocurrency).
		WithLang("Go", Extensively, smThermodynamics).
		WithLang("Java", Some, smBackend).
		WithLang("Javascript", Often, smFrontend).
		WithLang("Typescript", Often, smFullStack).
		WithLang("Kotlin", Rarely, smBackend).
		WithTechnologies("Kubernetes").
		WithSubjectMatters(smApi, smCryptocurrency).
		finalize()
	teiProjects = teiProjects.
		WithSummary("A catch-all for all projects at my position at TEI.").
		WithStatus(statusComplete).
		WithSubjectMatters(smStructuralEngineering, smCivilEngineering).
		finalize()
	taePhotoImporter = taePhotoImporter.
		WithSummary("Script to import photos taken in the field into the proper format and structure for use in engineering documentation").
		WithLang("Javascript", Some).
		WithStatus(statusComplete).
		WithSubjectMatters(smStructuralEngineering, smCivilEngineering).
		finalize()
	mafcProjects = mafcProjects.
		WithSummary("Lifeguarding stuff, I guess.").
		WithStatus(statusComplete).
		WithSubjectMatters(smFirstAid).
		finalize()
	ArrowNailProject = ArrowNailProject.
		WithLang("Javascript", Extensively).
		WithLang("Html", Regularly).
		WithLang("CSS", Regularly).
		WithSummary("Geospatial mapping app to track historical hail instances for use by a roofing company").
		WithStatus(statusShelved).
		WithTechnologies("React", "REST API").
		WithSubjectMatters(smApi, smGeospatial).
		WithCloudProvider("AWS", "Lambda", "RDS"). // TODO: MORE STUFF?
		finalize()
	WildlifeRProject = WildlifeRProject.
		WithStatus(statusComplete).
		WithLang("R", Often, smGeospatial).
		WithSummary("Utilized spatiotemporal data for wildlife in a specified area over a specified date range in order to produce population density maps").
		WithSubjectMatters(smGeospatial, smStatistics).
		finalize()
	jamfProject = jamfProject.
		WithSummary("JAMF companywide setup. " + fixmeLink).
		WithStatus(statusComplete).
		WithPlatforms("Jamf").
		WithSubjectMatters(smCybersecurity).
		finalize()
	smallImprovementsProject = smallImprovementsProject.
		WithLang("Javascript", Often).
		WithSummary("Created a Small Improvements Slack bot for Source Allies to track goal creation and achievement").
		WithCloudProvider("AWS", "Lambda", "DynamoDB", "SAM", "Cloudformation"). // TODO: AWS
		WithPlatforms("Slack", "Small Improvements", "Github").
		WithTechnologies("Github Actions", "REST API").
		WithStatus(statusComplete).
		WithSubjectMatters(smApi, smBackend, smServerless).
		finalize()
	internalResumeGeneratorProject = internalResumeGeneratorProject.
		WithSummary("Updated company internal resume generator").
		// TODO: any aws in here???
		WithStatus(statusComplete).
		WithLang("Java", Often, smBackend).
		WithLang("Html", Some).
		WithLang("CSS", Some).
		WithLang("Javascript", Some).
		WithTechnologies("Spring", "REST API").
		WithSubjectMatters(smApi, smFullStack).
		WithTechnologies("Github Actions").
		finalize()
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
