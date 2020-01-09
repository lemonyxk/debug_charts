/**
* @program: go
*
* @description:
*
* @author: lemo
*
* @create: 2020-01-06 20:45
**/

package debug_charts

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/Lemo-yxk/lemo"
	"github.com/Lemo-yxk/lemo/console"
	"github.com/Lemo-yxk/lemo/exception"
)

type update struct {
	Ts             int64
	BytesAllocated uint64
	GcPause        uint64
	Block          int
	Goroutine      int
	Heap           int
	Mutex          int
	ThreadCreate   int
}

const (
	maxCount int = 600
)

var (
	data            []update
	lastPause       uint32
	interval        = time.Millisecond * 500
	host            = "0.0.0.0"
	port            = 23456
	httpServer      = &lemo.HttpServer{Host: host, Port: port, AutoBind: true}
	webSocketServer = &lemo.WebSocketServer{Host: host, Port: port + 1, Path: "/debug/feed/", AutoBind: true}
)

func Interval(t time.Duration) {
	interval = t
}

func init() {

	var httpServerRouter = &lemo.HttpServerRouter{IgnoreCase: true}

	httpServer.Use(func(next lemo.HttpServerMiddle) lemo.HttpServerMiddle {
		return func(stream *lemo.Stream) {
			if stream.Request.Header.Get("Upgrade") == "websocket" {
				httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: fmt.Sprintf("%s:%d", host, port+1)}).ServeHTTP(stream.Response, stream.Request)
				return
			}
			next(stream)
		}
	})

	httpServerRouter.Group("/debug").Handler(func(handler *lemo.HttpServerRouteHandler) {
		handler.Get("/charts/").Handler(func(stream *lemo.Stream) exception.ErrorFunc {
			return exception.New(stream.EndString(render()))
		})
	})

	go httpServer.SetRouter(httpServerRouter).Start()

	var debugUrl = fmt.Sprintf("http://%s:%d/debug/charts/", host, port)

	console.Printf("you can open %s to watch.\n", debugUrl)

	var webSocketServerRouter = &lemo.WebSocketServerRouter{IgnoreCase: true}

	webSocketServerRouter.Group("/debug").Handler(func(handler *lemo.WebSocketServerRouteHandler) {
		handler.Route("/login").Handler(func(conn *lemo.WebSocket, receive *lemo.Receive) exception.ErrorFunc {
			return conn.JsonFormat(lemo.JsonPackage{Event: "listen", Message: lemo.JM("SUCCESS", 200, data)})
		})
	})

	webSocketServer.OnOpen = func(conn *lemo.WebSocket) {}
	webSocketServer.OnClose = func(fd uint32) {}
	webSocketServer.OnError = func(err exception.ErrorFunc) {}

	go webSocketServer.SetRouter(webSocketServerRouter).Start()

	go gatherData()

}

func gatherData() {

	nowUnix := time.Now().Unix()

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	u := update{
		Ts:           nowUnix * 1000,
		Block:        pprof.Lookup("block").Count(),
		Goroutine:    pprof.Lookup("goroutine").Count(),
		Heap:         pprof.Lookup("heap").Count(),
		Mutex:        pprof.Lookup("mutex").Count(),
		ThreadCreate: pprof.Lookup("threadcreate").Count(),
	}

	u.BytesAllocated = ms.Alloc

	if lastPause == 0 || lastPause != ms.NumGC {
		gcPause := ms.PauseNs[(ms.NumGC+255)%256]
		u.GcPause = gcPause
		lastPause = ms.NumGC
	}

	data = append(data, u)

	if len(data) > maxCount {
		data = data[len(data)-maxCount:]
	}

	time.Sleep(interval)

	webSocketServer.JsonFormatAll(lemo.JsonPackage{Event: "listen", Message: lemo.JM("SUCCESS", 200, []update{u})})

	go gatherData()
}
