package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

type DNSHandler struct {
	Service    DNSIface
	RecordName *string
	HostZoneID *string
	PubIP      *string
}

func DNSClient(region string) *route53.Route53 {
	return route53.New(session.New(), aws.NewConfig().WithRegion(region))
}

type DNSIface interface {
	ChangeResourceRecordSets(input *route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error)
}

func (h *DNSHandler) RecordSet() {
	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{ // Required
				{ // Required
					Action: aws.String("UPSERT"), // Required
					ResourceRecordSet: &route53.ResourceRecordSet{ // Required
						Name: aws.String(*h.RecordName), // Required
						Type: aws.String("CNAME"),       // Required
						ResourceRecords: []*route53.ResourceRecord{
							{ // Required
								Value: aws.String(*h.PubIP), // Required
							},
						},
						TTL:           aws.Int64(300),
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
	log.Println(result)

}
