# Package postgresql_generator
Build on top of `github.com/icrowley/fake` and `github.com/lib/pq`
For generating random data that are specified by `gentype`

# Environment
```
export PG_HOST=localhost
export PG_DB=postgresql
export PG_USER=pguser
export PG_PASS=p@ssw0rd
```

# Example

## PostgreSQL
```go
package main

import (
	"fmt"
	"log"
	"time"

	pg "github.com/hieutrtr/postgresql-generator"
)

// Users should be matched with table in postgresql server
// CREATE TABLE users (
//     user_id   serial,
//     name      varchar,
//     age       smallint,
//     friends   int,
//     salary    money,
//     ipv4      inet,
//     metadata  jsonb,
//     CONSTRAINT user_id PRIMARY KEY(user_id)
// );
type Users struct {
	name     string `gentype:"varchar"`
	age      uint8  `gentype:"smallint"`
	friends  int    `gentype:"int"`
	salary   int32  `gentype:"money"`
	ipv4     string `gentype:"cidr"`
	metadata string `gentype:"jsonb"`
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
			c := &pg.Config{}
			if err = c.Parse(); err != nil {
				panic(err.Error())
			}
			if err = pgCtrl.Connect(c); err != nil {
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

	go func(st interface{}, num int, sup chan<- string) {
		for i := 0; i < num; i++ {
			sup <- pg.GenInsertion(st)
		}
		close(sup)
	}(&Users{}, numQueries, sup)

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

## Kafka

```go
package main

import (
	"fmt"
	"log"
	"time"

	pg "github.com/hieutrtr/postgresql-generator"
)

type Ads struct {
	title string `gentype:"varchar"`
	ad_id uint8  `gentype:"smallint"`
}

const numQueries = 10000
const numWorkers = 20
const numConnection = 50

func main() {
	sup := make(chan string, numQueries)
	res := make(chan error, numQueries)
	for c := 0; c < numConnection; c++ {
		go func() {
			for w := 0; w < numWorkers; w++ {
				go func() {
					for q := range sup {
						pro, err := pg.NewProducer()
						if err != nil {
							panic(err.Error())
						}
						e := &pg.UploadEvent{
							Topic:   "test_kafka",
							Payload: q,
						}
						err = pro.Produce(e)
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

	go func(st interface{}, num int, sup chan<- string) {
		for i := 0; i < num; i++ {
			sup <- pg.GenJSON(st)
		}
		close(sup)
	}(&Ads{}, numQueries, sup)

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
