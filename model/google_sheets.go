package model

type Transaction struct {
	Date        string  `json:"date"`
	Type        string  `json:"type" required:"true" oneof:"Expense, Income"`
	Merchant    string  `json:"merchant"`
	SpendBy     string  `json:"spendby" oneof:"Sapna, Shamim"`
	EntryBy     string  `json:"entryby"`
	Amount      float64 `json:"amount" required:"true"`
	EarningType string  `json:"earningtype" oneof:"Online, Cash, Other"`
	Comments    string  `json:"comments"`
}

func (t Transaction) ToRow() []interface{} {
	return []interface{}{
		t.Date,
		t.Type,
		t.Merchant,
		t.SpendBy,
		t.Amount,
		t.EarningType,
		t.Comments,
		t.EntryBy,
	}
}
