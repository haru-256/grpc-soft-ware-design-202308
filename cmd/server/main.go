package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"buf.build/go/protovalidate"
	"connectrpc.com/connect"
	chatv1 "github.com/haru-256/grpc-soft-ware-design-202307/gen/go/proto/chat/v1"
	chatv1connect "github.com/haru-256/grpc-soft-ware-design-202307/gen/go/proto/chat/v1/chatv1connect"
	"github.com/haru-256/grpc-soft-ware-design-202307/internal"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var globalValidator protovalidate.Validator

func init() {
	v, err := protovalidate.New()
	if err != nil {
		slog.Error("Failed to create global validator", "error", err)
		os.Exit(1)
	}
	globalValidator = v
}

type ChatServer struct {
	validator protovalidate.Validator
	logger    *slog.Logger
}

func NewChatServer(logger *slog.Logger) (*ChatServer, error) {
	return &ChatServer{
		validator: globalValidator,
		logger:    logger,
	}, nil
}

func (c *ChatServer) Say(ctx context.Context, req *connect.Request[chatv1.SayRequest]) (*connect.Response[chatv1.SayResponse], error) {
	c.logger.Info(
		"Received request for chat RPC",
		"sentence", req.Msg.GetSentence(),
	)

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if err := c.validator.Validate(req.Msg); err != nil {
		slog.Error(" Request validation failed", "error", err)
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("validation error: %w", err))
	}

	res := &chatv1.SayResponse{
		Sentence:    fmt.Sprintf("You said %s", req.Msg.GetSentence()),
		RespondedAt: timestamppb.Now(),
	}

	if err := c.validator.Validate(res); err != nil {
		slog.Error("Response validation failed", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("validation error: %w", err))
	}

	// ここにビジネスロジックを実装
	return connect.NewResponse(res), nil
}

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	chatServer, err := NewChatServer(logger)
	if err != nil {
		slog.Error("Failed to create chat server", "error", err)
		return
	}
	logInterceptor := internal.NewLogger(logger)
	interceptors := connect.WithInterceptors(logInterceptor.NewUnaryInterceptor())

	mux := http.NewServeMux()
	// パスとハンドラを登録
	path, handler := chatv1connect.NewChatServiceHandler(chatServer, interceptors)
	mux.Handle(path, handler)

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	// サーバーを別のゴルーチンで起動
	go func() {
		logger.Info("Server is starting on port 8080")
		if serveErr := server.ListenAndServe(); serveErr != nil && serveErr != http.ErrServerClosed {
			logger.Error("Server failed to start", "error", serveErr)
			stop() // エラー時にコンテキストをキャンセル
		}
	}()

	// シグナルを待機
	<-ctx.Done()
	logger.Info("Shutdown signal received, starting graceful shutdown...")

	// グレースフルシャットダウンのためのタイムアウト付きコンテキスト
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// サーバーをグレースフルシャットダウン
	if shutdownErr := server.Shutdown(shutdownCtx); shutdownErr != nil {
		logger.Error("Server shutdown failed", "error", shutdownErr)
		return
	}

	logger.Info("Server stopped gracefully")
}
