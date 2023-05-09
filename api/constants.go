package api

type Method string

const (
	POST   Method = "Post"
	GET    Method = "Get"
	DELETE Method = "Delete"
	PUT    Method = "Put"
	PATCH  Method = "Patch"
)

func (method Method) ToString() string {
	return string(method)
}
