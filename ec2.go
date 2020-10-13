package main

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Ec2Handler defines the interface and parameters
type Ec2Handler struct {
	Service Ec2Iface
	Eni     *string
}

// Ec2Client return the aws client of EC2
func Ec2Client(region string) *ec2.EC2 {
	return ec2.New(session.New(), aws.NewConfig().WithRegion(region))
}

// Ec2Iface defines the features implemented in EC2
type Ec2Iface interface {
	DescribeNetworkInterfaces(*ec2.DescribeNetworkInterfacesInput) (*ec2.DescribeNetworkInterfacesOutput, error)
}

// PublicIP return the public ip associated with the task ENI
func (h *Ec2Handler) PublicIP() (*string, error) {
	input := &ec2.DescribeNetworkInterfacesInput{
		NetworkInterfaceIds: []*string{
			aws.String(*h.Eni),
		},
	}

	result, err := h.Service.DescribeNetworkInterfaces(input)
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
		return nil, err
	}

	if len(result.NetworkInterfaces) > 0 {
		return result.NetworkInterfaces[0].Association.PublicIp, err
	}
	return nil, errors.New("No ENI found")
	// (TODO)If public IP not found, attach an EIP
}
