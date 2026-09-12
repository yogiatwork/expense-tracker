include .env/local.env

### Server 
build:
	@echo "Building server..."
	@go build -o server main.go

run:
	@echo "Starting server..."
	go run main.go

compile:
	@echo "Compiling server..."
	GOOS=linux GOARCH=amd64 go build -o bin/main-linux-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/main-darwin-amd64 main.go

clean:
	@echo "Cleaning up..."
	@rm -r bin/*


### Tesing Tracker Application
th: tracker/health
tracker/health:
	@echo "Checking tracker health tests..."
	@curl -X GET \
		-H "Content-Type: application/json" \
		$(APP_URL)/ping | jq



### Whats App API

test1: test/send/message
test/send/message:
	@echo "Sending message..."
	@# Add your message sending logic here
	@curl -X POST \
		-H "Content-Type: application/json" \
		-H "Authorization: Bearer $(TEST_API_TOKEN)" \
		-d '{"messaging_product": "whatsapp", "to": "$(CLIENT_NUMBER)", "type": "text", "text": "{\"body\": \"/spend 233 texco sapna\"  }"}' \
		$(BASE_URL)/$(TEST_PHONE_ID)/messages | jq

test2: test/send/message2
test/send/message2:
	@echo "Sending message..."
	curl --request POST \
  --url https://graph.facebook.com/v25.0/$(BUSINESS_PHONE_ID)/messages \
  --header 'Authorization: Bearer $(BUSINESS_API_TOKEN)' \
  --header 'Content-Type: application/json' \
  --data '{"messaging_product": "whatsapp","to": "919868326190", "type": "text", "text": {"body": "What can I help you today?" }	}' | jq


##### Telegram API
webhook/register:
	@echo "Registering Telegram webhook..."
	curl -X POST \
		-H "Content-Type: application/json" \
		-d '{"url": "$(TELEGRAM_WEBHOOK_URL)"}' \
		$(TELEGRAM_API_URL)/bot$(TELEGRAM_BOT_TOKEN)/setWebhook | jq

webhook/info:
	@echo "Getting Telegram webhook info..."
	@curl -X GET \
		$(TELEGRAM_API_URL)/bot$(TELEGRAM_BOT_TOKEN)/getWebhookInfo | jq
