
## Using Cassowary as an Library

Cassowary can be imported and used as a module in your Go app. Start by fetching the dependency by using go mod:

```bash
$ go mod init test && go get github.com/rogerwelin/cassowary/pkg/client
```

**Example 1: Simple Load Test of an URL**  

```go
package main

import (
        "encoding/json"
	"fmt"

	cassowary "github.com/rogerwelin/cassowary/pkg/client"
)

func main() {
	cass := &cassowary.Cassowary{
		BaseURL:               "http://www.example.com",
		ConcurrencyLevel:      1,
		Requests:              10,
		DisableTerminalOutput: true,
	}
	metrics, err := cass.Coordinate()
	if err != nil {
		panic(err)
	}

        // print results
	fmt.Printf("%+v\n", metrics)

        // or print as json
	jsonMetrics, err := json.Marshal(metrics)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonMetrics))
}
```

**Example 2: Load Test an URL across multiple URL paths**  

The following code will make 30 requests across the 3 URL paths declared in the URLPaths field:

```go
package main

import (
        "encoding/json"
	"fmt"

	cassowary "github.com/rogerwelin/cassowary/pkg/client"
)

func main() {
	cass := &cassowary.Cassowary{
		BaseURL:               "http://www.example.com",
		ConcurrencyLevel:      2,
		Requests:              30,
		FileMode:	       true,
		URLPaths:	       []string{"/accounts", "/orders", "/customers"},
		DisableTerminalOutput: true,
	}
	metrics, err := cass.Coordinate()
	if err != nil {
		panic(err)
	}

        // print results
	fmt.Printf("%+v\n", metrics)

        // or print as json
	jsonMetrics, err := json.Marshal(metrics)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonMetrics))
}
```
