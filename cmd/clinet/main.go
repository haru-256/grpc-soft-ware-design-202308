package main

import (
	"context"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	chatv1 "github.com/haru-256/grpc-soft-ware-design-202307/gen/go/proto/chat/v1"
	"github.com/haru-256/grpc-soft-ware-design-202307/gen/go/proto/chat/v1/chatv1connect"
)

func main() {
	client := chatv1connect.NewChatServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
		connect.WithGRPC(),
	)
	ctx := context.Background()
	res, err := client.Say(
		ctx,
		connect.NewRequest(&chatv1.SayRequest{Sentence: "Hello, world!"}),
	)
	if err != nil {
		slog.Error("RPC failed", "error", err)
	}
	slog.Info("RPC response received", "response", res.Msg)
}
