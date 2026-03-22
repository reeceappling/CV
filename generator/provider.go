package main

import (
	"appli.ng/cv/generator/utils"
	"fmt"
	"strings"
)

var providers = map[string]*CloudProviderPage{}

func init() {
	//NewCloudProvider("AWS",
	//	"ApiGateway", "Aurora", "Autoscaling",
	//	"CloudFront", "CloudWatch", "Config", "Certificate Manager",
	//	"DynamoDB", "DAX",
	//	"EBS", "EC2", "ECR", "EFS", "ELB", "Elasticache", "ElasticSearch", "EventBridge",
	//	"Fargate",
	//	"IAM", "KMS",
	//	"Kinesis",
	//	"Lambda",
	//	"OpenSearch", "Organizations",
	//	"RDS", "Route53",
	//	"S3", "SNS", "SQS", "Secrets Manager", "STS", "SAM", "SSM",
	//	"VPC", "VPC Lattice", "VPN")
}

type CloudProviderPage struct {
	Name string
	*tracked
	Services map[string]*CloudServicePage
	SubjectMattersField
}

// TODO: Bytes() for cloud provider that also uses Services!
func (pg *CloudProviderPage) Bytes() []byte {
	b := strings.Builder{}
	b.WriteString(frontmatterFor(pg.Name, "Cloud Provider"))
	b.WriteString("# Services\n")
	for _, svc := range pg.Services {
		b.WriteString(fmt.Sprintf("- %s\n", Link(svc)))
	}
	b.WriteString("\n" + string(pg.tracked.Bytes("", nil, ""))) // TODO: ok?
	return []byte(b.String())
}

func (pg *CloudProviderPage) EntryType() string {
	return "Cloud Provider"
}
func (pg *CloudProviderPage) Dst() string {
	return dstFor("cv", "provider", withoutSpaces(pg.Name))
}

func (pg *CloudProviderPage) Title() string {
	if pg == nil {
		return noLinkText
	}
	return pg.Name
}

func NewCloudProvider(name string, services ...string) *CloudProviderPage {
	prov, exists := providers[name]
	if exists {
		panic("provider already exists")
	}
	svcs := make(map[string]*CloudServicePage, len(services))
	for _, svc := range services {
		serv, existing := cloudServices[svc]
		if !existing || serv == nil {
			serv = NewService(svc, name)
		}
		svcs[svc] = serv
	}
	prov = &CloudProviderPage{
		Name:     name,
		tracked:  newTracked(),
		Services: svcs,
		SubjectMattersField: SubjectMattersField{
			SubjectMatters: utils.SetFrom(smCloudComputing),
		},
	}
	subjectMatters[smCloudComputing][prov.EntryType()].Add() // TODO: ok?
	providers[name] = prov
	addLinkable(name, prov)
	return prov
}

func (pg *CloudProviderPage) AddServices(services ...*CloudServicePage) *CloudProviderPage {
	for _, serv := range services {
		svc, exists := pg.Services[serv.Name]
		if !exists {
			svc = &CloudServicePage{
				Name:                serv.Name,
				tracked:             newTracked(),
				SubjectMattersField: SubjectMattersField{SubjectMatters: map[SubjectMatter]struct{}{}},
			}
		}
		pg.Services[serv.Name] = svc
	}
	return pg
}
func (pg *CloudProviderPage) AddServicesByName(services ...string) []*CloudServicePage {
	out := make([]*CloudServicePage, len(services))
	for i, serv := range services {
		svc, exists := pg.Services[serv]
		if !exists || svc == nil {
			svc = NewService(serv, pg.Name)
		}
		pg.Services[serv] = svc
		out[i] = svc
	}
	return out
}
