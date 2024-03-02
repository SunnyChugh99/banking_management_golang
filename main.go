package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"

	db "github.com/SunnyChugh99/banking_management_golang/db/sqlc"
	"github.com/SunnyChugh99/banking_management_golang/gapi"
	"github.com/SunnyChugh99/banking_management_golang/pb"
	"github.com/SunnyChugh99/banking_management_golang/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}

	store := db.NewStore(conn)
	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot create server", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("Cannot create listen", err)
	}

	// Get the hostname and port of the running gRPC server
	host, port := extractHostnamePort(listener.Addr().String())
	log.Printf("start grpc server at %s", listener.Addr().String())

	// Get the list of registered services and their RPC methods
	serviceInfo := grpcServer.GetServiceInfo()

	// Print information about the running gRPC server
	fmt.Printf("gRPC Server is running on: %s:%s\n", host, port)
	fmt.Println("Registered Services:")
	for serviceName, methodInfo := range serviceInfo {
		fmt.Printf("- Service: %s\n", serviceName)
		fmt.Println("  RPC Methods:")
		for _, method := range methodInfo.Methods {
			fmt.Printf("    - %s\n", method.Name)
		}
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Cannot start grpc server", err)
	}
}
func extractHostnamePort(addr string) (string, string) {
	if strings.HasPrefix(addr, "[") {
		// IPv6 address format [host]:port
		hostPort := strings.Split(addr, "]:")
		host := strings.TrimPrefix(hostPort[0], "[")
		port := hostPort[1]
		return host, port
	}

	// IPv4 address format host:port
	parts := strings.Split(addr, ":")
	return parts[0], parts[1]
}

