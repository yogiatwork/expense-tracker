package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/yogiatwork/expense-tracker/config"
	"github.com/yogiatwork/expense-tracker/model"
)

func GetBotInfo() (int, error) {
	// make get api call to telegram api to get the bot info
	var botId int
	url := config.AppConfig.Telegram.BotUrl + config.AppConfig.Telegram.BotToken + "/getMe"
	res, err := http.Get(url)
	if err != nil {
		slog.Error("failed to get bot info", slog.String("error", err.Error()))
		return botId, err
	}
	defer res.Body.Close()

	var getMeResponse model.GetMeResponse
	if err := json.NewDecoder(res.Body).Decode(&getMeResponse); err != nil {
		slog.Error("failed to decode response", slog.String("error", err.Error()))
		return botId, err
	}
	slog.Info("bot info", slog.Any("bot_info", getMeResponse))

	if !getMeResponse.Ok {
		slog.Error("failed to get bot info", slog.String("error", "response not ok"))
		return botId, err
	}

	botId = getMeResponse.Result.Id
	slog.Info("bot id", slog.Int("bot_id", botId))
	return botId, nil
}

func ProcessMsg(msg model.Message) {
	msgTxt := msg.Text
	if len(msg.Entities) > 0 {
		for _, entity := range msg.Entities {
			if entity.Type == "bot_command" {
				cmd, txt := getMsgText(msgTxt)
				switch cmd {
				case "/help":
					userHelp(msg.Chat.Id)
				case "/spend":
					spendEntry(msg.Chat.Id, msg.From.FirstName, txt)
				case "/earned":
					earnedEntry(msg.Chat.Id, msg.From.FirstName, txt)
				// case "/balance":
				// 	getBalance(msg.Chat.Id, msg.From.FirstName, txt)
				// case "/spendtotal":
				// 	getSependTotal(msg.Chat.Id, msg.From.FirstName, txt)
				// case "/report":
				// 	getMonthlyReport(msg.Chat.Id, msg.From.FirstName, txt)
				default:
					unregisteredCmd(msg.Chat.Id)
				}
			}
		}
	} else {
		chooseCmd(msg.Chat.Id)
	}
}

func spendEntry(id int64, userFrom, m string) {
	slog.Info("received a message for expense " + m)
	ExpEntry, err := parseExpense(m)
	if err != nil {
		slog.Error("failed to parse expense entry " + err.Error())
		if err := sendMessage(id, "failed to parse expense entry, error: "+err.Error()); err != nil {
			slog.Error("failed to update google sheet")
		}
		return
	}
	ExpEntry.EntryBy = userFrom

	if err := SheetService.AppendSheet(ExpEntry.ToTransaction()); err != nil {
		slog.Error("failed to update google sheet " + err.Error())
		if err := sendMessage(id, "failed to update google sheet, error: "+err.Error()); err != nil {
			slog.Error("failed to update google sheet")
		}
		return
	}

	if err := sendMessage(id, "expense entry made in google sheet"); err != nil {
		slog.Error("failed to update google sheet")
	}
}

func parseExpense(m string) (model.ExpenseEntry, error) {
	var expense model.ExpenseEntry
	parts := strings.Split(m, " ")
	if len(parts) < 3 {
		return expense, errors.New("invalid expense format")
	}

	// parse amount
	amount, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return expense, errors.New("invalid amount format")
	}
	expense.Amount = amount

	// parse merchant name like "amazon" or "flipkart" or "walmart"
	if parts[1] == "" {
		return expense, errors.New("merchant name is required")
	}
	expense.Merchant = strings.Title(strings.TrimSpace(parts[1]))

	// parse spend by name like "Sapna" or "Shamim"
	if strings.ToLower(parts[2]) != "sapna" && strings.ToLower(parts[2]) != "shamim" {
		return expense, errors.New("spend by must be either Sapna or Shamim")
	}
	expense.SpendBy = strings.Title(parts[2])

	msgLen := len(parts)
	switch {
	case msgLen == 3:
		expense.Comments = ""
	case msgLen > 3:
		expense.Comments = strings.Join(parts[3:], " ")
	}

	return expense, nil
}

