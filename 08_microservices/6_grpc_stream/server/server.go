package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"gitlab.com/vk-go/lectures-2022-2/08_microservices/6_grpc_stream/translit"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("cant listet port", err)
	}

	server := grpc.NewServer()

	clientsWriter := []func(string){}

	tr := NewTr()
	tr.SetSendCallback = func(f func(string)) {
		clientsWriter = append(clientsWriter, f)
	}
	translit.RegisterTransliterationServer(server, tr)

	fmt.Println("starting server at :8081")
	go server.Serve(lis)

	for {
		for _, f := range clientsWriter {
			time.Sleep(time.Second)
			f("123")
		}
	}
}
