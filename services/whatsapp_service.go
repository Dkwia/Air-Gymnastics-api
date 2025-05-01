package services

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"gopkg.in/yaml.v3"
)

type WhatsAppConfig struct {
	ClientName string `yaml:"client_name"`
	LogLevel   string `yaml:"log_level"`
}

type WhatsAppService struct {
	client *whatsmeow.Client
}

var whatsAppService *WhatsAppService

func NewWhatsAppService() (*WhatsAppService, error) {
	if whatsAppService != nil {
		return whatsAppService, nil
	}

	config, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	logger := waLog.Stdout("WhatsApp", config.LogLevel, true)

	storeDir := filepath.Join("data", "sessions")
	os.MkdirAll(storeDir, os.ModePerm)
	
	container, err := sqlstore.New("sqlite3", "file:"+filepath.Join(storeDir, "store.db")+"?_foreign_keys=on", logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %v", err)
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %v", err)
	}

	client := whatsmeow.NewClient(deviceStore, logger)
	client.AddEventHandler(eventHandler)

	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			return nil, fmt.Errorf("failed to connect: %v", err)
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				fmt.Println("QR code:", evt.Code)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		err = client.Connect()
		if err != nil {
			return nil, fmt.Errorf("failed to connect: %v", err)
		}
	}

	whatsAppService = &WhatsAppService{client: client}
	return whatsAppService, nil
}

func loadConfig() (*WhatsAppConfig, error) {
	configPath := filepath.Join("data", "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config WhatsAppConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Received message:", v.Message.GetConversation())
	}
}

func (s *WhatsAppService) SendMessage(phone, message string) (string, error) {
	recipient, err := types.ParseJID(phone)
	if err != nil {
		return "", fmt.Errorf("failed to parse phone number: %v", err)
	}

	msgID, err := s.client.SendMessage(context.Background(), recipient, &waProto.Message{
		Conversation: &message,
	})
	if err != nil {
		return "", fmt.Errorf("failed to send message: %v", err)
	}

	return msgID.ID, nil
}

func (s *WhatsAppService) ValidatePhone(phone string) bool {
	_, err := types.ParseJID(phone)
	return err == nil
}

func (s *WhatsAppService) Disconnect() {
	s.client.Disconnect()
}
