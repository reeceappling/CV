package tags

import "appli.ng/cv/generator/utils"

type Tag string

var tags = utils.Set[Tag]{}

func NewTag(name string) Tag {
	out := Tag(name)
	tags.Add(out)
	return out
}

type Field []Tag

func (tags Field) AsStrings() []string {
	out := make([]string, len(tags))
	for i, tag := range tags {
		out[i] = string(tag)
	}
	return out
}

var (
	AgTech                = NewTag("AgTech")
	Consulting            = NewTag("Consulting")
	Robotics              = NewTag("Robotics")
	Cryptocurrency        = NewTag("Cryptocurrency") // TODO: USE SOMEWHERE
	ClusterComputing      = NewTag("Cluster Computing")
	DistributedComputing  = NewTag("Distributed Computing")
	Containerization      = NewTag("Containerization")
	IAC                   = NewTag("Infrastructure As Code")
	DevOps                = NewTag("DevOps")
	CiCd                  = NewTag("CI-CD")
	Statistics            = NewTag("Statistics")
	Topology              = NewTag("Topology")
	ParticlePhysics       = NewTag("Particle Physics")
	NuclearEngineering    = NewTag("Nuclear Engineering")
	CivilEngineering      = NewTag("Civil Engineering")
	StructuralEngineering = NewTag("Structural Engineering")
	Telecom               = NewTag("Telecommunications")
	SWE                   = NewTag("Software Engineering")
	Thermodynamics        = NewTag("Thermodynamics")
	FluidMechanics        = NewTag("Fluid Mechanics")
	Cybersecurity         = NewTag("Cybersecurity")
	Cryptography          = NewTag("Cryptography")
	EmbeddedSystems       = NewTag("Embedded Systems")
	Electronics           = NewTag("Electronics")
	AI                    = NewTag("AI")
	LinearAlgebra         = NewTag("Linear Algebra")
	Graphics              = NewTag("Graphics")
	GPU                   = NewTag("GPU")
	Mycology              = NewTag("Mycology")
	Frontend              = NewTag("Frontend")
	Backend               = NewTag("Backend")
	FullStack             = NewTag("Full Stack")
	CloudComputing        = NewTag("Cloud Computing")
	Observability         = NewTag("Observability")
	Serverless            = NewTag("Serverless")
	Logging               = NewTag("Logging")
	FirstAid              = NewTag("First Aid")
	Scripting             = NewTag("Scripting")
	Statics               = NewTag("Statics")
	Networking            = NewTag("Networking") // TODO: maybe get rid of
	Charity               = NewTag("Charity")
	Entertainment         = NewTag("Entertainment")
	Education             = NewTag("Education")
	Construction          = NewTag("Construction")
	PublicHealth          = NewTag("Public Health")
	Ecology               = NewTag("Ecology")
)