func earnedEntry(id int64, userFrom, m string) {
	slog.Info("received a message for earned " + m)
	EarnedEntry, err := parseEarned(m)
	if err != nil {
		slog.Error("failed to parse earned entry " + err.Error())
		if err := sendMessage(id, "failed to parse earned entry, error: "+err.Error()); err != nil {
			slog.Error("failed to update google sheet")
		}
		return
	}
	EarnedEntry.EntryBy = userFrom

	if err := SheetService.AppendSheet(EarnedEntry.ToTransaction()); err != nil {
		slog.Error("failed to update google sheet " + err.Error())
		if err := sendMessage(id, "failed to update google sheet, error: "+err.Error()); err != nil {
			slog.Error("failed to update google sheet")
		}
		return
	}

	if err := sendMessage(id, "earned entry made in google sheet"); err != nil {
		slog.Error("failed to update google sheet")
	}
}

func parseEarned(m string) (model.EarnedEntry, error) {
	var earned model.EarnedEntry
	parts := strings.Split(m, " ")
	if len(parts) < 2 {
		return earned, errors.New("invalid earned format")
	}

	// parse amount
	amount, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return earned, errors.New("invalid amount format")
	}
	earned.Amount = amount

	// parse earning type like "Online" or "Cash" or "Other"
	if strings.ToLower(parts[1]) != "online" && strings.ToLower(parts[1]) != "cash" && strings.ToLower(parts[1]) != "other" {
		return earned, errors.New("earning type must be either Online, Cash or Other")
	}
	earned.EarningType = strings.Title(parts[1])

	// parse comments if any
	msgLen := len(parts)
	switch {
	case msgLen == 2:
		earned.Comments = ""
	case msgLen > 2:
		earned.Comments = strings.Join(parts[2:], " ")
	}

	return earned, nil
}

func unregisteredCmd(id int64) {
	if err := sendMessage(id, "unknow command, try help for more details"); err != nil {
		slog.Error("failed to update google sheet")
	}
}

func chooseCmd(id int64) {
	if err := sendMessage(id, "choose from a command to interact"); err != nil {
		slog.Error("failed to update google sheet")
	}
}

func userHelp(id int64) {
	helpMsg := `Please select from one of the registered command
	/expense - for expense entry
	/help - for availble help`

	if err := sendMessage(id, helpMsg); err != nil {
		slog.Error("failed to send help message")
	}
}

// response back to the telegram user with help message
func sendMessage(id int64, txt string) error {

	url := config.AppConfig.Telegram.BotUrl + config.AppConfig.Telegram.BotToken + "/sendMessage"
	botRes := model.BotRespone{
		ChatId: id,
		Text:   txt,
	}

	// serilize
	b, err := json.Marshal(botRes)
	if err != nil {
		slog.Error("failed to serilize with error " + err.Error())
		return err
	}

	payload := bytes.NewBuffer(b)
	slog.Info("send message payload : " + payload.String())
	slog.Info("send message url : " + url)
	// send a POST request
	resp, err := http.Post(url, "application/json", payload)
	if err != nil {
		slog.Error("failed to post bot response with error " + err.Error())
		return err
	}
	resp.Body.Close()

	slog.Info("send message status : " + resp.Status)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("failed to read body " + err.Error())
	}

	//Convert the body to type string
	sb := string(body)
	slog.Info("send message reponse :" + sb)

	return nil
}

func getMsgText(m string) (string, string) {
	cmd, message, ok := strings.Cut(m, " ")
	if !ok && len(m) > 0 {
		return "", m
	}

	return strings.TrimSpace(strings.ToLower(cmd)), strings.TrimSpace(message)
}
