package receiver

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/710leo/urlooker/modules/alarm/g"
)

func Start() {
	addr := g.Config.Rpc.Listen

	server := rpc.NewServer()
	server.Register(new(Alarm))

	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatalln("listen occur error", e)
	} else {
		log.Println("listen on..", addr)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("listener accept error:", err)
			continue
		}
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
