package main

import (
	flag "flag"
	net "net"
	fmt "fmt"
	log "log"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", lis.Addr())
	}
	fmt.Println("Hello Server")
}
