package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type WhatsAppService struct {
	client *whatsmeow.Client
}

func NewWhatsAppService() (*WhatsAppService, error) {
	storeDir := filepath.Join("data", "sessions")
	os.MkdirAll(storeDir, os.ModePerm)

	container, err := sqlstore.New("sqlite3", "file:"+filepath.Join(storeDir, "store.db")+"?_foreign_keys=on", waLog.Stdout("DB", "ERROR", true))
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %v", err)
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %v", err)
	}

	client := whatsmeow.NewClient(deviceStore, nil)
	svc := &WhatsAppService{client: client}

	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
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
		return "", fmt.Errorf("invalid phone: %v", err)
	}

	msg := &waProto.Message{Conversation: &message}
	resp, err := svc.client.SendMessage(context.Background(), recipient, msg)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (svc *WhatsAppService) Disconnect() {
	svc.client.Disconnect()
}
