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

	Record = os.Getenv("HOSTHEADER")
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
	// Verify Env
	if Record == "" {
		log.Fatalln("Missing Environment Variable HOSTHEADER record set")
		os.Exit(1)
	}

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

	pubIP, err := (&Ec2Handler{
		Service: Ec2Client(os.Getenv("AWS_DEFAULT_REGION")),
		Eni:     eni,
	}).PublicIp()
	if err != nil {
		log.Fatalf("Failed to get PublicIP: %v", err)
	}
	log.Printf("Public IP: %v", *pubIP)

	// Create route53 Record set

}
