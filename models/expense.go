package models

//id, paidby, amount, user
type Expense struct {
	Id      int     `json:"id"`
	PaidBy  int     `json:"paid_by"`
	Amount  float64 `json:"amount"`
	UserIds []int   `json:"user_ids"`
}
