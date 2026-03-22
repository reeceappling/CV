package main

import (
	"appli.ng/cv/generator/utils"
	"strings"
)

var cloudServices = map[string]*CloudServicePage{}

type CloudServicePage struct {
	Name     string
	provider string // TODO: add provider to page????
	*tracked
	SubjectMattersField
}

func (pg *CloudServicePage) EntryType() string {
	return "Cloud Service"
}

func (pg *CloudServicePage) WithSubjectMatters(sms ...SubjectMatter) *CloudServicePage {
	if pg == nil {
		return nil
	}
	pg.SubjectMatters = withSubjectMatters(pg, pg.SubjectMatters, sms...)
	return pg
}

func (pg *CloudServicePage) Dst() string {
	return dstFor("cv", "service", withoutSpaces(pg.Name))
}

func (pg *CloudServicePage) Title() string {
	if pg == nil {
		return noLinkText
	}
	return pg.Name
}

func NewService(name string, provider string) *CloudServicePage {
	out := &CloudServicePage{
		Name:     name,
		provider: provider,
		tracked:  newTracked(),
		SubjectMattersField: SubjectMattersField{
			SubjectMatters: utils.SetFrom(smCloudComputing),
		},
	}
	subjectMatters[smCloudComputing][out.EntryType()].Add() // TODO: ok?
	cloudServices[name] = out
	addLinkable(strings.ToLower(name), out)
	return out
}

func setupServiceSubjectMatters() {
	NewService("Cloudwatch", "AWS").
		WithSubjectMatters(smObservability)
	NewService("Lambda", "AWS").
		WithSubjectMatters(smServerless, smCloudComputing)
	NewService("Fargate", "AWS").
		WithSubjectMatters(smServerless, smCloudComputing)
	NewService("Aurora", "AWS").
		WithSubjectMatters(smServerless, smDatabase)
	//WithTags("Database") // TODO: ok???
	NewService("DynamoDB", "AWS").
		WithSubjectMatters(smServerless, smDatabase)
	NewService("RDS", "AWS").
		WithSubjectMatters(smDatabase)
	NewService("DocumentDB", "AWS").
		WithSubjectMatters(smDatabase)
	NewService("DAX", "AWS").
		WithSubjectMatters(smCaching)
	NewService("Elasticache", "AWS"). // TODO: use!
						WithSubjectMatters(smCaching)
	//NewService("DocumentDB", "AWS"). // TODO: use if used
	//	WithSubjectMatters(smServerless)
	NewService("IAM", "AWS").
		WithSubjectMatters(smCybersecurity)
	NewService("ECS", "AWS").
		WithSubjectMatters(smCloudComputing, smClusterComputing)
	// TODO: ASGs? CapacityProviders?
	// TODO: Glue?
	NewService("EC2", "AWS").
		WithSubjectMatters(smCloudComputing)
	NewService("Secrets Manager", "AWS").
		WithSubjectMatters(smCybersecurity)
	NewService("Key Management Service", "AWS").
		WithSubjectMatters(smCybersecurity)
	NewService("Cloudfront", "AWS").
		WithSubjectMatters(smFullStack) // TODO: OK?
	NewService("ECS", "AWS").
		WithSubjectMatters(smCloudComputing, smClusterComputing)
	NewService("ECR", "AWS").
		WithSubjectMatters(smCloudComputing, smClusterComputing, smContainerization)
	NewService("SAM", "AWS").
		WithSubjectMatters(smServerless)
	NewService("EBS", "AWS").
		WithSubjectMatters(smBackend)
	//NewService("EFS", "AWS"). // TODO: use somewhere when used
	//	WithSubjectMatters(smBackend)
	//NewService("EKS", "AWS"). // TODO: use somewhere when used
	//	WithSubjectMatters(smBackend, smClusterComputing)
	NewService("Kinesis", "AWS"). // TODO: use somewhere when used
					WithSubjectMatters(smBackend, smClusterComputing) // TODO: event driven???
	NewService("Cloudformation", "AWS"). // TODO: use somewhere
						WithSubjectMatters(smIAC)
	NewService("Chime", "AWS"). // TODO: use somewhere
					WithSubjectMatters(smCommunication)
	NewService("SQS", "AWS"). // TODO: use somewhere
					WithSubjectMatters(smBackend) // TODO: event driven???
	NewService("SNS", "AWS"). // TODO: use somewhere
					WithSubjectMatters(smBackend) // TODO: event driven???
	//NewService("SES", "AWS"). // TODO: use somewhere when used
	//	WithSubjectMatters(smBackend)
	NewService("VPC", "AWS"). // TODO: use somewhere when used
					WithSubjectMatters(smBackend, smNetworking)
	NewService("VPN", "AWS"). // TODO: use somewhere when used
					WithSubjectMatters(smBackend, smNetworking)
	//NewService("MQ", "AWS"). // TODO: use somewhere
	//				WithSubjectMatters(smBackend) // TODO: event driven???
	//NewService("Cognito", "AWS"). // TODO: use somewhere when used
	//				WithSubjectMatters(smCommunication)
	//NewService("Athena", "AWS"). // TODO: use somewhere when used
	//	WithSubjectMatters(smBackend)
	//NewService("Bedrock", "AWS"). // TODO: use somewhere when used
	//	WithSubjectMatters(smBackend, smAI)
	NewService("Corretto", "AWS"). // TODO: use somewhere when used
					WithSubjectMatters(smBackend)
	NewService("S3", "AWS").
		WithSubjectMatters(smBackend)

}
