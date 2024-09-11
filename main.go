package main

import (
  "log"
  "fmt"
  "os"
  "net/http"
  "time"
)

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

func main() {
  log.Print("Starting up...")
  
  // TODO: parse arguments for server name, port, concurrency, 
  // number of requests, verbosity, etc.
  serverName := "localhost"
  serverProto := "http"
  serverPort := 7302


  // Create a structure to record the results and timing of each test
  results := make(map[string]*TestResult,10) // TODO: Use a real size, not a guess

  // Perform a check to make sure the arguments point to a valid server
  log.Print("Testing connectivity")
  if pass, duration, err := connectTest(serverProto, serverName, serverPort); pass && err != nil {
    log.Fatal("Configuration error (are the arguments correct?): %s", err)
    os.Exit(1)
  } else {
    log.Print("connectivity test passed!")
    results["connectivity"] = &TestResult{Pass: pass, Duration: duration}  
  }
  
  
  // TODO: Execute tests as requested (number, concurrency, etc.)
  
  // TODO: Test each JSFS method and record the results
  // TODO: POST a file
  
  // TODO: HEAD the file
  // TODO: GET the file
  // TODO: PUT a chane to the file
  // TODO: DELETE the file

  // TODO: Consider how tests might be written external to this code

  // Display results
  fmt.Printf("Test\t\t\tPass\tDuration\n")
  for k, v := range results {
    fmt.Printf("%s\t\t%v\t%dms\n", k, v.Pass, v.Duration)
  }
  
  log.Print("All done!")
}
