package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/dzoxploit/go-grpc-crud/entities"
	"github.com/dzoxploit/go-grpc-crud/protobuf/golang_protobuf_contact"
	"github.com/dzoxploit/go-grpc-crud/protobuf/server"
	"github.com/dzoxploit/go-grpc-crud/repos"

	"google.golang.org/grpc"

	"github.com/gorilla/mux"
)

func main() {
	LoadAppConfig()

	// Create Brand Repository
	var contactRepo repos.GenericRepo[entities.Contact] = repos.NewContactRepo()

	// push RPC server as goroutine
	go StartRPCServer(&contactRepo)

	// Initialize the router
	router := mux.NewRouter().StrictSlash(true)

	RegisterBrandRoutes(router, contactRepo)

	// Start the server
	log.Println(fmt.Sprintf("Starting Server on port %s", AppConfig.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", AppConfig.Port), router))
}

func RegisterBrandRoutes(router *mux.Router, contactRepo repos.GenericRepo[entities.Contact]) {
	NewGenericRouter[entities.Contact, *repos.ContactRepo]("/api/contacts", router, &contactRepo)
}

func StartRPCServer(contactRepo *repos.GenericRepo[entities.Contact]) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", AppConfig.RPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	golang_protobuf_contact.RegisterCRUDServer(s, server.NewCRUDServiceServer(contactRepo))

	log.Printf("gRPC server listening on port %v\n", AppConfig.RPCPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}