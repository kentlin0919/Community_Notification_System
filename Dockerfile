# syntax=docker/dockerfile:1

FROM golang:1.23 AS base
WORKDIR /app
ENV CGO_ENABLED=1
ENV GO111MODULE=on

# 開發用階段：預先下載依賴並安裝 air
FROM base AS dev
RUN go install github.com/cosmtrek/air@v1.51.0
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD ["air", "-c", ".air.toml"]

# 建置階段：產生靜態二進位檔
FROM base AS builder
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/community-notification ./...

# 執行階段：精簡映像
FROM gcr.io/distroless/base-debian12 AS runner
WORKDIR /app
COPY --from=builder /app/bin/community-notification /app/community-notification
EXPOSE 9080
CMD ["/app/community-notification"]
