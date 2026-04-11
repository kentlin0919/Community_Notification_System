FROM golang:1.26.1-alpine

ARG AIR_VERSION=v1.61.7
ARG DLV_VERSION=v1.24.2

# 安裝必要套件與工具
RUN apk add --no-cache git build-base libc6-compat

# 安裝 Air 與 Delve
RUN go install github.com/air-verse/air@${AIR_VERSION}
RUN go install github.com/go-delve/delve/cmd/dlv@${DLV_VERSION}

WORKDIR /app

# 預先下載相依套件以利用 Docker 快取
COPY go.mod go.sum ./
RUN go mod download

# 預設啟動 Air
CMD ["air", "-c", ".air.toml"]
