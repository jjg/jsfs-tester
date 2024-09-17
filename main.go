package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

// TODO: this should probably be more dynamic
const testCount = 6

type TestRequest struct {
	Protocol   string
	ServerName string
	ServerPort int
	FileName   string
}

type TestResult struct {
	Pass     bool
	Duration int64
}

func doRun(run int, q chan int, r chan<- map[string]*TestResult, tr *TestRequest) {
	runResult := make(map[string]*TestResult, testCount)

	// Add the run to the queue (this will block if the channel is full)
	q <- run

	// Test each JSFS method and record the results
	log.Print("Testing POST")
	pass, duration, err := postTest(tr)
	if err != nil {
		log.Printf("Error running POST test: %s", err)
	}
	runResult["post"] = &TestResult{Pass: pass, Duration: duration}

	log.Print("Testing HEAD")
	pass, duration, err = headTest(tr)
	if err != nil {
		log.Printf("Error running HEAD test: %s", err)
	}
	runResult["head"] = &TestResult{Pass: pass, Duration: duration}

	log.Print("Testing GET")
	pass, duration, err = getTest(tr)
	if err != nil {
		log.Printf("Error running GET test: %s", err)
	}
	runResult["get"] = &TestResult{Pass: pass, Duration: duration}

	log.Print("Testing PUT")
	pass, duration, err = putTest(tr)
	if err != nil {
		log.Printf("Error running PUT test: %s", err)
	}
	runResult["put"] = &TestResult{Pass: pass, Duration: duration}

	log.Print("Testing DELETE")
	pass, duration, err = deleteTest(tr)
	if err != nil {
		log.Printf("Error running DELETE test: %s", err)
	}
	runResult["delete"] = &TestResult{Pass: pass, Duration: duration}

	// TODO: Include auth (token, key, etc.) in these tests
	// TODO: Consider how tests might be written external to this code

	// Send the results of the run through the results channel
	r <- runResult

	// Remove this run from the queue
	log.Printf("Removing run %d from the queue", <-q)
}

func printRunResult(r map[string]*TestResult) {
	// TODO: It would be cool if these always printed in the same order
	fmt.Printf("Test\t\tPass\tDuration\n")
	for k, v := range r {
		fmt.Printf("%s\t\t%v\t%dms\n", k, v.Pass, v.Duration)
	}
}

func main() {

	// Log to a file
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Unable to set logfile:", err.Error())
	}
	log.SetOutput(file)

	log.Print("Starting up...")

	// TODO: parse cli arguments for server name, port,
	// concurrency, number of test runs, verbosity, etc.
	serverName := "localhost"
	serverProto := "http"
	serverPort := 7302
	runs := 5
	concurrency := 1

	// Perform a check to make sure the arguments point to a valid server
	log.Print("Testing connectivity")
	if pass, duration, err := connectTest(serverProto, serverName, serverPort); pass && err != nil {
		log.Fatal("Configuration error (are the arguments correct?): %s", err)
		os.Exit(1)
	} else {
		log.Printf("Connectivity test passed! pass: %v, duration: %d", pass, duration)
	}

	log.Printf("Testing %d runs with concurrency of %d", runs, concurrency)

	// Measure overall test duration
	startTime := time.Now()

	// Create a structure to record the results and timing of each test
	runResults := make(map[int]map[string]*TestResult, runs)

	// Generate a unique filename to use for testing
	// TODO: Consider making it unique-er
	fileName := fmt.Sprintf("test-%d.json", time.Now().UnixMilli())

	// Assemble the request parameters
	testRequest := &TestRequest{
		Protocol:   serverProto,
		ServerName: serverName,
		ServerPort: serverPort,
		FileName:   fileName,
	}

	// Create a channel to serve as a concurrency queue
	q := make(chan int, concurrency)

	// Create a channel to receive run results
	r := make(chan map[string]*TestResult)

	// Start each run
	for run := 0; run < runs; run++ {
		go doRun(run, q, r, testRequest)
	}

	// Gather run results from the channel
	// TODO: There may be a smarter way to do this
	for i := runs - 1; i >= 0; i-- {
		runResults[i] = <-r
		printRunResult(runResults[i])
	}

	// Display results
	methodAccumulator := make(map[string]int64, 5)

	for i := 0; i < runs; i++ {
		fmt.Printf("Run %d:\n", i+1)
		printRunResult(runResults[i])

		// Accumulate results to display averages
		for k, v := range runResults[i] {
			methodAccumulator[k] = methodAccumulator[k] + v.Duration
		}
	}

	// Display overall averages for each method
	fmt.Println("Average duration by method for all runs:")
	fmt.Println("Method\tDuration")
	for k, v := range methodAccumulator {
		fmt.Printf("%s:\t%dms\n", k, v/int64(runs))
	}

	fmt.Printf("%d test runs with concurrency of %d finished in %dms\n", runs, concurrency, time.Since(startTime).Milliseconds())

	log.Print("All done!")
}
