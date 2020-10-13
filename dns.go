package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

// DNSHandler defines the service interface and parameters
type DNSHandler struct {
	Service    DNSIface
	RecordName *string
	HostZoneID *string
	PubIP      *string
}

//DNSClient returns the aws client of route53
func DNSClient(region string) *route53.Route53 {
	return route53.New(session.New(), aws.NewConfig().WithRegion(region))
}

// DNSIface define the features implemented in route53
type DNSIface interface {
	ChangeResourceRecordSets(input *route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error)
	WaitUntilResourceRecordSetsChanged(input *route53.GetChangeInput) error
}

// RecordSet creates or update record set in route53
func (h *DNSHandler) RecordSet() error {
	log.Printf("Create or update record set: %v", *h.RecordName)
	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{ // Required
				{ // Required
					Action: aws.String("UPSERT"), // Required
					ResourceRecordSet: &route53.ResourceRecordSet{ // Required
						Name: aws.String(*h.RecordName), // Required
						Type: aws.String("A"),           // Required
						ResourceRecords: []*route53.ResourceRecord{
							{ // Required
								Value: aws.String(*h.PubIP), // Required
							},
						},
						TTL:           aws.Int64(60),
						Weight:        aws.Int64(100),
						SetIdentifier: aws.String(*h.PubIP),
					},
				},
			},
			Comment: aws.String("Record Name for Public ECS task"),
		},
		HostedZoneId: aws.String(*h.HostZoneID),
	}
	result, err := h.Service.ChangeResourceRecordSets(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}
	}
	return h.Service.WaitUntilResourceRecordSetsChanged(&route53.GetChangeInput{
		Id: aws.String(*result.ChangeInfo.Id),
	})
}
