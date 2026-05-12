# 系統常數與列舉 (System Constants & Enums)

## 1. 平台定義 (Platform Types)
後端使用以下列舉值區分客戶端類型，用於版本管理與推播過濾：
- `1`: iOS (Mobile)
- `2`: Android (Mobile)
- `3`: Web (Admin Dashboard)
- `4`: macOS (Desktop)
- `5`: Windows (Desktop)

## 2. 版本控制
- **MinVersion**: 最低支援版本。低於此版本之請求將由 `VersionCheckMiddleware` 攔截。
- **ForceUpdate**: 布林值，決定是否強制跳轉至商店更新。
