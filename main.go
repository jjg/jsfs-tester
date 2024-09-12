package main

import (
  "log"
  "fmt"
  "os"
  "net/http"
  "time"
  "bytes"
)

// TODO: this should probably be more dynamic
const testCount = 6

type TestResult struct {
  Pass bool
  Duration int64
}

func connectTest(proto string, name string, port int) (bool, int64, error){
  startTime := time.Now()

  reqUrl := fmt.Sprintf("%s://%s:%d", proto, name, port)
  req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
  if err != nil {
    return false, time.Since(startTime).Milliseconds(), err
  }
  res, err := http.DefaultClient.Do(req)
  if err != nil {
    return false, time.Since(startTime).Milliseconds(), err
  }

  log.Printf("Got a response code (%d), server is up!", res.StatusCode)
  
  // if we got this far, success!
  return true, time.Since(startTime).Milliseconds(), nil
}

func postTest(proto string, name string, port int, fileName string) (bool, int64, error){
  startTime := time.Now()

  // TODO: Don't hard-code the access key
  reqUrl := fmt.Sprintf("%s://%s:%d/%s?access_key=foo", proto, name, port, fileName)
  jsonBody := []byte(`{"message": "I am a teapot"}`)
  bodyReader := bytes.NewReader(jsonBody)
  
  req, err := http.NewRequest(http.MethodPost, reqUrl, bodyReader)
  if err != nil {
    return false, time.Since(startTime).Milliseconds(), err
  }
  res, err := http.DefaultClient.Do(req)
  if err != nil {
    return false, time.Since(startTime).Milliseconds(), err
  }

  log.Printf("POST got a response code %d", res.StatusCode)

  // TODO: Might need to check for variations of 2xx
  if res.StatusCode == 200 {
    return true, time.Since(startTime).Milliseconds(), nil
  } else {
    return false, time.Since(startTime).Milliseconds(), nil
  }
}

// TODO: This should probably accept a struct or something easier
// for the caller to not mess-up
func doRun(run int, q chan int, r chan<- map[string]*TestResult, serverProto string, serverName string, serverPort int, fileName string){
  runResult := make(map[string]*TestResult, testCount)

  // Add the run to the queue (this will block if the channel is full)
  q <- run

  // Test each JSFS method and record the results    
  log.Print("Testing POST")
  pass, duration, err := postTest(serverProto, serverName, serverPort, fileName)
  if err != nil {
    log.Printf("Error running POST test: %s", err)
  }
  runResult["post"] = &TestResult{Pass: pass, Duration: duration}

  // TODO: HEAD the file
  // TODO: GET the file
  // TODO: PUT a chane to the file
  // TODO: DELETE the file

  // TODO: Include auth (token, key, etc.) in these tests
  // TODO: Consider how tests might be written external to this code

  // Send the results of the run through the results channel
  r <- runResult

  // Remove this run from the queue
  log.Printf("Removing run %d from the queue", <-q)  
}

func main() {

  log.Print("Starting up...")
  
  // TODO: parse cli arguments for server name, port, 
  // concurrency, number of test runs, verbosity, etc.
  serverName := "localhost"
  serverProto := "http"
  serverPort := 7302
  runs := 5         // TODO: there is a bug when this is set to 1
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
  
  // Create a channel to serve as a concurrency queue
  q := make(chan int, concurrency)
  
  // Create a channel to receive run results
  r := make(chan map[string]*TestResult)

  // Start each run
  for run := 0; run < runs; run++ {

    go doRun(run, q, r, serverProto, serverName, serverPort, fileName)
    /*
    go func() {
      runResult := make(map[string]*TestResult, testCount)
      
      // Add the run to the queue (this will block if the channel is full)
      q <- run

      // Test each JSFS method and record the results    
      log.Print("Testing POST")
      pass, duration, err := postTest(serverProto, serverName, serverPort, fileName)
      if err != nil {
        log.Printf("Error running POST test: %s", err)
      }
      runResult["post"] = &TestResult{Pass: pass, Duration: duration}
 
      // TODO: HEAD the file
      // TODO: GET the file
      // TODO: PUT a chane to the file
      // TODO: DELETE the file

      // TODO: Include auth (token, key, etc.) in these tests
      // TODO: Consider how tests might be written external to this code

      // Send the results of the run through the results channel
      r <- runResult

      // Remove this run from the queue
      log.Printf("Removing run %d from the queue", <-q)
    }()
    */
  }

  // Gather run results from the channel
  // TODO: There may be a smarter way to do this
  for i:=runs;i>0;i-- {
    runResults[i] = <-r
    
    // TODO: Consider outputting the results of each run as they complete
    // instead of all at once at the end of the program
  }

  // Display results
  for i:=0;i<runs;i++ {
    fmt.Printf("Run %d:\n", i)
    fmt.Printf("Test\t\tPass\tDuration\n")
    for k, v := range runResults[i] {
      fmt.Printf("%s\t\t%v\t%dms\n", k, v.Pass, v.Duration)
    }
  }

  // TODO: Display overall average results
  fmt.Printf("Total testing duration: %dms\n", time.Since(startTime).Milliseconds())
  
  log.Print("All done!")
}
