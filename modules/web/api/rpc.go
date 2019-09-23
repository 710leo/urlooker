package api

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/710leo/urlooker/modules/web/g"
)

type Web int

func Start() {
	addr := g.Config.Rpc.Listen

	server := rpc.NewServer()
	server.Register(new(Web))

	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatalln("listen occur error", e)
	} else {
		log.Println("listening on", addr)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("listener accept occur error", err)
			continue
		}
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

func (this *Web) Ping(req interface{}, reply *string) error {
	*reply = "ok"
	return nil
}
