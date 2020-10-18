package main

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/route53"
)

// DNSHandler defines the service interface and parameters
type DNSHandler struct {
	Service        CfnIface
	RecordName     *string
	HostedZoneName *string
	PubIP          *string
}

// DNSClient returns the aws client of route53
func DNSClient(region string) *cfn.CloudFormation {
	return cfn.New(session.New(), aws.NewConfig().WithRegion("ap-southeast-2"))
}

// CfnIface expose the CFN feature that used to deploy DNS
type CfnIface interface {
	CreateStack(input *cfn.CreateStackInput) (*cfn.CreateStackOutput, error)
	UpdateStack(input *cfn.UpdateStackInput) (*cfn.UpdateStackOutput, error)
	WaitUntilStackCreateComplete(input *cfn.DescribeStacksInput) error
	WaitUntilStackUpdateComplete(input *cfn.DescribeStacksInput) error
}

func (h *DNSHandler) cfnTemplate() string {
	template := cloudformation.NewTemplate()

	template.Resources["BastionDNS"] = &route53.RecordSet{
		HostedZoneName:  *h.HostedZoneName,
		Name:            *h.RecordName,
		Type:            "A",
		ResourceRecords: []string{*h.PubIP},
		TTL:             "60",
	}
	// Template string
	j, err := template.JSON()
	if err != nil {
		log.Panicf("CFN Failed to generate JSON: %s\n", err)
	} else {
		log.Printf("%s\n", string(j))
	}
	return string(j)
}

// RecordSet creates or update record set in route53
func (h *DNSHandler) RecordSet() error {
	tmpBody := h.cfnTemplate()
	stackName := *h.RecordName
	_, err := h.Service.CreateStack(&cfn.CreateStackInput{
		StackName:    aws.String(stackName),
		TemplateBody: aws.String(tmpBody),
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "AlreadyExistsException") {
			log.Println("Stack already exists, update stack")
			output, err := h.Service.UpdateStack(&cfn.UpdateStackInput{
				StackName:    aws.String(stackName),
				TemplateBody: aws.String(tmpBody),
			})
			if err != nil {
				log.Fatalf("Failed to update stack: %v", err)
				return err
			}
			log.Println(*output.StackId)
			err = h.Service.WaitUntilStackUpdateComplete(&cfn.DescribeStacksInput{StackName: aws.String(stackName)})
			if err != nil {
				log.Println(err)
				return err
			}
			return err
		}
		log.Printf("Failed to create stack %v", err)
		return err
	}

	// Wait until stack is created
	err = h.Service.WaitUntilStackCreateComplete(&cfn.DescribeStacksInput{StackName: aws.String(stackName)})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
