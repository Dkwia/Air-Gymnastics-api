package models

type WhatsAppInfo struct {
	Phone     string `json:"phone" binding:"required,e164"` 
	OptIn     bool   `json:"opt_in"`                       
	LastMsgID string `json:"last_msg_id,omitempty"`        
}

type User struct {
	ID          string       `json:"id"`
	Username    string       `json:"username" binding:"required"`
	WhatsApp    WhatsAppInfo `json:"whatsapp" binding:"required"`
}

var Users = []User{
	{
		ID:       "1",
		Username: "user1",
		WhatsApp: WhatsAppInfo{
			Phone: "+12345678901",
			OptIn: true,
		},
	},
	{
		ID:       "2",
		Username: "user2",
		WhatsApp: WhatsAppInfo{
			Phone: "+12345678902",
			OptIn: false,
		},
	},
}
