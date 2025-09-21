# gRPC学習プロジェクト - Soft Ware Design 2023年7月号

Soft Ware Design 2023年7月号のgRPC特集を基に、最新のツールとライブラリを使用してgRPCサーバーを構築する学習プロジェクトです。

<img src="https://gihyo.jp/assets/images/cover/2023/thumb/TH800_642307.jpg" width=50%>

- <https://gihyo.jp/magazine/SD/archive/2023/202307>

## 📚 学習内容

このプロジェクトでは以下のツールとライブラリを学習し、実装しました：

### 1. **Buf CLI** - Protocol Buffers管理ツール

### 2. **Connect-Go** - gRPCとHTTP/JSONの統合フレームワーク

### 3. **Protovalidate** - Protocol Buffersのバリデーション

### 4. **GitHub Actions** - CI/CDパイプライン

---

## 🛠 技術スタック

| カテゴリ | ツール/ライブラリ | バージョン | 用途 |
|---------|------------------|-----------|------|
| **Protocol Buffers** | Buf CLI | v1.x | .protoファイル管理・コード生成 |
| **gRPCフレームワーク** | Connect-Go | v1.x | HTTP/gRPC統合サーバー |
| **バリデーション** | Protovalidate | v0.x | リクエスト/レスポンス検証 |
| **CI/CD** | GitHub Actions | - | 自動テスト・Lint・破壊的変更検出 |

---

## 🏗 プロジェクト構造

```
.
├── .github/workflows/          # GitHub Actions設定
│   └── buf.yml                # Protocol Buffers CI
├── cmd/server/                 # サーバー実装
│   ├── main.go                # メインサーバーロジック
│   └── main_test.go           # テストファイル
├── proto/                      # Protocol Buffers定義
│   ├── address_book/v1/       # アドレス帳サービス
│   └── chat/v1/               # チャットサービス
├── gen/                        # 生成されたコード
│   └── go/proto/              # Go言語用生成コード
├── buf.yaml                    # Buf設定
├── buf.gen.yaml               # コード生成設定
├── Makefile                   # ビルドタスク
└── go.mod                     # Go依存関係
```

---

## 🔧 1. Buf CLI - Protocol Buffers管理ツール

### できること

Buf CLIは従来のprotocコンパイラーの上位互換として、以下のことができます：

#### 1. **品質管理**

- **.protoファイルのLint**: スタイルガイド準拠チェック、命名規則違反検出
- **破壊的変更検出**: APIの後方互換性を自動チェック
- **依存関係管理**: .protoファイル間の依存関係を解析・管理

#### 2. **効率的なコード生成**

- **複数言語対応**: Go, Java, Python, TypeScript等への一括生成
- **リモートプラグイン**: ローカルにコンパイラー不要
- **並列処理**: 高速なコード生成

#### 3. **モジュール管理**

- **リモートモジュール**: buf.build上の公開プロトを利用
- **バージョン管理**: プロトファイルのバージョニング

### 使い方

#### 基本セットアップ

1. **buf.yaml作成** - プロジェクトルートに配置

```yaml
version: v2
modules:
  - path: proto          # protoファイルのディレクトリ
lint:
  use:
    - STANDARD          # 標準的なlintルール適用
breaking:
  use:
    - STANDARD          # 標準的な破壊的変更検出
```

2. **buf.gen.yaml作成** - コード生成設定

```yaml
version: v2
managed:
  enabled: true
  override:
    - file_option: go_package_prefix
      value: github.com/your-repo/gen/go    # Goパッケージのベースパス

plugins:
  - remote: buf.build/protocolbuffers/go   # Goコード生成プラグイン
    out: gen/go
    opt: paths=source_relative             # 相対パス使用
```

#### 日常的なコマンド

```bash
# Lintチェック実行
buf lint

# コード生成
buf generate

# 破壊的変更チェック（mainブランチと比較）
buf breaking --against main

# フォーマット実行
buf format --write
```

#### 実践的な活用例

