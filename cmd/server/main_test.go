package main

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"connectrpc.com/connect"
	chatv1 "github.com/haru-256/grpc-soft-ware-design-202307/gen/go/proto/chat/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestNewChatServer(t *testing.T) {
	logger := newTestLogger()
	server, err := NewChatServer(logger)
	require.NoError(t, err, "NewChatServer() should not return error")
	require.NotNil(t, server, "NewChatServer() should not return nil server")
	assert.NotNil(t, server.validator, "NewChatServer() should create server with non-nil validator")
	assert.NotNil(t, server.logger, "NewChatServer() should create server with non-nil logger")

	// Test that the global validator is properly assigned
	assert.Same(t, globalValidator, server.validator, "Server should use the global validator")
}

func TestChatServer_Say(t *testing.T) {
	logger := newTestLogger()
	server, err := NewChatServer(logger)
	require.NoError(t, err, "Failed to create server")

	tests := []struct {
		name        string
		input       string
		wantContain string
		wantError   bool
	}{
		{
			name:        "valid input",
			input:       "Hello",
			wantContain: "You said Hello",
			wantError:   false,
		},
		{
			name:        "valid input with spaces",
			input:       "Hello World",
			wantContain: "You said Hello World",
			wantError:   false,
		},
		{
			name:        "empty input - should fail validation",
			input:       "",
			wantContain: "",
			wantError:   true,
		},
		{
			name:        "long input",
			input:       "This is a longer sentence to test the service",
			wantContain: "You said This is a longer sentence to test the service",
			wantError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// リクエスト作成
			req := connect.NewRequest(&chatv1.SayRequest{
				Sentence: tt.input,
			})

			// Say メソッド実行
			resp, sayErr := server.Say(context.Background(), req)

			// エラーチェック
			if tt.wantError {
				assert.Error(t, sayErr, "Say() should return error for invalid input")
				return
			}

			require.NoError(t, sayErr, "Say() should not return error for valid input")
			require.NotNil(t, resp, "Say() should not return nil response")
			require.NotNil(t, resp.Msg, "Say() should not return response with nil message")

			// 期待される文字列が含まれているかチェック
			assert.Equal(t, tt.wantContain, resp.Msg.GetSentence(), "Say() should return expected sentence")

			// RespondedAt が設定されているかチェック
			assert.NotNil(t, resp.Msg.GetRespondedAt(), "Say() should set RespondedAt")
		})
	}
}

func TestChatServer_Say_ContextCancellation(t *testing.T) {
	logger := newTestLogger()
	server, err := NewChatServer(logger)
	require.NoError(t, err, "Failed to create server")

	// Test with canceled context
	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	req := connect.NewRequest(&chatv1.SayRequest{
		Sentence: "Hello",
	})

	// The current implementation checks context after logging, so it should return an error
	_, err = server.Say(canceledCtx, req)
	assert.Equal(t, context.Canceled, err, "Error should be context.Canceled")
}

func TestChatServer_Say_ContextTimeout(t *testing.T) {
	logger := newTestLogger()
	server, err := NewChatServer(logger)
	require.NoError(t, err, "Failed to create server")

	// Test with already timed out context
	timeoutCtx, cancel := context.WithTimeout(context.Background(), -1) // Already expired
	defer cancel()

	req := connect.NewRequest(&chatv1.SayRequest{
		Sentence: "Hello",
	})

	_, err = server.Say(timeoutCtx, req)
	if err != nil {
		assert.Equal(t, context.DeadlineExceeded, err, "Error should be context.DeadlineExceeded")
	} else {
		t.Log("Say() completed before context timeout was checked")
	}
}

func TestChatServer_Say_RequestHeaders(t *testing.T) {
	logger := newTestLogger()
	server, err := NewChatServer(logger)
	require.NoError(t, err, "Failed to create server")

	req := connect.NewRequest(&chatv1.SayRequest{
		Sentence: "Hello with headers",
	})

	// ヘッダー追加
	req.Header().Set("User-Agent", "test-client")
	req.Header().Set("X-Request-ID", "test-123")

	resp, err := server.Say(context.Background(), req)
	require.NoError(t, err, "Say() with headers should not fail")
	require.NotNil(t, resp, "Say() should not return nil response")

	expected := "You said Hello with headers"
	assert.Equal(t, expected, resp.Msg.GetSentence(), "Say() should return expected sentence")
}

func TestChatServer_Say_ResponseValidation(t *testing.T) {
	logger := newTestLogger()
	server, err := NewChatServer(logger)
	require.NoError(t, err, "Failed to create server")

	req := connect.NewRequest(&chatv1.SayRequest{
		Sentence: "Test response validation",
	})

	resp, err := server.Say(context.Background(), req)
	require.NoError(t, err, "Say() should not fail")
	require.NotNil(t, resp, "Response should not be nil")
	require.NotNil(t, resp.Msg, "Response message should not be nil")

	// Verify response structure
	assert.NotEmpty(t, resp.Msg.GetSentence(), "Response sentence should not be empty")
	assert.NotNil(t, resp.Msg.GetRespondedAt(), "Response timestamp should not be nil")
	assert.True(t, resp.Msg.GetRespondedAt().IsValid(), "Response timestamp should be valid")
}

func TestGlobalValidator(t *testing.T) {
	// Test that the global validator is properly initialized
	assert.NotNil(t, globalValidator, "Global validator should be initialized")

	// Test validation with a valid message
	validMsg := &chatv1.SayRequest{Sentence: "Valid message"}
	err := globalValidator.Validate(validMsg)
	assert.NoError(t, err, "Valid message should pass validation")

	// Test validation with an invalid message
	invalidMsg := &chatv1.SayRequest{Sentence: ""}
	err = globalValidator.Validate(invalidMsg)
	assert.Error(t, err, "Empty sentence should fail validation")
}
