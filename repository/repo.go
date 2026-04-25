package repository

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/sona-123/splitwise_clone/models"
)

type Repo struct {
	DB *sql.DB
}

func (r *Repo) SaveUser(name string) (models.User, error) {
	var u models.User
	query := "INSERT INTO users(name) VALUES($1) RETURNING id, name"
	err := r.DB.QueryRow(query, name).Scan(&u.Id, &u.Name)
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
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repo) SaveGroup(name string) (models.Group, error) {
	var g models.Group
	query := "INSERT INTO groups(name) VALUES($1) RETURNING id, name"
	err := r.DB.QueryRow(query).Scan(&g.ID, &g.Name)
	return g, err
}

func (r *Repo) GetExpensesByGroup(groupID int) ([]models.Expense, error) {
	query := `SELECT e.id, e.paid_by, e.amount, array_agg(p.user_id)
FROM expenses e
JOIN participants p
ON e.id=p.expense_id 
WHERE e.group_id=$1
GROUP BY e.id, e.paid_by, e.amount`

	rows, err := r.DB.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var expenses []models.Expense
	for rows.Next() {
		var e models.Expense
		var tempIDs pq.Int64Array
		if err := rows.Scan(&e.Id, &e.PaidBy, &e.Amount, &tempIDs); err != nil {
			return nil, err
		}
		e.UserIds = make([]int, len(tempIDs))
		for i, v := range tempIDs {
			e.UserIds[i] = int(v)
		}
		expenses = append(expenses, e)
	}
	return expenses, nil
}