**1. CI/CDでの品質管理**

```yaml
# GitHub Actions例
- name: Buf Lint
  run: buf lint
- name: Breaking Change Detection
  run: buf breaking --against origin/main
```

**2. 複数言語への同時生成**

```yaml
plugins:
  - remote: buf.build/protocolbuffers/go
    out: gen/go
  - remote: buf.build/protocolbuffers/python
    out: gen/python
  - remote: buf.build/bufbuild/connect-web
    out: gen/typescript
```

    opt: paths=source_relative

- remote: buf.build/grpc/go
    out: gen/go
    opt: paths=source_relative

```

### 重要な学習ポイント

#### `paths=source_relative` vs `go_package_prefix`

```bash
# paths=source_relative使用時
proto/chat/v1/chat.proto → gen/go/proto/chat/v1/

# go_package_prefix使用時（paths=source_relativeなし）
proto/chat/v1/chat.proto → gen/go/chat/v1/
```

**重要**: `paths=source_relative` を使用すると、`go_package_prefix` は実質的に無視され、ファイルの物理パスがimportパスになる。

#### Breaking Change検出

```bash
# リモートブランチとの比較
buf breaking --against "https://github.com/${GITHUB_REPOSITORY}.git#branch=main,subdir=proto"
```

---

## 🌐 2. Connect-Go - HTTP/gRPC統合フレームワーク

### できること

Connect-Goは従来のgRPCの制限を取り払い、以下のことを実現します：

#### 1. **マルチプロトコル対応**

- **同一サービス**: gRPC, gRPC-Web, HTTP/JSON を同時サポート
- **自動変換**: プロトコル間の変換を透明に処理
- **統一API**: 1つのサービス実装で複数クライアントに対応

#### 2. **フロントエンド統合**

- **ブラウザ直接呼び出し**: JavaScriptから直接gRPCサービス呼び出し
- **TypeScript型生成**: Protocol Buffersから型安全なクライアント生成
- **CORS対応**: Webアプリケーションでの利用

#### 3. **開発効率向上**

- **cURL テスト**: コマンドラインでAPIテスト可能
- **RESTライク**: 既存のHTTPツールで開発・デバッグ
- **軽量**: 追加の依存関係最小限

### 使い方

#### 基本的なサーバー実装

1. **サービス構造体定義**

```go
type ChatServer struct {
    validator protovalidate.Validator
}

// Connect-Goのインターフェースを実装
func (s *ChatServer) Say(
    ctx context.Context, 
    req *connect.Request[chatv1.SayRequest],
) (*connect.Response[chatv1.SayResponse], error) {
    
    // リクエストからメッセージ取得
    message := req.Msg.GetSentence()
    
    // レスポンス作成
    response := &chatv1.SayResponse{
        Sentence: fmt.Sprintf("You said: %s", message),
        RespondedAt: timestamppb.Now(),
    }
    
    return connect.NewResponse(response), nil
}
```

2. **HTTPサーバー設定**

```go
func main() {
    // サービスインスタンス作成
    chatServer := &ChatServer{}
    
    // HTTPルーター作成
    mux := http.NewServeMux()
    
    // Connect-Goハンドラー生成
    path, handler := chatv1connect.NewChatServiceHandler(chatServer)
    mux.Handle(path, handler)
    
    // HTTP/2対応サーバー起動
    http.ListenAndServe("localhost:8080", h2c.NewHandler(mux, &http2.Server{}))
}
```

#### アーキテクチャ理解

**リクエストフロー**

```text
クライアント → HTTPサーバー → ルーター(mux) → Connect-Goハンドラー → サービス実装
```

- **HTTPサーバー**: ネットワーク接続管理
- **mux**: パスベースルーティング (`/chat.v1.ChatService/Say`)
- **handler**: プロトコル変換（HTTP/JSON ↔ gRPC）
- **サービス**: ビジネスロジック実装

#### 実践的な使用例

**1. エラーハンドリング**

```go
func (s *ChatServer) Say(ctx context.Context, req *connect.Request[chatv1.SayRequest]) (*connect.Response[chatv1.SayResponse], error) {
    if req.Msg.GetSentence() == "" {
        return nil, connect.NewError(
            connect.CodeInvalidArgument,
            errors.New("sentence cannot be empty"),
        )
    }
    // ... 処理続行
}
```

**2. ヘッダー処理**

```go
func (s *ChatServer) Say(ctx context.Context, req *connect.Request[chatv1.SayRequest]) (*connect.Response[chatv1.SayResponse], error) {
    // リクエストヘッダー取得
    userAgent := req.Header().Get("User-Agent")
    
    // レスポンスヘッダー設定
    resp := connect.NewResponse(&chatv1.SayResponse{...})
    resp.Header().Set("X-Response-Time", time.Now().String())
    
    return resp, nil
}
```

**3. 複数クライアント対応テスト**

```bash
# gRPCクライアント（grpcurl使用）
grpcurl -plaintext -d '{"sentence":"Hello"}' localhost:8080 chat.v1.ChatService/Say

