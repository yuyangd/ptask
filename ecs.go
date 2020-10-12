package main

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// Get TaskEni

type EcsHandler struct {
	Service ECSIface
	Cluster *string
	TaskArn *string
}

func EcsClient(region string) *ecs.ECS {
	return ecs.New(session.New(), aws.NewConfig().WithRegion(region))
}

type ECSIface interface {
	DescribeTasks(*ecs.DescribeTasksInput) (*ecs.DescribeTasksOutput, error)
}

func (h *EcsHandler) TaskEni() (*string, error) {
	input := &ecs.DescribeTasksInput{
		Cluster: aws.String(*h.Cluster),
		Tasks: []*string{
			aws.String(*h.TaskArn),
		},
	}
	result, err := h.Service.DescribeTasks(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				log.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeClientException:
				log.Println(ecs.ErrCodeClientException, aerr.Error())
			case ecs.ErrCodeInvalidParameterException:
				log.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
			case ecs.ErrCodeClusterNotFoundException:
				log.Println(ecs.ErrCodeClusterNotFoundException, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
			return nil, err
		}
	}

	if len(result.Tasks) == 0 {
		log.Println("No Task found")
		return nil, errors.New("No Task found in the cluster")
	}
	return result.Tasks[0].Attachments[0].Details[1].Value, nil

}
