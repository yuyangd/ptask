package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	// FargateMetadataEndpoint refer to https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-metadata-endpoint-v2.html
	// FargateMetadataEndpoint = "http://169.254.170.2/v2/metadata"
	FargateMetadataEndpoint = "http://localhost:8080/v2/metadata"

	httpClient = &http.Client{Timeout: 10 * time.Second}
)

type TaskData struct {
	Cluster string `json:"Cluster"`
	TaskARN string `json:"TaskARN"`
}

func responseJson(url string, target interface{}) error {
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func main() {
	// Get Task
	td := new(TaskData)
	responseJson(FargateMetadataEndpoint, td)
	log.Println(td.Cluster)
	log.Println(td.TaskARN)

	// Get Task ENI
	eni, err := (&EcsHandler{
		Service: EcsClient(os.Getenv("AWS_DEFAULT_REGION")),
		Cluster: &td.Cluster,
		TaskArn: &td.TaskARN,
	}).TaskEni()
	if err != nil {
		log.Fatalf("Failed to get ECS task ENI: %v", err)
	}
	log.Printf("Task ENI provisioned: %v", *eni)

	// Get Public IP

	// If public IP not found, attach an EIP

	//
}
