package components

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strings"
)

type HttpServer struct {
	IHttp
	version             string
	port                []string
	Handler             interface{}
	HandlerFuncCallBack func(tvl, obj reflect.Value) gin.HandlerFunc
}

func Pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func NewHttpServer(version string, port ...string) *HttpServer {
	return &HttpServer{version: version, port: port}
}

func GetRoutePath(objName, objFunc string) string {
	return strings.ToLower(objName + "/" + objFunc)
}

func (h *HttpServer) HandlerFuncObj(tvl, obj reflect.Value) gin.HandlerFunc {
	if h.HandlerFuncCallBack != nil {
		return h.HandlerFuncCallBack(tvl, obj)
	}

	return func(c *gin.Context) {
		v := tvl.Call([]reflect.Value{obj, reflect.ValueOf(c)})
		if len(v) != 2 {
			c.JSON(http.StatusNotFound, gin.H{"code": -100, "message": "request param len is error", "data": ""})
			return
		}
		code := v[0].Int()
		if code == 0 {
			c.JSON(http.StatusOK, gin.H{"code": v[0].Interface(), "message": "success", "data": v[1].Interface()})
		} else {
			c.JSON(http.StatusOK, gin.H{"code": v[0].Interface(), "message": v[1].Interface()})
		}
	}
}

func (h *HttpServer) SetHandlerFuncCallback(gh func(tvl, obj reflect.Value) gin.HandlerFunc) {
	h.HandlerFuncCallBack = gh
}

func (h *HttpServer) BindHandler(handler interface{}) {
	h.Handler = handler
}

func (h *HttpServer) Start() error {
	//gin初始化
	r := gin.Default()
	r.GET("/ping", Pong)
	typ := reflect.TypeOf(h.Handler)
	val := reflect.ValueOf(h.Handler)
	//t := reflect.Indirect(val).Type()
	//objectName := t.Name()

	numOfMethod := val.NumMethod()
	for i := 0; i < numOfMethod; i++ {
		method := typ.Method(i)
		r.GET(GetRoutePath(h.version, method.Name), h.HandlerFuncObj(method.Func, val))
		r.POST(GetRoutePath(h.version, method.Name), h.HandlerFuncObj(method.Func, val))
	}
	return r.Run(h.port...) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func (h *HttpServer) Stop() {
}
