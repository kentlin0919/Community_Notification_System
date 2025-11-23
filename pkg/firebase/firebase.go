package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

var (
	// FcmClient is a client for sending FCM messages.
	FcmClient *messaging.Client
)

// InitFirebase initializes the Firebase app and the messaging client.
func InitFirebase() {
	opt := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Printf("error initializing app: %v\n", err)
		return
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Printf("error getting Messaging client: %v\n", err)
		return
	}

	FcmClient = client
}
