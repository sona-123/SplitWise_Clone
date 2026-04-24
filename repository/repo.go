package repository

import (
	"database/sql"

	"github.com/sona-123/splitwise_clone/models"
)

type Repo struct {
	DB *sql.DB
}

func (r *Repo) SaveUser(name string) (models.User, error) {
	var u models.User
	query := "INSERT INTO users(name) VALUES($1) RETURNING id, name"
	err := r.DB.QueryRow(query, name).Scan(&u)
	return u, err
}

func (r *Repo) SaveExpense(exp models.Expense) error {
	var expID int
	query := "INSERT INTO expenses(paid_by, amount) VALUES($1, $2) RETURNING id"
	err := r.DB.QueryRow(query, exp.PaidBy, exp.Amount).Scan(&expID)
	if err != nil {
		return err
	}
	for _, uid := range exp.UserIds {
		query1 := "INSERT INTO participants(expense_id, user_id) VALUES($1, $2)"
		_, err := r.DB.Exec(query1, expID, uid)
		return err
	}
}