# HTTP/JSONクライアント（curl使用）
curl -X POST http://localhost:8080/chat.v1.ChatService/Say \
  -H "Content-Type: application/json" \
  -d '{"sentence":"Hello"}'

# JavaScriptクライアント（ブラウザ）
const client = createPromiseClient(ChatService, createConnectTransport({
  baseUrl: "http://localhost:8080",
}));
const response = await client.say({sentence: "Hello"});
```

### インターセプター (Interceptors)

Connect-Goのインターセプターは、リクエスト/レスポンスのライフサイクルに割り込んで横断的な処理を実装するための仕組みです。

#### インターセプターとは

インターセプターは **Middleware** や **フィルター** と同じような概念で、以下の目的で使用されます：

- **ログ出力**: リクエスト開始/終了のログ
- **認証・認可**: トークン検証、ユーザー権限チェック
- **メトリクス**: レスポンス時間、エラー率の測定
- **リトライ**: 失敗時の自動再試行
- **レート制限**: API使用量の制御
- **バリデーション**: リクエスト/レスポンスの検証

#### インターセプターの実装

**1. 基本的なインターセプター構造**

```go
type Logger struct {
    logger *slog.Logger
}

func NewLogger(logger *slog.Logger) *Logger {
    return &Logger{logger: logger}
}

// Unary（単発）リクエスト用インターセプター
func (l *Logger) NewUnaryInterceptor() connect.UnaryInterceptorFunc {
    return func(next connect.UnaryFunc) connect.UnaryFunc {
        return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
            // リクエスト前処理
            startTime := time.Now()
            l.logStart(ctx, req)
            
            // 次の処理（他のインターセプター or 実際のハンドラー）を実行
            res, err := next(ctx, req)
            
            // レスポンス後処理
            l.logEnd(ctx, req, startTime, err)
            
            return res, err
        }
    }
}
```

**2. ログ出力の詳細実装**

```go
func (l *Logger) logStart(ctx context.Context, req connect.AnyRequest) {
    l.logger.LogAttrs(
        ctx,
        slog.LevelInfo,
        "start processing request",
        slog.String("procedure", req.Spec().Procedure),          // RPC名
        slog.String("user-agent", req.Header().Get("User-Agent")), // クライアント情報
        slog.String("http_method", req.HTTPMethod()),             // HTTPメソッド
        slog.String("peer", req.Peer().Addr),                    // クライアントアドレス
    )
}

