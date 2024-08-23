package __

import context "context"

type MyServer struct{}

var _ ExampleServiceServer = (*MyServer)(nil)

func (s *MyServer) SayHello(ctx context.Context, req *HelloRequest) (*HelloResponse, error) {
	return &HelloResponse{Message: "Hello, " + req.Name}, nil
}

func (s *MyServer) mustEmbedUnimplementedExampleServiceServer() {}
