package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type TaskData struct {
	Cluster    string            `json:"Cluster"`
	TaskARN    string            `json:"TaskARN"`
	Family     string            `json:"Family"`
	Containers map[string]string `json:"Containers"`
}

type Health struct {
	Status string `json:"status"`
}

func userHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	t := TaskData{
		Cluster: "arn:aws:ecs:ap-southeast-2:567418462583:cluster/ecs-cluster-fin-exp-dev",
		TaskARN: "arn:aws:ecs:ap-southeast-2:567418462583:task/ecs-cluster-fin-exp-dev/dfc8752c57134e17afee8696be98cf78",
		Family:  "myapp",
		Containers: map[string]string{
			"DockerId": "870b6c89b84778963577874678b34edbb1adac7a739fb33ce1ab39af9526be46",
		},
	}

	json.NewEncoder(w).Encode(t)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	health := Health{
		Status: "ok",
	}
	json.NewEncoder(w).Encode(health)
}

func main() {
	http.HandleFunc("/v2/metadata", userHandler)
	http.HandleFunc("/healthz", healthHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