func (l *Logger) logEnd(ctx context.Context, req connect.AnyRequest, startTime time.Time, err error) {
    duration := time.Since(startTime)
    logFields := []slog.Attr{
        slog.String("procedure", req.Spec().Procedure),
        slog.String("code", connect.CodeOf(err).String()),
        slog.Duration("duration", duration),
    }

    if err != nil {
        // エラー詳細をログ出力
        var connectErr *connect.Error
        if errors.As(err, &connectErr) {
            logFields = append(logFields, slog.String("code", connectErr.Code().String()))
        }
        logFields = append(logFields, slog.String("error", err.Error()))
        l.logger.LogAttrs(ctx, slog.LevelError, "finished with error", logFields...)
    } else {
        l.logger.LogAttrs(ctx, slog.LevelInfo, "finished", logFields...)
    }
}
```

**3. インターセプターの適用**

```go
func main() {
    logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
    
    // チャットサーバー作成
    chatServer, err := NewChatServer(logger)
    if err != nil {
        slog.Error("Failed to create chat server", "error", err)
        return
    }
    
    // インターセプター作成
    logInterceptor := internal.NewLogger(logger)
    interceptors := connect.WithInterceptors(logInterceptor.NewUnaryInterceptor())
    
    // ハンドラー作成時にインターセプターを適用
    path, handler := chatv1connect.NewChatServiceHandler(chatServer, interceptors)
    mux.Handle(path, handler)
}
```

#### インターセプターチェーン

複数のインターセプターを組み合わせることができます：

```go
// 複数のインターセプターを組み合わせ
interceptors := connect.WithInterceptors(
    authInterceptor.NewUnaryInterceptor(),     // 認証
    metricsInterceptor.NewUnaryInterceptor(),  // メトリクス
    logInterceptor.NewUnaryInterceptor(),      // ログ
)
```

**実行順序**:

```text
Request → Auth → Metrics → Log → Handler → Log → Metrics → Auth → Response
```

#### 実際のログ出力例

サーバー起動後のログ出力：

```
time=2025-09-22T00:02:53.328+09:00 level=INFO msg="start processing request" procedure=/chat.v1.ChatService/Say user-agent=test-client http_method=POST peer=127.0.0.1:54321
time=2025-09-22T00:02:53.335+09:00 level=INFO msg="Received request for chat RPC" sentence=Hello
time=2025-09-22T00:02:53.336+09:00 level=INFO msg="finished" procedure=/chat.v1.ChatService/Say code=OK duration=8.123ms
```

#### インターセプターの利点

1. **横断的関心事の分離**: ビジネスロジックとインフラ処理を分離
2. **再利用性**: 複数のサービスで同じインターセプターを利用可能
3. **テスタビリティ**: インターセプターとビジネスロジックを独立してテスト
4. **設定の柔軟性**: 環境に応じてインターセプターを組み合わせ可能

### 重要な学習ポイント

#### mux, server, handler の関係

```text
HTTP Request → HTTP Server → mux (Router) → handler (Connect) → Service (ビジネスロジック)
```

1. **HTTP Server**: ネットワーク層での接続管理
2. **mux**: URLパスに基づくルーティング
3. **handler**: HTTP ↔ gRPC プロトコル変換
4. **Service**: 実際のビジネスロジック

#### h2c (HTTP/2 cleartext) の必要性

- **gRPC要件**: gRPCはHTTP/2プロトコルベース
- **開発環境**: TLS証明書なしでHTTP/2を使用
- **本番環境**: 通常はTLS付きHTTP/2を使用

#### インターセプター vs 従来のミドルウェア

| 特徴 | Connect-Go インターセプター | 従来のHTTPミドルウェア |
|------|---------------------------|----------------------|
| **適用レベル** | gRPCメソッドレベル | HTTPリクエストレベル |
| **型安全性** | Protocol Buffers型対応 | 汎用的なHTTP処理 |
| **エラー処理** | gRPCエラーコード対応 | HTTPステータスコード |
| **メタデータ** | gRPCメタデータアクセス | HTTPヘッダーのみ |
| **ストリーミング** | ストリーミングRPC対応 | リクエスト/レスポンス単位 |

---

## ✅ 3. Protovalidate - Protocol Buffersバリデーションライブラリ

### できること

Protovalidateは、Protocol Buffersメッセージに対する強力なバリデーション機能を提供します：

#### 1. **スキーマレベルバリデーション**

- **.protoファイルでのルール定義**: バリデーション要件をAPIスキーマに埋め込み
- **コンパイル時チェック**: 不正なバリデーションルールを早期発見
- **型安全性**: Protocol Buffersの型システムと連携

#### 2. **多様なバリデーションルール**

- **文字列**: 長さ、正規表現、プリフィックス/サフィックス
- **数値**: 範囲、比較演算子
- **繰り返しフィールド**: 要素数、ユニーク制約
- **カスタムルール**: 独自のバリデーション関数

#### 3. **ランタイム検証**

- **自動検証**: サーバー/クライアント両方で同一ルール適用
- **詳細エラー**: フィールドパス付きエラーメッセージ
- **パフォーマンス**: 高速なバリデーション処理

### 使い方

#### Protocol Buffersでのルール定義

**基本的なバリデーション**

```proto
syntax = "proto3";

