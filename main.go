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
	FargateMetadataEndpoint = "http://169.254.170.2/v2/metadata"
	// FargateMetadataEndpoint = "http://localhost:8080/v2/metadata"

	httpClient = &http.Client{Timeout: 10 * time.Second}
	region     = os.Getenv("AWS_DEFAULT_REGION")
	recordName = os.Getenv("HOSTHEADER")
	hostZoneID = os.Getenv("HOSTZONEID")
)

// TaskData represents ARNs required to describe an ECS Fargate task
type TaskData struct {
	Cluster string `json:"Cluster"`
	TaskARN string `json:"TaskARN"`
}

func responseJSON(url string, target interface{}) error {
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func main() {
	// Verify Env
	if recordName == "" {
		log.Fatalln("Missing Environment Variable HOSTHEADER record set")
		os.Exit(1)
	}
	if hostZoneID == "" {
		log.Fatalln("Missing Environment Variable HostZoneID record set")
		os.Exit(1)
	}

	// Get Task
	td := new(TaskData)
	responseJSON(FargateMetadataEndpoint, td)
	log.Println(td.Cluster)
	log.Println(td.TaskARN)

	// Get Task ENI
	eni, err := (&EcsHandler{
		Service: EcsClient(region),
		Cluster: &td.Cluster,
		TaskArn: &td.TaskARN,
	}).TaskEni()
	if err != nil {
		log.Fatalf("Failed to get ECS task ENI: %v", err)
	}
	log.Printf("Task ENI provisioned: %v", *eni)

	// Get Public IP

	pubIP, err := (&Ec2Handler{
		Service: Ec2Client(region),
		Eni:     eni,
	}).PublicIp()
	if err != nil {
		log.Fatalf("Failed to get PublicIP: %v", err)
	}
	log.Printf("Public IP: %v", *pubIP)

	// Create route53 Record set

	err = (&DNSHandler{
		Service:    DNSClient(region),
		RecordName: &recordName,
		HostZoneID: &hostZoneID,
		PubIP:      pubIP,
	}).RecordSet()

	if err != nil {
		log.Printf("Error creating Route53 record: %v", err)
	}
}
