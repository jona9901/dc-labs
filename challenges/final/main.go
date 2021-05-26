package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

//	"github.com/jona9901/dc-labs/challenges/third-partial/api"
	"github.com/jona9901/dc-labs/challenges/third-partial/controller"
	"github.com/jona9901/dc-labs/challenges/third-partial/scheduler"
)

func main() {
	log.Println("Welcome to the Distributed and Parallel Image Processing System")

	// Start Controller
	go controller.Start()

	// Start Scheduler
	jobs := make(chan scheduler.Job)
	go scheduler.Start(jobs)
	// Send sample jobs
	sampleJob := scheduler.Job{Address: "localhost:50051", RPCName: "hello"}

	for {
		sampleJob.RPCName = fmt.Sprintf("hello-%v", rand.Intn(10000))
		jobs <- sampleJob
		time.Sleep(time.Second * 5)
	}

	// API
	// Here's where your API setup will be

	//go api.Start()
}
