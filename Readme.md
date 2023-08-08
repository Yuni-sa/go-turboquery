# turboquery

turboquery is a Go package that allows you to query multiple databases or replicas with the same query and returns the fastest result while canceling the remaining queries. It helps you find the fastest queries among multiple database sources if you don't care about extra load to every source and replication lag.

## Installation

To use turboquery in your Go project, you need to have Go installed and set up. Then, you can install the package using the following command:

```shell
go get github.com/Yuni-sa/go-turboquery
```

## Usage

Here's an example of how to use turboquery with a mysql cluster:

```go
package main

import (
	"fmt"

	"github.com/Yuni-sa/go-turboquery"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	replicas := []turboquery.ReplicaInfo{
		{
			Name: "replica1",
			DSN:  "replica1_connection_string",
		},
		{
			Name: "replica2",
			DSN:  "replica2_connection_string",
		},
		// Add more replicas as needed
	}

	query := "SELECT * FROM your_table"
	result := turboquery.MultiQuery(replicas, query)

	fmt.Println(result)
}
```