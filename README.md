# Package postgresql-generator
Build on top of `github.com/icrowley/fake` and `github.com/lib/pq`
For generating random data that are specified by `pgtype`

# Environment
```
export PG_HOST=localhost
export PG_DB=postgresql
export PG_USER=pguser
export PG_PASS=p@ssw0rd
```

# Example
```go
package main

import (
	"fmt"
	"log"
	"time"

	pg "github.com/hieutrtr/postgresql-generator"
)

type Users struct {
	name     string `pgtype:"varchar"`
	age      uint8  `pgtype:"smallint"`
	friends  int    `pgtype:"int"`
	salary   int32  `pgtype:"money"`
	ipv4     string `pgtype:"cidr"`
	metadata string `pgtype:"jsonb"`
}

const numQueries = 10000
const numWorkers = 20
const numConnection = 50

func main() {
	sup := make(chan string, numQueries)
	res := make(chan error, numQueries)
	for c := 0; c < numConnection; c++ {
		go func() {
			pgCtrl, err := pg.NewPG()
			if err != nil {
				panic(err.Error())
			}
			c := &Config{}
			if err = c.parse(); err != nil {
				panic(err.Error())
			}
			if err = pgCtrl.connect(c); err != nil {
				panic(err.Error())
			}

			for w := 0; w < numWorkers; w++ {
				go func() {
					for q := range sup {
						err := pgCtrl.Execute(q)
						if err != nil {
							res <- fmt.Errorf("%s %s\n", q, err.Error())
						}
						res <- nil
					}
				}()
			}
		}()
	}
	start := time.Now()
	pg.NewSupplier(&Users{}, numQueries, sup)
	for i := 0; i < numQueries; i++ {
		r := <-res
		if r != nil {
			fmt.Println(r)
		}
	}
	elapsed := time.Since(start)
	log.Printf("Insertion took %s", elapsed)
}
```
