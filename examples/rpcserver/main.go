package main

import "github.com/osgochina/dmicro/server"

func main() {
	s := server.NewRpcServer("test_one",
		server.OptEnableHeartbeat(true),
		server.OptListenAddress("127.0.0.1:8199"),
	)
	s.ListenAndServe()
}
