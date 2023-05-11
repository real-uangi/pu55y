// Package api @author uangi 2023-05
package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/real-uangi/pu55y/plog"
)

type restfulApi struct {
	method    Method
	uri       string
	processor gin.HandlerFunc
}

type Server struct {
	apiMap    map[string]restfulApi
	logEnable bool
	port      string
}

// Run 启动
func (server *Server) Run() error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	//http日志开关
	if true == server.logEnable {
		formatter := func(param gin.LogFormatterParams) string {
			var msg = fmt.Sprintf("[%d]%s takes %d ", param.StatusCode, param.Path, param.Latency)
			return plog.GetLine(plog.LvInfo, msg, param.TimeStamp)
		}
		router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
			Formatter: formatter,
			Output:    nil,
			SkipPaths: nil,
		}))
	}
	for _, api := range server.apiMap {
		switch api.method {
		case GET:
			{
				router.GET(api.uri, api.processor)
				break
			}
		case POST:
			{
				router.POST(api.uri, api.processor)
				break
			}
		case DELETE:
			{
				router.DELETE(api.uri, api.processor)
				break
			}
		case PUT:
			{
				router.PUT(api.uri, api.processor)
				break
			}
		case PATCH:
			{
				router.PATCH(api.uri, api.processor)
				break
			}
		}
	}
	return router.Run(":" + server.port)
}

// AddApi 添加API
func (server *Server) AddApi(method Method, uri string, processor gin.HandlerFunc) {
	var api restfulApi
	api.uri = uri
	api.method = method
	api.processor = processor
	if server.apiMap == nil {
		server.apiMap = make(map[string]restfulApi)
	}
	server.apiMap[method.ToString()+uri] = api
}

// ListenPort 监听端口
func (server *Server) ListenPort(port string) {
	server.port = port
	plog.Info("Starting server listen on port " + port)
}

// SetLog 开启日志
func (server *Server) SetLog(enable bool) {
	server.logEnable = enable
}
