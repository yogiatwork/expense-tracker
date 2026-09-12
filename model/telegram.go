package model

import "time"

type GetMeResponse struct {
	Ok     bool   `json:"ok"`
	Result Result `json:"result"`
}

type Result struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	IsBot     bool   `json:"is_bot"`
}

type MessageFrom struct {
	Id           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type MessageChat struct {
	Id       int64  `json:"id"`
	Type     string `json:"type"`
	IsFourm  bool   `json:"is_forum"`
	IsDirect bool   `json:"is_direct_messages"`
}

type MessageEntity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
}

type Message struct {
	MessageId int             `json:"message_id"`
	Date      int             `json:"date"`
	Text      string          `json:"text"`
	Entities  []MessageEntity `json:"entities"`
	From      MessageFrom     `json:"from"`
	Chat      MessageChat     `json:"chat"`
}

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type BotRespone struct {
	ChatId int64  `json:"chat_id"`
	Text   string `json:"text"`
}

// Entries for Telegram Bot
// UserExpense is a struct that holds the expense data from the user
type ExpenseEntry struct {
	Amount   float64 `json:"amount" required:"true"`
	Merchant string  `json:"merchant" required:"true"`
	SpendBy  string  `json:"spend_by" required:"true" oneof:"Sapna, Shamim"`
	EntryBy  string  `json:"entry_by"`
	Comments string  `json:"comments"`
}

func (e ExpenseEntry) ToTransaction() Transaction {
	return Transaction{
		Date:     time.Now().Format("2006-01-02"),
		Type:     "Expense",
		Amount:   e.Amount,
		Merchant: e.Merchant,
		SpendBy:  e.SpendBy,
		EntryBy:  e.EntryBy,
		Comments: e.Comments,
	}
}

// EarnedEntry is a struct that holds the earned data from the user
type EarnedEntry struct {
	Amount      float64 `json:"amount" required:"true"`
	EarningType string  `json:"earning_type" required:"true" oneof:"Online, Cash, Other"`
	EntryBy     string  `json:"entry_by"`
	Comments    string  `json:"comments"`
}

func (e EarnedEntry) ToTransaction() Transaction {
	return Transaction{
		Date:        time.Now().Format("2006-01-02"),
		Type:        "Income",
		Amount:      e.Amount,
		EarningType: e.EarningType,
		EntryBy:     e.EntryBy,
		Comments:    e.Comments,
	}
}
