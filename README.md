# pu55y
`don't think too much, it's just a base package for go projects`

### Can & Can't
- support configuration file (but only json now)
- support `Restful Api`
- And more in development
- Not Support discovery services like Nacos or some kind of spring projects because I don't want it to become too heavy, but you can still import those form other packages.

### Properties
- it's configurable using configuration file, see `conf.json`

### Support Stacks
- redis
- postgresql



## How To Use
- import & download 
```shell
go get github.com/real-uangi/pu55y
```
- edit main function, see `test/run_test.go`
```go
package main

func main() {

	server := runner.Prepare()
	
	server.AddApi(api.GET, "/example", func(context *gin.Context) {
		//your function here
	})
	//...
	runner.Run()
	
}
```