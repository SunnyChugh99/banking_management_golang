package gapi

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgent="grpcgateway-user-agent"
	XForwardedHost="x-forwarded-host"
	userAgent="user-agent"
)

type Metadata struct{
	UserAgent string
	ClientIp string

}

// map[grpcgateway-accept:[*/*] grpcgateway-cache-control:[no-cache] grpcgateway-content-type:[application/json] grpcgateway-user-agent:[PostmanRuntime/7.36.3] x-forwarded-for:[::1] x-forwarded-host:[localhost:8080]]


func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtd := &Metadata{}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok{
		fmt.Println(md)
		if useragent:=md.Get(grpcGatewayUserAgent); len(useragent) >0 {
			mtd.UserAgent = string(useragent[0])
		}
		if useragent:=md.Get(userAgent); len(useragent) >0 {
			mtd.UserAgent = string(useragent[0])
		}
		if clientIps:=md.Get(XForwardedHost); len(clientIps) >0{
			mtd.ClientIp = string(clientIps[0])
		}
	}

	if p, ok := peer.FromContext(ctx); ok{
		mtd.ClientIp = p.Addr.String()

	}
	fmt.Println("final")
	fmt.Println(mtd)
	return mtd
}