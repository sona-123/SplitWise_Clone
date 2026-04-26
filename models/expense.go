package models

import "time"

type ExpenseShare struct {
	UserID int     `json:"user_id"`
	Amount float64 `json:"amount"`
}

//id, paidby, amount, user
type Expense struct {
	Id           int            `json:"id"`
	PaidBy       int            `json:"paid_by"` //The user who paid
	Amount       float64        `json:"amount"`
	GroupID      int            `json:"group_id"`
	UserIds      []int          `json:"user_ids"` //The participants who are involved
	Description  string         `json:"description"`
	Category     string         `json:"category"`
	ReceiptImage string         `json:"receipt_image"`
	SplitType    string         `json:"split_type"`
	CreatedAt    time.Time      `json:"created_at"`
	Shares       []ExpenseShare `json:"shares,omitempty"`
}
