package repository

import "github.com/sona-123/splitwise_clone/models"

func (r *Repo) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, name, email, profile_pic, auth_provider FROM users WHERE email=$1`

	var user models.User
	err := r.DB.QueryRow(query, email).Scan(
		&user.Id, &user.Name, &user.Email, &user.ProfilePic, &user.AuthProvider,
	)

	if err != nil {
		return nil, err
	}
	return &user, nil
}
