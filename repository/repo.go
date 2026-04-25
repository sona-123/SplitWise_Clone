package repository

import (
	"database/sql"
	"fmt"

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
	// 1. Start transaction
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	//2. Rollback safety
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	var expID int

	//3. Insert expense
	query := "INSERT INTO expenses(group_id, paid_by, amount) VALUES($1, $2, $3) RETURNING id"
	err = tx.QueryRow(query, exp.GroupID, exp.PaidBy, exp.Amount).Scan(&expID)
	if err != nil {
		return err
	}

	//4. Insert participants
	for _, uid := range exp.UserIds {
		query1 := "INSERT INTO participants(expense_id, user_id) VALUES($1, $2)"
		_, err := tx.Exec(query1, expID, uid)
		if err != nil {
			return err
		}
	}

	//5. Commit
	err = tx.Commit()
	return err
}

func (r *Repo) SaveGroup(name string) (models.Group, error) {
	var g models.Group
	fmt.Println(name)
	query := `INSERT INTO groups(name) VALUES($1) RETURNING id, name`
	err := r.DB.QueryRow(query, name).Scan(&g.ID, &g.Name)
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

func (r *Repo) AddUserToGroup(groupID int, userID int) error {
	query := "INSERT INTO group_members (group_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING"
	_, err := r.DB.Exec(query, groupID, userID)
	if err != nil {
		return err
	}
	return err
}
