package main

import "appli.ng/cv/generator/utils"

type Tag string

var tags utils.Set[Tag]

func NewTag(name string) Tag {
	out := Tag(name)
	tags.Add(out)
	return out
}

var (
	tagRobotics              = NewTag("Robotics")
	tagCryptocurrency        = NewTag("Cryptocurrency") // TODO: USE SOMEWHERE
	tagClusterComputing      = NewTag("Cluster Computing")
	tagDistributedComputing  = NewTag("Distributed Computing")
	tagContainerization      = NewTag("Containerization")
	tagIAC                   = NewTag("Infrastructure As Code")
	tagDevOps                = NewTag("DevOps")
	tagCiCd                  = NewTag("CI-CD")
	tagStatistics            = NewTag("Statistics")
	tagTopology              = NewTag("Topology")
	tagParticlePhysics       = NewTag("Particle Physics")
	tagNuclearEngineering    = NewTag("Nuclear Engineering")
	tagCivilEngineering      = NewTag("Civil Engineering")
	tagStructuralEngineering = NewTag("Structural Engineering")
	tagThermodynamics        = NewTag("Thermodynamics")
	tagFluidMechanics        = NewTag("Fluid Mechanics")
	tagCybersecurity         = NewTag("Cyberecurity")
	tagCryptography          = NewTag("Cryptography")
	tagEmbeddedSystems       = NewTag("Embedded Systems")
	tagElectronics           = NewTag("Electronics")
	tagAI                    = NewTag("AI")
	tagLinearAlgebra         = NewTag("Linear Algebra")
	tagGraphics              = NewTag("Graphics")
	tagGPU                   = NewTag("GPU")
	tagMycology              = NewTag("Mycology")
	tagFrontend              = NewTag("Frontend")
	tagBackend               = NewTag("Backend")
	tagFullStack             = NewTag("Full Stack")
	tagCloudComputing        = NewTag("Cloud Computing")
	tagObservability         = NewTag("Observability")
	tagServerless            = NewTag("Serverless")
	tagLogging               = NewTag("Logging")
	tagFirstAid              = NewTag("First Aid")
	tagScripting             = NewTag("Scripting")
	tagEducation             = NewTag("Education")
	tagStatics               = NewTag("Statics")
	tagNetworking            = NewTag("Networking") // TODO: maybe get rid of
)
