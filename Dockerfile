# ステージ1: ビルド環境
FROM golang:1.20 AS builder

# 依存関係をコピーおよびダウンロード
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピーしてビルド
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

# ステージ2: 実行環境
FROM golang:1.20

# ビルドされたバイナリをコピー
COPY --from=builder /app/app /app

# ポート8080をエクスポート
EXPOSE 8080

# アプリケーションを実行
CMD ["/app"]