import "buf/validate/validate.proto";

message UserRequest {
  // 必須フィールド（空文字不可）
  string name = 1 [(buf.validate.field).string.min_len = 1];
  
  // メールアドレス形式
  string email = 2 [(buf.validate.field).string.pattern = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"];
  
  // 年齢（18歳以上120歳以下）
  int32 age = 3 [(buf.validate.field).int32 = {gte: 18, lte: 120}];
  
  // タイムスタンプ（必須）
  google.protobuf.Timestamp created_at = 4 [(buf.validate.field).timestamp = {}];
}
```

**応用的なバリデーション**

```proto
message ProductRequest {
  // 複数条件
  string product_name = 1 [(buf.validate.field).string = {
    min_len: 3,
    max_len: 50,
    pattern: "^[a-zA-Z0-9\\s-]+$"
  }];
  
  // 配列バリデーション
  repeated string tags = 2 [(buf.validate.field).repeated = {
    min_items: 1,
    max_items: 10,
    unique: true,
    items: {string: {min_len: 1, max_len: 20}}
  }];
  
  // 価格（正の値のみ）
  double price = 3 [(buf.validate.field).double.gt = 0];
  
  // 現在時刻より前の日時のみ許可
  google.protobuf.Timestamp manufactured_at = 4 [(buf.validate.field).timestamp.lt_now = true];
}
```

#### Go実装

**バリデーター初期化**

```go
type Server struct {
    validator protovalidate.Validator
}

func NewServer() (*Server, error) {
    // バリデーター作成
    v, err := protovalidate.New()
    if err != nil {
        return nil, fmt.Errorf("failed to create validator: %w", err)
    }
    
    return &Server{validator: v}, nil
}
```

**リクエストバリデーション**

```go
func (s *Server) CreateUser(ctx context.Context, req *connect.Request[pb.UserRequest]) (*connect.Response[pb.UserResponse], error) {
    // 入力バリデーション
    if err := s.validator.Validate(req.Msg); err != nil {
        // バリデーションエラーをクライアントに返す
        return nil, connect.NewError(
            connect.CodeInvalidArgument, 
            fmt.Errorf("validation failed: %w", err),
        )
    }
    
    // ビジネスロジック処理
    user := createUser(req.Msg)
    
    // レスポンスもバリデーション（推奨）
    response := &pb.UserResponse{User: user}
    if err := s.validator.Validate(response); err != nil {
        return nil, connect.NewError(
            connect.CodeInternal,
            fmt.Errorf("response validation failed: %w", err),
        )
    }
    
    return connect.NewResponse(response), nil
}
```

#### 実践的な活用例

**1. 条件付きバリデーション**

```proto
message OrderRequest {
  string order_type = 1; // "online" or "store"
  
  // オンライン注文の場合のみ必須
  string shipping_address = 2 [(buf.validate.field).string = {
    min_len: 1,
    // カスタム条件：order_type が "online" の場合のみ適用
  }];
  
  // 店舗注文の場合のみ必須
  string store_location = 3;
}
```

**2. カスタムエラーメッセージ**

```go
func (s *Server) validateWithCustomMessages(msg proto.Message) error {
    if err := s.validator.Validate(msg); err != nil {
        // バリデーションエラーをユーザーフレンドリーに変換
        var validationErr *protovalidate.ValidationError
        if errors.As(err, &validationErr) {
            return fmt.Errorf("入力エラー: %s フィールドが不正です", validationErr.Field)
        }
        return err
    }
    return nil
}
```

**3. パフォーマンス最適化**

```go
// バリデーターは再利用可能（スレッドセーフ）
var globalValidator protovalidate.Validator

func init() {
    v, err := protovalidate.New()
    if err != nil {
        panic(err)
    }
    globalValidator = v
}

// 複数のゴルーチンから安全に使用可能
func validateConcurrently(msg proto.Message) error {
    return globalValidator.Validate(msg)
}
```

        RespondedAt: timestamppb.Now(),
    }
    
    // レスポンスバリデーション
    if err := s.validator.Validate(res); err != nil {
        return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("response validation error: %w", err))
    }
    
    return connect.NewResponse(res), nil
}

```

### 重要な学習ポイント

#### バリデーション構文の注意点

```proto
# ❌ 間違い - timestampはメッセージ型
[(buf.validate.field).timestamp = true]

# ✅ 正しい - 空のメッセージを使用
[(buf.validate.field).timestamp = {}]

# ✅ 制約付き例
[(buf.validate.field).timestamp = {
  lt_now: true,
  gte: {seconds: 1609459200}
}]
```

---

## 🧪 4. Testify - Goテストライブラリ

### できること

Testifyは標準のGoテストを大幅に改善し、以下のことを可能にします：

#### 1. **アサーションの簡素化**

- **require vs assert**: 失敗時の動作を明確に制御
- **豊富な比較関数**: Equal, Contains, Greater, Nil など
- **自動型変換**: 異なる型間の比較を自動処理

#### 2. **テストの可読性向上**

- **明確なエラーメッセージ**: 期待値と実際値を自動表示
- **カスタムメッセージ**: コンテキスト情報の追加
- **構造化出力**: 複雑なオブジェクトの差分表示

#### 3. **モックとスタブ**

- **mock パッケージ**: インターフェースの自動モック生成
- **呼び出し検証**: メソッド呼び出し回数・引数の検証
- **戻り値制御**: 動的な戻り値設定

### 使い方

#### require vs assert の戦略的使い分け

**require**: 後続処理に影響する致命的チェック

```go
func TestUserService(t *testing.T) {
    // データベース接続は必須（失敗時は後続テスト不可）
    db, err := connectDB()
    require.NoError(t, err, "DB接続は必須")
    require.NotNil(t, db, "DBインスタンス必須")
    
    // サービス初期化も必須
    service := NewUserService(db)
    require.NotNil(t, service, "サービス初期化必須")
    
    // 以降のテストが意味を持つ
}
```

**assert**: 検証したいが、失敗しても他の検証は継続

```go
func TestUserValidation(t *testing.T) {
    user := &User{Name: "John", Age: 25, Email: "john@example.com"}
    
    // 複数の検証を並行実行（1つ失敗しても他は継続）
    assert.Equal(t, "John", user.Name, "名前チェック")
    assert.Equal(t, 25, user.Age, "年齢チェック")
    assert.Contains(t, user.Email, "@", "メール形式チェック")
    assert.True(t, user.IsValid(), "全体妥当性チェック")
}
```

#### テーブル駆動テストの実装

**基本パターン**

```go
func TestCalculator(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
        wantErr  bool
    }{
        {"正常ケース", 2, 3, 5, false},
        {"ゼロ計算", 0, 5, 5, false},
        {"負数計算", -2, 3, 1, false},
        {"オーバーフロー", math.MaxInt, 1, 0, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Add(tt.a, tt.b)
            
            if tt.wantErr {
                assert.Error(t, err, "エラーが期待される")
                return
            }
            
            require.NoError(t, err, "エラーは期待されない")
            assert.Equal(t, tt.expected, result, "計算結果が一致")
        })
    }
}
```

#### Connect-Go サービスのテスト

**実際のgRPCサービステスト**

```go
func TestChatService_Integration(t *testing.T) {
    // サーバー初期化
    server, err := NewChatServer()
    require.NoError(t, err, "サーバー作成は必須")
    
    tests := []struct {
        name          string
        input         *chatv1.SayRequest
        expectedMsg   string
        expectedCode  connect.Code
        setupHeaders  map[string]string
    }{
        {
            name:         "正常リクエスト",
            input:        &chatv1.SayRequest{Sentence: "Hello"},
            expectedMsg:  "You said Hello",
            expectedCode: connect.CodeOK,
        },
        {
            name:         "空文字列エラー",
            input:        &chatv1.SayRequest{Sentence: ""},
            expectedCode: connect.CodeInvalidArgument,
        },
        {
            name:         "ヘッダー付きリクエスト",
            input:        &chatv1.SayRequest{Sentence: "Test"},
            expectedMsg:  "You said Test",
            expectedCode: connect.CodeOK,
            setupHeaders: map[string]string{
                "User-Agent": "test-client",
                "X-Request-ID": "test-123",
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // リクエスト作成
            req := connect.NewRequest(tt.input)
            
            // ヘッダー設定
            for k, v := range tt.setupHeaders {
                req.Header().Set(k, v)
            }
            
            // サービス呼び出し
            resp, err := server.Say(context.Background(), req)
            
            // エラーケースの検証
            if tt.expectedCode != connect.CodeOK {
                assert.Error(t, err, "エラーが期待される")
                assert.Equal(t, tt.expectedCode, connect.CodeOf(err), "期待するエラーコード")
                return
            }
            
            // 正常ケースの検証
            require.NoError(t, err, "エラーは期待されない")
            require.NotNil(t, resp, "レスポンスは必須")
            require.NotNil(t, resp.Msg, "メッセージは必須")
            
            assert.Equal(t, tt.expectedMsg, resp.Msg.GetSentence(), "メッセージ内容")
            assert.NotNil(t, resp.Msg.GetRespondedAt(), "タイムスタンプ設定")
            assert.True(t, resp.Msg.GetRespondedAt().IsValid(), "有効なタイムスタンプ")
        })
    }
}
```

#### モック活用例

**外部依存をモック化**

```go
//go:generate mockery --name=UserRepository --output=mocks
type UserRepository interface {
    GetUser(id string) (*User, error)
    SaveUser(*User) error
}

func TestUserService_WithMock(t *testing.T) {
    // モック作成
    mockRepo := &mocks.UserRepository{}
    service := NewUserService(mockRepo)
    
    // モックの期待値設定
    expectedUser := &User{ID: "123", Name: "John"}
    mockRepo.On("GetUser", "123").Return(expectedUser, nil)
    
    // テスト実行
    user, err := service.GetUser("123")
    
    // 検証
    require.NoError(t, err)
    assert.Equal(t, expectedUser, user)
    
    // モック呼び出し確認
    mockRepo.AssertExpectations(t)
    mockRepo.AssertCalled(t, "GetUser", "123")
}
```

#### パフォーマンステスト

**ベンチマークテスト**

```go
func BenchmarkChatService_Say(b *testing.B) {
    server, err := NewChatServer()
    require.NoError(b, err)
    
    req := connect.NewRequest(&chatv1.SayRequest{
        Sentence: "Benchmark test message",
    })
    
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            _, err := server.Say(context.Background(), req)
            if err != nil {
                b.Fatal(err)
            }
        }
    })
}
```

    },
    {
        name:        "empty input - should fail validation",
        input:       "",
        wantError:   true,
    },
}

