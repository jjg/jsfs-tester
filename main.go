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


func main() {

  log.Print("Starting up...")
  
  // TODO: parse arguments for server name, port, concurrency, 
  // number of test runs, verbosity, etc.
  serverName := "localhost"
  serverProto := "http"
  serverPort := 7302
  runs := 5
  concurrency := 1

  // Create a structure to record the results and timing of each test
  runResults := make(map[int]map[string]*TestResult, runs)
  //finalResults := make(map[string]*TestResult, testCount)

  // Perform a check to make sure the arguments point to a valid server
  log.Print("Testing connectivity")
  if pass, duration, err := connectTest(serverProto, serverName, serverPort); pass && err != nil {
    log.Fatal("Configuration error (are the arguments correct?): %s", err)
    os.Exit(1)
  } else {
    log.Print("Connectivity test passed! pass: %v, duration: %d", pass, duration)
  }
  
  // TODO: Execute tests as requested (number, concurrency, etc.)
  // TODO: Include auth (token, key, etc.) in these tests
  fileName := fmt.Sprintf("test-%d.json", time.Now().UnixMilli())
  
  // TODO: Test each JSFS method and record the results
  // TODO: POST a file
  log.Printf("Testing %d runs with concurrency %d", runs, concurrency)

  // TODO: Can't do concurrency until we decide how we want to report on it,
  // (show all results, show average, etc.)
  //c := make(chan int)
  for run := 0; run < runs; run++ {
    //go func() {
      runResults[run] = make(map[string]*TestResult, testCount)
      
      log.Print("Testing POST")
      pass, duration, err := postTest(serverProto, serverName, serverPort, fileName)
      if err != nil {
        log.Printf("Error running POST test: %s", err)
      }
      runResults[run]["post"] = &TestResult{Pass: pass, Duration: duration}
 
      
      // TODO: HEAD the file
      // TODO: GET the file
      // TODO: PUT a chane to the file
      // TODO: DELETE the file

      // TODO: Consider how tests might be written external to this code
    //}()
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
  
  log.Print("All done!")
}
