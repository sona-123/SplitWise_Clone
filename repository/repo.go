package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/sona-123/splitwise_clone/models"
)

type Repo struct {
	DB *sql.DB
}

func (r *Repo) SaveUser(name, hashedPassword, email, profilePic string) (models.User, error) {
	var u models.User
	query := "INSERT INTO users(name, password, email, profile_pic) VALUES($1, $2, $3, $4) RETURNING id, name, email, profile_pic"
	err := r.DB.QueryRow(query, name, hashedPassword, email, profilePic).Scan(&u.Id, &u.Name, &u.Email, &u.ProfilePic)
	fmt.Println(err)
	return u, err
}

func (r *Repo) GetUserByID(id int) (models.User, error) {
	var u models.User
	err := r.DB.QueryRow("SELECT id, name, password FROM users WHERE id = $1", id).Scan(
		&u.Id,
		&u.Name,
		&u.Password,
	)
	fmt.Println(err)
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
		fmt.Println(err)
		_ = tx.Rollback() // safe ignore in defer
	}()
	var expID int

	//3. Insert expense
	query := "INSERT INTO expenses(group_id, paid_by, amount, description, category, receipt_image, split_type) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	err = tx.QueryRow(query, exp.GroupID, exp.PaidBy, exp.Amount, exp.Description, exp.Category, exp.ReceiptImage, exp.SplitType).Scan(&expID)
	if err != nil {
		return err
	}

	//4. Insert participants
	if exp.SplitType == "manual" {
		for _, share := range exp.Shares {
			query1 := "INSERT INTO participants(expense_id, user_id, share_amount) VALUES($1, $2, $3)"
			_, err := tx.Exec(query1, expID, share.UserID, share.Amount)
			if err != nil {
				return err
			}
		}
	} else {
		if len(exp.UserIds) == 0 {
			return errors.New("no participants")
		}
		shareAmt := exp.Amount / float64(len(exp.UserIds))
		for _, uid := range exp.UserIds {
			query1 := "INSERT INTO participants(expense_id, user_id, share_amount) VALUES($1, $2, $3)"
			_, err := tx.Exec(query1, expID, uid, shareAmt)
			if err != nil {
				return err
			}
		}
	}

	//5. Commit
	err = tx.Commit()
	fmt.Println(err)
	return err
}

func (r *Repo) SaveGroup(name string, creatorID int) (models.Group, error) {
	var g models.Group
	tx, err := r.DB.Begin()
	if err != nil {
		return models.Group{}, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := `INSERT INTO groups(name, created_by) VALUES($1, $2) RETURNING id, name`
	err = tx.QueryRow(query, name, creatorID).Scan(&g.ID, &g.Name)
	if err != nil {
		return models.Group{}, err
	}
	memberQuery := "INSERT INTO group_members(group_id, user_id) VALUES ($1, $2)"
	_, err = tx.Exec(memberQuery, g.ID, creatorID)
	if err != nil {
		return models.Group{}, err
	}
	err = tx.Commit()
	if err != nil {
		return models.Group{}, err
	}
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
	query := `INSERT INTO group_members (group_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.DB.Exec(query, groupID, userID)
	if err != nil {
		return err
	}
	return err
}

// GetTotalPaidByUser calculated the sum of all the expenses paid by this user
func (r *Repo) GetTotalPaidByUser(userID int) (float64, error) {
	var total float64
	query := "SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE paid_by=$1"
	err := r.DB.QueryRow(query, userID).Scan(&total)
	return total, err
}

// GetTotalOwedByUser calculates the sum of all shares assigned to this user
func (r *Repo) GetTotalOwedByUser(userID int) (float64, error) {
	var total float64
	query := "SELECT COALESCE(SUM(share_amount), 0) FROM participants WHERE user_id = $1"
	err := r.DB.QueryRow(query, userID).Scan(&total)
	return total, err
}
