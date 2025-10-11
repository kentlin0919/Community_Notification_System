# GEMINI.md

## Project Overview

This project is a Community Notification System built in Go. It utilizes the Gin web framework for handling HTTP requests and GORM for database interactions with PostgreSQL. The system features a versioned API (v1 and v2), JWT-based authentication, and Swagger for API documentation.

**Key Technologies:**

- **Language:** Go
- **Web Framework:** Gin
- **Database:** PostgreSQL
- **ORM:** GORM
- **Authentication:** JWT
- **API Documentation:** Swagger

**Architecture:**

The project follows a typical layered architecture for a Go web service:

- `main.go`: The application entry point, responsible for initializing the server, database, and middleware.
- `routers/`: Defines the API routes and groups them by version.
- `app/controller/`: Contains the business logic for handling requests.
- `app/models/`: Defines the data structures and models.
- `app/repositories/`: Implements the database operations.
- `database/`: Manages the database connection and schema creation.
- `middlewares/`: Includes custom middleware for CORS, JWT authentication, and cookies.
- `configs/`: Handles application configuration.
- `docs/`: Contains the auto-generated Swagger documentation.

## Building and Running

To build and run the project, follow these steps:

1.  **Install Dependencies:**

    ```bash
    go mod tidy
    ```

2.  **Set Environment Variables:**
    Create a `.env` file in the root directory with the following variables:

    ```
    PORT=:9080
    DB_HOST=127.0.0.1
    DB_USER=postgres
    DB_PASSWORD=your_password
    DB_NAME=db_Community
    DB_PORT=5432
    DB_TIMEZONE=Asia/Shanghai
    JWTPASSWORD=your_jwt_secret
    ```

3.  **Run the Application:**
    ```bash
    go run main.go
    ```

The application will be accessible at `http://localhost:9080`.

## Development Conventions

- **API Versioning:** The API is versioned under the `/api/v1` and `/api/v2` paths.
- **Authentication:** Most routes are protected by JWT authentication. The token should be provided in the `Authorization` header as a Bearer token.
- **Database:** The application uses GORM to interact with a PostgreSQL database. Database tables are created automatically on startup.
- **API Documentation:** API documentation is available at `http://localhost:9080/swagger/index.html`.

## 提交與 Pull Request 規範

- Commit：簡短、命令式、標明範疇（如：`fix(chart): clamp pan at edges`）
- 語言 ： 請以繁體中文書寫
- 請將相關變更分組，避免不相關的重構
- PR 必須包含：
  - 變更摘要與原因
  - UI 變更請附截圖/GIF
  - 手動測試步驟與影響模組/路徑
  - 若適用請附上相關 issue 或任務 ID
- 分支命名：`feature/<name>`、`fix/<name>`、`chore/<name>`
- 每個 PR 需至少一位審核者
- 請保持 PR 規模小且聚焦（盡量少於 300 行）
- 請將所有 commit 經過逐行分析後放進/docs 中的 commit*summary_2025*{當月月份}.md 的文件中
- commit*summary_2025*{當月月份}.md 請以時間新到舊排序

## 使用套件的版本

- 確保環境安裝腳本有確認是否安裝必要套件

## 軟體架構

- clean architecture

## 開發文件管理

- 必須放在 /docs
- 必須依照種類建立資料夾並依照種類存放
- 除了 swagger 相關的不處理
- 自動處理相關條件

## 產生開發文件

- 放到 /docs
- 每行程式碼逐一分析
- 必須要產生時序圖
- 必須產生 class 圖
