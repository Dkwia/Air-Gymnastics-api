package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type WhatsAppService struct {
	client *whatsmeow.Client
}

func NewWhatsAppService() (*WhatsAppService, error) {
	storeDir := filepath.Join("data", "sessions")
	os.MkdirAll(storeDir, os.ModePerm)

	ctx := context.Background()
	container, err := sqlstore.New(ctx, "sqlite3", "file:"+filepath.Join(storeDir, "store.db")+"?_foreign_keys=on", waLog.Stdout("DB", "ERROR", true))
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %v", err)
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %v", err)
	}

	client := whatsmeow.NewClient(deviceStore, nil)
	client.Log = waLog.Stdout("Client", "DEBUG", true)

	svc := &WhatsAppService{client: client}

	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(ctx)
		err = client.Connect()
		if err != nil {
			return nil, err
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			}
		}
	} else {
		if err := client.Connect(); err != nil {
			return nil, fmt.Errorf("failed to connect: %v", err)
		}
	}

	return svc, nil
}

func (svc *WhatsAppService) ValidatePhone(phone string) bool {
	_, err := types.ParseJID(phone)
	return err == nil
}

func (svc *WhatsAppService) SendMessage(phone, message string) (string, error) {
	recipient, err := types.ParseJID(phone)
	if err != nil {
		log.Printf("Failed to parse JID for phone %s: %v", phone, err)
		return "", fmt.Errorf("invalid phone: %v", err)
	}

	log.Printf("Parsed JID: %v", recipient)

	msg := &waProto.Message{
		Conversation: &message,
	}

	resp, err := svc.client.SendMessage(context.Background(), recipient, msg)
	if err != nil {
		log.Printf("Failed to send message to %s: %v", phone, err)
		return "", err
	}

	return resp.ID, nil
}

func (svc *WhatsAppService) Disconnect() {
	svc.client.Disconnect()
}
