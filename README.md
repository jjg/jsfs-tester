# jsfsmoke

Automated smoke tester for JSFS-compatible servers.

## Usage

Running `jsfsmoke` with no arguments will perform a single test against a local JSFS server and report some basic statistics.  

### Output
```
Average duration by method for all runs:
Method	Duration
head:	1ms
get:	2ms
put:	3ms
delete:	10ms
post:	2ms
1 test runs with concurrency of 1 finished in 20ms
```

In addition to the console output, `jsfsmoke` keeps a more detailed record of the test in the file `log.txt`.

```
2024/09/17 15:35:02 Starting up...
2024/09/17 15:35:02 Testing connectivity
2024/09/17 15:35:02 Got a response code (404), server is up!
2024/09/17 15:35:02 Connectivity test passed! pass: true, duration: 2
2024/09/17 15:35:02 Testing 1 runs with concurrency of 1
2024/09/17 15:35:02 Testing POST
2024/09/17 15:35:02 POST got a response code 200
2024/09/17 15:35:02 Testing HEAD
2024/09/17 15:35:02 HEAD got a response code 200
2024/09/17 15:35:02 Testing GET
2024/09/17 15:35:02 GET got a response code 200
2024/09/17 15:35:02 Testing PUT
2024/09/17 15:35:02 POST got a response code 200
2024/09/17 15:35:02 Testing DELETE
2024/09/17 15:35:02 DELETE got a response code 204
2024/09/17 15:35:02 HEAD got a response code 404
2024/09/17 15:35:02 Removing run 0 from the queue
2024/09/17 15:35:02 All done!
```

### Options

`jsfsmoke`'s' behavior can be customized by providing additional arguments.

```
   -c int
     	Number of tests to run at the same time (concurrency) (default 1)
   -n int
     	Number of test runs to execute (default 1)
   -name string
     	Name or IP address of the JSFS server to test (default "localhost")
   -port int
     	Port number the JSFS server is listening on (default 7302)
   -protocol string
     	Protocol to use for requests (http or https) (default "http")
   -v	Increase verbosity
```


