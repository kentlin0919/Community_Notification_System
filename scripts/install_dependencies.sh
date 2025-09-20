#!/usr/bin/env bash
set -euo pipefail

# 安裝專案所需的開發環境與工具：
# - Go 1.23（或更新版本）
# - PostgreSQL 14
# - Air (hot reload)
# - Swag CLI (Swagger 產生器)
# - 專案 Go 依賴

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GO_VERSION="1.23.0"
GO_TARBALL="go${GO_VERSION}.linux-amd64.tar.gz"
GO_DOWNLOAD_URL="https://go.dev/dl/${GO_TARBALL}"

print_step() {
  printf "\n==> %s\n" "$1"
}

command_exists() {
  command -v "$1" >/dev/null 2>&1
}

install_go_tools() {
  if ! command_exists go; then
    echo "未偵測到 go，可於安裝結束後手動安裝或重新執行腳本。"
    return
  fi

  print_step "安裝 Go CLI 工具 (air, swag)"
  GO_BIN_DIR="$(go env GOPATH 2>/dev/null)/bin"
  go install github.com/cosmtrek/air@v1.51.0
  go install github.com/swaggo/swag/cmd/swag@v1.16.4

  if [ -d "${GO_BIN_DIR}" ] && [[ ":$PATH:" != *":${GO_BIN_DIR}:"* ]]; then
    echo "提醒：請將 ${GO_BIN_DIR} 加入 PATH，例如在 shell 設定檔加入："
    echo "export PATH=\"${GO_BIN_DIR}:\$PATH\""
  fi
}

download_go_dependencies() {
  if ! command_exists go; then
    echo "未偵測到 go，可於安裝結束後手動安裝或重新執行腳本。"
    return
  fi

  print_step "下載 Go 函式庫依賴 (Gin, GORM, JWT, Swagger)"
  (
    cd "${PROJECT_ROOT}"
    go mod download github.com/gin-gonic/gin@v1.10.0
    go mod download gorm.io/gorm@v1.30.0
    go mod download gorm.io/driver/postgres@v1.5.11
    go mod download gorm.io/driver/sqlite@v1.6.0
    go mod download github.com/golang-jwt/jwt@v3.2.2+incompatible
    go mod download github.com/swaggo/gin-swagger@v1.6.0
  )
}

install_macos() {
  if ! command_exists brew; then
    echo "尚未安裝 Homebrew，請先依據 https://brew.sh 指示安裝後再執行此腳本。"
    exit 1
  fi

  print_step "使用 Homebrew 安裝套件"
  brew update
  brew install go@1.23 postgresql@14 air || true
  brew upgrade go@1.23 postgresql@14 air || true

  if [ -d "/opt/homebrew/opt/go@1.23/bin" ] && [[ ":$PATH:" != *":/opt/homebrew/opt/go@1.23/bin:"* ]]; then
    echo "提醒：可將以下行加入 shell 設定檔以使用 go@1.23："
    echo "export PATH=\"/opt/homebrew/opt/go@1.23/bin:$PATH\""
  fi

  if [ -d "/usr/local/opt/go@1.23/bin" ] && [[ ":$PATH:" != *":/usr/local/opt/go@1.23/bin:"* ]]; then
    echo "提醒：可將以下行加入 shell 設定檔以使用 go@1.23："
    echo "export PATH=\"/usr/local/opt/go@1.23/bin:$PATH\""
  fi
}

install_linux() {
  if command_exists apt-get; then
    print_step "更新 APT 套件索引"
    sudo apt-get update

    print_step "安裝基本開發套件與 PostgreSQL"
    sudo apt-get install -y build-essential wget tar postgresql postgresql-contrib

    if ! command_exists go; then
      print_step "安裝 Go ${GO_VERSION}"
      wget -q "${GO_DOWNLOAD_URL}" -O "/tmp/${GO_TARBALL}"
      sudo rm -rf /usr/local/go
      sudo tar -C /usr/local -xzf "/tmp/${GO_TARBALL}"
      rm -f "/tmp/${GO_TARBALL}"
      if [[ ":$PATH:" != *":/usr/local/go/bin:"* ]]; then
        echo "提醒：請將 /usr/local/go/bin 加入 PATH，例如在 ~/.bashrc 加入："
        echo 'export PATH="/usr/local/go/bin:$HOME/go/bin:$PATH"'
      fi
    fi
  else
    echo "尚未支援此 Linux 套件管理器，請依 README 手動安裝必要套件。"
    exit 1
  fi
}

main() {
  OS_NAME="$(uname -s)"
  case "${OS_NAME}" in
    Darwin)
      install_macos
      ;;
    Linux)
      install_linux
      ;;
    *)
      echo "暫不支援的作業系統：${OS_NAME}，請參考 README.md 手動安裝。"
      exit 1
      ;;
  esac

  download_go_dependencies
  install_go_tools

  print_step "安裝流程完成"
  echo "請確認 PostgreSQL 已啟動，並於專案根目錄建立 .env 後執行 go run main.go。"
}

main "$@"
