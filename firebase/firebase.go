package firebase

import (
	"context"
	"fmt"
	"log"
	"sync"


	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type FirebaseService struct {
	client   *messaging.Client
	initOnce sync.Once
}

var instance *FirebaseService
var once sync.Once

// GetFirebaseService returns a singleton instance of FirebaseService
func GetFirebaseService() *FirebaseService {
	once.Do(func() {
		instance = &FirebaseService{}
	})
	return instance
}

// Initialize initializes the Firebase Messaging client
func (f *FirebaseService) Initialize(ctx context.Context, serviceKey string) error {
	var err error
	f.initOnce.Do(func() {

		// Convert service key to byte array and create Firebase App
		credentials := []byte(serviceKey)
		opt := option.WithCredentialsJSON(credentials)
		app, e := firebase.NewApp(ctx, nil, opt)
		if e != nil {
			err = fmt.Errorf("failed to initialize Firebase app: %w", e)
			return
		}

		// Initialize Messaging client
		f.client, e = app.Messaging(ctx)
		if e != nil {
			err = fmt.Errorf("failed to initialize Firebase Messaging client: %w", e)
			return
		}

		log.Println("Firebase Messaging client initialized successfully!")
	})

	return err
}

// GetClient returns the Messaging client
func (f *FirebaseService) GetClient() (*messaging.Client, error) {
	if f.client == nil {
		return nil, fmt.Errorf("Firebase Messaging client is not initialized. Call Initialize first")
	}
	return f.client, nil
}

// SendNotification sends a notification to a specific device
func (f *FirebaseService) SendNotification(ctx context.Context, token, title, body string) error {
	client, err := f.GetClient()
	if err != nil {
		return err
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Token: token,
	}

	_, err = client.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	log.Println("Notification sent successfully!")
	return nil
}