```

---

## 🚀 5. GitHub Actions - CI/CDパイプライン

### Protocol Buffers CI

#### `.github/workflows/buf.yml`

```yaml
name: Protobuf CI

on:
  push:
    branches: [main]
  pull_request:

jobs:
  breaking-change:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Install Buf CLI
        uses: bufbuild/buf-action@v1
        with:
          version: stable
      - name: breaking-change
        run: buf breaking --against "https://github.com/${GITHUB_REPOSITORY}.git#branch=main,subdir=proto"
```

### 重要な学習ポイント

#### Breaking Change検出の仕組み

1. **プルリクエスト時**: 現在のブランチとmainブランチを比較
2. **比較対象**: `subdir=proto` でprotoディレクトリのみ
3. **検出例**: フィールド削除、型変更、フィールド番号変更など

---

## 🔄 6. グレースフルシャットダウン

### 実装

```go
func main() {
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()
    
    server := &http.Server{
        Addr: "localhost:8080",
        Handler: h2c.NewHandler(mux, &http2.Server{}),
    }
    
    // 非同期でサーバー起動
    go func() {
        if serveErr := server.ListenAndServe(); serveErr != nil && serveErr != http.ErrServerClosed {
            slog.Error("Server failed", "error", serveErr)
            stop()
        }
    }()
    
    // シグナル待機
    <-ctx.Done()
    slog.Info("Shutting down...")
    
    // グレースフルシャットダウン
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if shutdownErr := server.Shutdown(shutdownCtx); shutdownErr != nil {
        slog.Error("Shutdown failed", "error", shutdownErr)
    }
}
```

### 動作フロー

1. **シグナル受信** (`Ctrl+C`, `SIGTERM`)
2. **新規接続停止**
3. **既存リクエスト完了待機** (最大30秒)
4. **サーバー停止**

---

## 🎯 主要な学習成果

### 1. **Protocol Buffers管理の現代的アプローチ**

- 従来のprotocではなくBuf CLIを使用
- リント・破壊的変更検出の自動化
- リモートプラグインによる柔軟なコード生成

### 2. **gRPCとHTTPの統合**

- Connect-Goによる単一サービスでの複数プロトコル対応
- ブラウザ・cURL対応のAPI設計
- HTTP/2の理解と実装

### 3. **型安全なバリデーション**

- Protocol Buffers定義でのバリデーション記述
- コンパイル時型安全性の確保
- リクエスト/レスポンス両方の検証

### 4. **テスト駆動開発**

- testifyによるクリーンなテストコード
- require/assertの適切な使い分け
- 実際のgRPCサービステスト手法

---

## 🚀 クイックスタート

### 前提条件

- Go 1.21+
- Buf CLI

### セットアップ

```bash
# リポジトリクローン
git clone https://github.com/haru-256/grpc-soft-ware-design-202307.git
cd grpc-soft-ware-design-202307

