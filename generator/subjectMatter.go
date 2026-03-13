package main

import (
	"appli.ng/cv/generator/utils"
	"strings"
)

var subjectMatters map[SubjectMatter]map[string]utils.Set[string] = nil // Map of SM to map of linkType to []link

type SubjectMatter string

func (sm SubjectMatter) Link() string {
	return linkFor(string(sm), "cv", "subjectMatter", withoutSpaces(string(sm)))
}
func (pg *SubjectMatter) EntryType() string {
	return "Subject Matter"
}
func (sm SubjectMatter) Bytes() []byte {
	b := strings.Builder{}
	b.WriteString(frontmatterFor(string(sm)))
	b.WriteString("View the backlinks on this page to see everything that references this subject matter\n")
	// TODO: backlinks
	return []byte(b.String())
}

func NewSubjectMatter(sm string) SubjectMatter {
	out := SubjectMatter(sm)
	if _, exists := subjectMatters[out]; !exists {
		subjectMatters[out] = map[string]utils.Set[string]{}
	}
	return out
}

type SubjectMattersField struct {
	SubjectMatters utils.Set[SubjectMatter]
}

// TODO: ensure works properly!
func withSubjectMatters(pg Linkable, setIn utils.Set[SubjectMatter], sms ...SubjectMatter) utils.Set[SubjectMatter] {
	// handle full stack
	for _, smIn := range sms {
		if smIn == smFullStack {
			sms = append(sms, smFrontend, smBackend)
			break
		}
	}
	for _, sm := range sms {
		if current, exists := subjectMatters[sm]; exists {
			newInternalMap := current
			if currentSet, setExists := current[pg.EntryType()]; setExists {
				currentSet.Add(pg.Link())
				newInternalMap[pg.EntryType()] = currentSet // TODO: use new set?
			} else {
				newInternalMap[pg.EntryType()] = utils.SetFrom(pg.Link())
			}
			subjectMatters[sm] = newInternalMap
		} else {
			subjectMatters[sm] = map[string]utils.Set[string]{
				pg.EntryType(): utils.SetFrom(pg.Link()),
			}
		}
	}
	setIn.Add(sms...)
	return setIn
}

var (
	smRobotics              = NewSubjectMatter("Robotics")
	smCryptocurrency        = NewSubjectMatter("Cryptocurrency")
	smClusterComputing      = NewSubjectMatter("Cluster Computing")
	smDistributedComputing  = NewSubjectMatter("Distributed Computing")
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
	smCloudComputing        = NewSubjectMatter("Cloud Computing")
	smObservability         = NewSubjectMatter("Observability")
	smServerless            = NewSubjectMatter("Serverless")
	smLogging               = NewSubjectMatter("Logging")
	smFirstAid              = NewSubjectMatter("First Aid")
	smScripting             = NewSubjectMatter("Scripting")
	smEducation             = NewSubjectMatter("Education")
	smStatics               = NewSubjectMatter("Statics")
	smNetworking            = NewSubjectMatter("Networking") // TODO: maybe get rid of
)

func setupSubjectMatters() {
	subjectMatters = map[SubjectMatter]map[string]utils.Set[string]{
		smCloudComputing: map[string]utils.Set[string]{
			"Cloud Provider": {},
			"Cloud Service":  {},
		},
	}
	setupLanguageSubjectMatters()
	setupServiceSubjectMatters()
	setupTechSubjectMatters()
	setupPlatformSubjectMatters()
}
