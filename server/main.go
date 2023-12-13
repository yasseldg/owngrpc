// Copyright (c) 2022 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/yasseldg/owngrpc/certs"
	pb "github.com/yasseldg/owngrpc/proto/helloworld"
	"github.com/yasseldg/simplego/sLog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	ctx := context.Background()

	flag.Parse()

	if err := run(ctx); err != nil {
		sLog.Error("run: %s", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	sLog.Info("Initializing GRPC server")
	defer sLog.Info("Exit GRPC server")

	opts, err := authOpts()
	if err != nil {
		return fmt.Errorf("authOpts: %s", err)
	}

	s := grpc.NewServer(opts...)

	pb.RegisterGreeterServer(s, &server{})

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		return fmt.Errorf("failed to listen: %s", err)
	}

	sLog.Info("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	sLog.Warn("Total greetings: %d", _count)
	return nil
}

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

var _count int

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	_count++
	sLog.Info("Received %4d: %v", _count, in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func authOpts() ([]grpc.ServerOption, error) {
	cert, err := tls.LoadX509KeyPair(certs.Path("server_cert.pem"), certs.Path("server_key.pem"))
	if err != nil {
		return nil, fmt.Errorf("failed to load key pair: %s", err)
	}

	return []grpc.ServerOption{
		// The following grpc.ServerOption adds an interceptor for all unary
		// RPCs. To configure an interceptor for streaming RPCs, see:
		// https://godoc.org/google.golang.org/grpc#StreamInterceptor
		grpc.UnaryInterceptor(ensureValidToken),
		// Enable TLS for all incoming connections.
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	}, nil
}

// valid validates the authorization.
func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	// Perform the token validation here. For the sake of this example, the code
	// here forgoes any of the usual OAuth2 token validation and instead checks
	// for a token matching an arbitrary string.
	return token == "some-secret-token"
}

// ensureValidToken ensures a valid token exists within a request's metadata. If
// the token is missing or invalid, the interceptor blocks execution of the
// handler and returns an error. Otherwise, the interceptor invokes the unary
// handler.
func ensureValidToken(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	// The keys within metadata.MD are normalized to lowercase.
	// See: https://godoc.org/google.golang.org/grpc/metadata#New
	if !valid(md["authorization"]) {
		return nil, errInvalidToken
	}
	// Continue execution of handler after ensuring a valid token.
	return handler(ctx, req)
}
