package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yasseldg/owngrpc/certs"
	pb "github.com/yasseldg/owngrpc/proto/helloworld"
	"github.com/yasseldg/simplego/sLog"
	"golang.org/x/oauth2"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	// Intenta obtener la dirección del entorno, sino utiliza el valor predeterminado
	if envAddr, exists := os.LookupEnv("ServerAddress"); exists {
		addr = &envAddr
	}

	flag.Parse()

	for {
		if err := run(); err != nil {
			sLog.Error("run: %s", err)
			os.Exit(1)
		}

		time.Sleep(5 * time.Second)
	}
}

func run() error {
	sLog.Info("Initializing GRPC server")
	defer sLog.Info("Exit GRPC server")

	opts, err := authOpts()
	if err != nil {
		return fmt.Errorf("authOpts: %s", err)
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, opts...)
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx := context.Background()

	return loopSayHello(ctx, c)
}

func authOpts() ([]grpc.DialOption, error) {
	// Set up the credentials for the connection.
	perRPC := oauth.TokenSource{TokenSource: oauth2.StaticTokenSource(fetchToken())}

	creds, err := credentials.NewClientTLSFromFile(certs.Path("ca_cert.pem"), "io.polygon.moon")
	if err != nil {
		return nil, fmt.Errorf("failed to load credentials: %s", err)
	}

	return []grpc.DialOption{
		// In addition to the following grpc.DialOption, callers may also use
		// the grpc.CallOption grpc.PerRPCCredentials with the RPC invocation
		// itself.
		// See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
		grpc.WithPerRPCCredentials(perRPC),
		// oauth.TokenSource requires the configuration of transport
		// credentials.
		grpc.WithTransportCredentials(creds),
	}, nil
}

// fetchToken simulates a token lookup and omits the details of proper token
// acquisition. For examples of how to acquire an OAuth2 token, see:
// https://godoc.org/golang.org/x/oauth2
func fetchToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: "some-secret-token",
	}
}

var _count int

func loopSayHello(ctx context.Context, c pb.GreeterClient) error {
	names := []string{*name, "Alice", "Bob", "Carol", "Dave", "Eve", "Frank", "Grace", "Heidi", "Ivan", "Judy", "Mallory", "Oscar", "Peggy", "Sybil", "Trent", "Victor", "Walter"}

	for i := 1; i <= 1000; i++ {
		for _, name := range names {
			r, err := c.SayHello(ctx, &pb.HelloRequest{Name: fmt.Sprintf("%s %d", name, i)})
			if err != nil {
				return fmt.Errorf("could not greet: %v", err)
			}
			_count++

			sLog.Info("Greeting %d: %s", _count, r.GetMessage())
		}
	}

	return nil
}
