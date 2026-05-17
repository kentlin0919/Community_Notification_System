package firebase

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

var (
	// FcmClient is a client for sending FCM messages.
	FcmClient *messaging.Client
	// AuthClient is a client for Firebase Authentication.
	AuthClient *auth.Client
)

// InitFirebase initializes the Firebase app and the messaging client.
func InitFirebase() {
	if _, err := os.Stat("serviceAccountKey.json"); err != nil {
		log.Printf("Firebase 未初始化：找不到 serviceAccountKey.json，推播與認證功能將停用")
		return
	}

	opt := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Printf("Firebase 初始化失敗：%v，推播與認證功能將停用", err)
		return
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Printf("Firebase Messaging 初始化失敗：%v，推播功能將停用", err)
	} else {
		FcmClient = client
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Printf("Firebase Auth 初始化失敗：%v，認證功能將停用", err)
	} else {
		AuthClient = authClient
	}
}
