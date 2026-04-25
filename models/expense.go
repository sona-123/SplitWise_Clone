package models

//id, paidby, amount, user
type Expense struct {
	Id      int     `json:"id"`
	PaidBy  int     `json:"paid_by"`
	Amount  float64 `json:"amount"`
	GroupID int     `json:"group_id"`
	UserIds []int   `json:"user_ids"`
}