# 依存関係インストール
go mod download

# Protocol Buffers コード生成
buf generate

# サーバー起動
go run cmd/server/main.go
```

### テスト実行

```bash
# ユニットテスト
go test ./...

# Lint実行
make lint

# cURLでテスト
curl -X POST http://localhost:8080/chat.v1.ChatService/Say \
  -H "Content-Type: application/json" \
  -d '{"sentence": "Hello World"}'
```

---

## 📖 参考資料

- [Buf CLI Documentation](https://buf.build/docs/)
- [Connect-Go Documentation](https://connectrpc.com/docs/go/)
- [Protovalidate Documentation](https://github.com/bufbuild/protovalidate)
- [Software Design 2023年7月号](https://gihyo.jp/magazine/SD/archive/2023/202307)

---

## 🏷 学習のポイント

このプロジェクトを通じて、Protocol Buffersを中心とした現代的なAPI開発手法を習得しました。特に：

- **ツールチェーンの理解**: Buf CLIによる効率的な開発フロー
- **プロトコル統合**: gRPCとHTTPの統一的な扱い
- **型安全性**: コンパイル時からランタイムまでの一貫した型チェック
- **CI/CD統合**: 破壊的変更の自動検出による安全なAPI進化

これらの知識は、スケーラブルで保守性の高いAPIサーバー開発に直接活用できます。
