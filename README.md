# pu55y
`don't think too much, it's just a base package for go projects`

### Can & Can't
- support configuration file (but only json now)
- support `Restful Api`
- support multiple datasource
- And more in future ...
- **Not** Support discovery services like Nacos or some kind of spring projects because I don't want it to become too heavy, however, you can still import those form other packages.

### Properties
- it's configurable using configuration file, see `conf.json`

### Support Stacks
- redis
- postgresql



## How To Use
- import & download `change @x.x.x to the version you need`
```shell
go get github.com/real-uangi/pu55y@x.x.x
```
- edit main function, see `test/run_test.go`
```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/real-uangi/pu55y/api"
	"github.com/real-uangi/pu55y/runner"
)

func main() {
	
	r := runner.Prepare()
	
	r.AddApi(api.GET, "/example", func(context *gin.Context) {
		//your function here
	})
	//...
	
	r.Run()
	
}
```
- database operations
```go
package db

import "github.com/real-uangi/pu55y/datasource"

func dbQuery() {
	
	ds := datasource.Get(xx) //xx : database name in your config
	
	rows, err := ds.Query(xxxxx)
	defer rows.Close()
	for rows.Next() {
		//... original database/sql operations in go sdk
	}
	
}
```

- redis operations

```go
package redis

import "github.com/real-uangi/pu55y/rdb"

var ctx = context.Background()

func redisOps() {

	// provided functions
	rdb.Set(x,xx)
	rdb.TryLock(x,xx,xxx)
	//...
	
	//for origin redis/v9 operations
	client := rdb.GetClient()
	client.Ping(ctx)
	//...

}
```