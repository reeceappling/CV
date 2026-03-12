package main

import (
	"appli.ng/cv/generator/utils"
	"strings"
)

var subjectMatters = utils.Set[SubjectMatter]{}

type SubjectMatter string

func (sm SubjectMatter) Link() string {
	return linkFor(string(sm), "cv", "subjectMatter", withoutSpaces(string(sm)))
}
func (sm SubjectMatter) Bytes() []byte {
	b := strings.Builder{}
	b.WriteString(frontmatterFor(string(sm)))
	b.WriteString("NEED TO INSERT BACKLINKS\n") // TODO: THIS
	return []byte(b.String())
}

func NewSubjectMatter(sm string) SubjectMatter {
	out := SubjectMatter(sm)
	subjectMatters.Add(out)
	return out
}

var (
	smRobotics              = NewSubjectMatter("Robotics")
	smCryptocurrency        = NewSubjectMatter("Cryptocurrency") // TODO: USE SOMEWHERE
	smClusterComputing      = NewSubjectMatter("Cluster Computing")
	smDistributedComputing  = NewSubjectMatter("Cluster Computing")
	smContainerization      = NewSubjectMatter("Containerization")
	smIAC                   = NewSubjectMatter("Infrastructure As Code")
	smDevOps                = NewSubjectMatter("DevOps")
	smCiCd                  = NewSubjectMatter("CI-CD")
	smStatistics            = NewSubjectMatter("Statistics")
	smTopology              = NewSubjectMatter("Topology")
	smParticlePhysics       = NewSubjectMatter("Particle Physics")
	smNuclearEngineering    = NewSubjectMatter("Nuclear Engineering")
	smCivilEngineering      = NewSubjectMatter("Civil Engineering")
	smStructuralEngineering = NewSubjectMatter("Structural Engineering")
	smThermodynamics        = NewSubjectMatter("Thermodynamics")
	smFluidMechanics        = NewSubjectMatter("Fluid Mechanics")
	smCybersecurity         = NewSubjectMatter("Cyberecurity")
	smCryptography          = NewSubjectMatter("Cryptography")
	smEmbeddedSystems       = NewSubjectMatter("Embedded Systems")
	smElectronics           = NewSubjectMatter("Electronics")
	smAI                    = NewSubjectMatter("AI")
	smLinearAlgebra         = NewSubjectMatter("Linear Algebra")
	smGraphics              = NewSubjectMatter("Graphics")
	smGPU                   = NewSubjectMatter("GPU")
	smMycology              = NewSubjectMatter("Mycology")
	smFrontend              = NewSubjectMatter("Frontend")
	smBackend               = NewSubjectMatter("Backend")
	smFullStack             = NewSubjectMatter("Full Stack")
)
