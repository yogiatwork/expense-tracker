package config

import (
	"github.com/magiconair/properties"
)

var AppConfig Config

type Config struct {
	Telegram struct {
		BotUrl   string `properties:"bot.url"`
		BotToken string `properties:"bot.token"`
		BotEmail string `properties:"bot.client_email"`
		BotKey   string `properties:"bot.private_key"`
	} `properties:"telegram"`

	GoogleSheets struct {
		CredentialsFile string `properties:"credentials.file"`
		SpreadsheetId   string `properties:"spreadsheet.id"`
	} `properties:"google.sheets"`

	WhatsappApi struct {
		VerifyToken string `properties:"verify.token"`
	} `properties:"whatsapp.api"`
}

func LoadConfig() error {
	// Load the configuration from the properties file
	p := properties.MustLoadFile(".env/config.properties", properties.UTF8)
	if err := p.Decode(&AppConfig); err != nil {
		return err
	}
	return nil
}
