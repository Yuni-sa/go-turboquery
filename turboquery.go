package turboquery

import (
	"context"
	"database/sql"
	"log"
)

type ReplicaInfo struct {
	Name string
	DSN  string
}

type Conn struct {
	Name     string
	Endpoint *sql.DB
}

type Result struct {
	DatabaseName string
	Columns      []string
	Rows         [][]string
}

func MultiQuery(replicas []ReplicaInfo, query string) Result {
	conns := make([]Conn, len(replicas))

	for i, replica := range replicas {
		db, err := sql.Open("mysql", replica.DSN)
		if err != nil {
			log.Fatalf("failed to connect to %s: %v", replica.Name, err)
		}
		defer db.Close()

		conns[i] = Conn{
			Name:     replica.Name,
			Endpoint: db,
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan Result, len(conns))
	done := make(chan struct{})

	defer func() {
		cancel() // Ensure cancel is always called before returning
		close(done)
	}()

	for _, conn := range conns {
		go func(conn Conn) {
			select {
			case ch <- Query(ctx, conn, query):
				cancel() // Cancel other queries
			case <-done:
				// Do nothing if a result has already been received
			}
		}(conn)
	}

	result := <-ch

	return result
}

func Query(ctx context.Context, c Conn, query string) Result {
	rows, err := c.Endpoint.QueryContext(ctx, query)
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Panic(err)
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var rowsData [][]string

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Panic(err)
		}

		var rowData []string
		for _, value := range values {
			if value == nil {
				rowData = append(rowData, "")
			} else {
				rowData = append(rowData, string(value))
			}
		}
		rowsData = append(rowsData, rowData)
	}

	if err := rows.Err(); err != nil {
		log.Panic(err)
	}

	return Result{DatabaseName: c.Name, Columns: columns, Rows: rowsData}
}
