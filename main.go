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

	"github.com/lemonyxk/kitty/v2"
	"github.com/lemonyxk/kitty/v2/router"
	"github.com/lemonyxk/kitty/v2/socket"
	"github.com/lemonyxk/kitty/v2/socket/http"
	"github.com/lemonyxk/kitty/v2/socket/http/server"
	server2 "github.com/lemonyxk/kitty/v2/socket/websocket/server"
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

var (
	maxCount        = 600
	data            []update
	interval        = time.Millisecond * 500
	ip              = "0.0.0.0"
	port            = 23456
	httpServer      = &server.Server{Addr: fmt.Sprintf("%s:%d", ip, port)}
	webSocketServer = &server2.Server{Addr: fmt.Sprintf("%s:%d", ip, port+1), Path: "/debug/feed/"}
	lastPause       uint32
)

func Interval(t time.Duration) {
	interval = t
}

func Port(p int) {
	port = p
}

func Ip(h string) {
	ip = h
}

func MaxCount(n int) {
	maxCount = n
}

func Start() {

	var httpServerRouter = kitty.NewHttpServerRouter()

	httpServer.Use(func(next server.Middle) server.Middle {
		return func(stream *http.Stream) {
			if stream.Request.Header.Get("Upgrade") == "websocket" {
				httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: fmt.Sprintf("%s:%d", ip, port+1)}).ServeHTTP(stream.Response, stream.Request)
				return
			}
			next(stream)
		}
	})

	httpServerRouter.Group("/debug").Handler(func(handler *router.Handler[*http.Stream]) {
		handler.Get("/charts/").Handler(func(stream *http.Stream) error {
			return stream.EndString(render())
		})
	})

	go httpServer.SetRouter(httpServerRouter).Start()

	var debugUrl = fmt.Sprintf("http://%s:%d/debug/charts/", ip, port)

	fmt.Printf("you can open %s to watch.\n", debugUrl)

	var webSocketServerRouter = kitty.NewWebSocketServerRouter()

	webSocketServerRouter.Group("/debug").Handler(func(handler *router.Handler[*socket.Stream[server2.Conn]]) {
		handler.Route("/login").Handler(func(stream *socket.Stream[server2.Conn]) error {
			return stream.Conn.JsonEmit("listen", http.JsonFormat{Status: "SUCCESS", Code: 200, Msg: data})
		})
	})

	webSocketServer.OnOpen = func(conn server2.Conn) {}
	webSocketServer.OnClose = func(conn server2.Conn) {}
	webSocketServer.OnError = func(err error) {}

	go webSocketServer.SetRouter(webSocketServerRouter).Start()

	go gatherData()

}

func gatherData() {

	var ticker = time.NewTicker(interval)

	for range ticker.C {
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

		webSocketServer.JsonEmitAll("listen", http.JsonFormat{Status: "SUCCESS", Code: 200, Msg: []update{u}})
	}
}
