package config

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties"
)

var AppConfig Config

type Config struct {
	Telegram struct {
		BotUrl   string `properties:"bot.url"`
		BotToken string `properties:"bot.token,default:''"`
		BotEmail string `properties:"bot.client_email"`
		BotKey   string `properties:"bot.private_key,default:''"`
		BotTest  string `properties:"bot.test,default:''"`
	} `properties:"telegram"`

	GoogleSheets struct {
		SpreadsheetId string `properties:"spreadsheet.id"`
	} `properties:"google.sheets"`

	WhatsappApi struct {
		VerifyToken string `properties:"verify.token,default:''"`
	} `properties:"whatsapp.api"`
}

func LoadConfig() error {
	// Load the configuration from the properties file
	defaultPropPath := "default.properties"
	if gin.Mode() == gin.DebugMode {
		slog.Info("Running in debug mode, loading configuration from .env/config.properties")
		defaultPropPath = ".env/config.properties"
	}

	p := properties.MustLoadFiles([]string{defaultPropPath}, properties.UTF8, true)
	if err := p.Decode(&AppConfig); err != nil {
		return err
	}

	if gin.Mode() == gin.ReleaseMode {
		slog.Info("Running in release mode, loading configuration from environment variables")
		// Load the configuration from the environment variables

		token := os.Getenv("BOT_TOKEN")
		pvtKey := os.Getenv("BOT_PRIVATE_KEY")

		if token == "" || pvtKey == "" {
			slog.Error("BOT_TOKEN and BOT_PRIVATE_KEY environment variables are required in release mode")
		}

		AppConfig.Telegram.BotTest = os.Getenv("BOT_TEST")
		slog.Info("got a value of Test Bot Env", "test_value", os.Getenv("BOT_TEST"), "config_value", AppConfig.Telegram.BotTest)

		AppConfig.Telegram.BotToken = token
		AppConfig.Telegram.BotKey = pvtKey
		AppConfig.WhatsappApi.VerifyToken = os.Getenv("WHATSAPP_VERIFY_TOKEN")
	}

	return nil
}
