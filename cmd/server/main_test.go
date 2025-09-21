package main

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	chatv1 "github.com/haru-256/grpc-soft-ware-design-202307/gen/go/proto/chat/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewChatServer(t *testing.T) {
	server, err := NewChatServer()
	require.NoError(t, err, "NewChatServer() should not return error")
	require.NotNil(t, server, "NewChatServer() should not return nil server")
	assert.NotNil(t, server.validator, "NewChatServer() should create server with non-nil validator")
}

func TestChatServer_Say(t *testing.T) {
	server, err := NewChatServer()
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
	server, err := NewChatServer()
	require.NoError(t, err, "Failed to create server")

	// キャンセル可能なコンテキスト作成
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 即座にキャンセル

	req := connect.NewRequest(&chatv1.SayRequest{
		Sentence: "Hello",
	})

	// キャンセルされたコンテキストでも正常に動作することを確認
	// (現在の実装ではコンテキストキャンセレーションを特別に処理していない)
	resp, err := server.Say(ctx, req)
	require.NoError(t, err, "Say() with canceled context should not fail")
	require.NotNil(t, resp, "Say() should not return nil response with canceled context")
}

func TestChatServer_Say_RequestHeaders(t *testing.T) {
	server, err := NewChatServer()
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
