package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type ServiceStatus struct {
	Name   string
	Status string
	URL    string
}

type Operator struct {
	services []ServiceStatus
	client   *http.Client
}

func NewOperator() *Operator {
	return &Operator{
		services: []ServiceStatus{
			{
				Name:   "chat-server",
				Status: "unknown",
				URL:    "http://chat-server:8080/health",
			},
			{
				Name:   "moderation-service",
				Status: "unknown",
				URL:    "http://moderation-service:5001/health",
			},
		},
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (o *Operator) Start() {
	log.Println("Starting operator...")
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Initial check
	o.checkServices()

	// Periodic checks using for range
	for range ticker.C {
		o.checkServices()
	}
}

func (o *Operator) checkServices() {
	for i, service := range o.services {
		status := o.checkServiceHealth(service.URL)
		if status != service.Status {
			o.services[i].Status = status
			log.Printf("Service %s status changed to: %s\n", service.Name, status)
		}
	}
}

func (o *Operator) checkServiceHealth(url string) string {
	resp, err := o.client.Get(url)
	if err != nil {
		log.Printf("Error checking %s: %v\n", url, err)
		return "unhealthy"
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return "healthy"
	}

	log.Printf("Unhealthy status code from %s: %d\n", url, resp.StatusCode)
	return "unhealthy"
}

func main() {
	fmt.Println("Chat System Operator v0.1")
	operator := NewOperator()
	operator.Start()
}


