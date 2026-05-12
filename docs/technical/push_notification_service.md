# 推播服務技術規格 (Push Notification Technical Spec)

## 1. Firebase Cloud Messaging (FCM) 整合
- **發送方式**: 後端採用 FCM HTTP v1 API。
- **認證**: 使用 Firebase Admin SDK 或 Google Service Account JSON。
- **Payload 結構**:
  ```json
  {
    "message": {
      "token": "device_token",
      "notification": {
        "title": "標題",
        "body": "內容"
      },
      "data": {
        "click_action": "FLUTTER_NOTIFICATION_CLICK",
        "feature": "parcel",
        "id": "123"
      }
    }
  }
  ```

## 2. 頻道管理 (Channels)
- **Parcel**: 重要等級 High，支援聲音提醒。
- **Announcement**: 重要等級 Default。
- **System**: 系統通知。
