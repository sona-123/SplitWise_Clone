package business

import (
	"github.com/sona-123/splitwise_clone/models"
	"github.com/sona-123/splitwise_clone/repository"
)

type Service struct {
	Repo *repository.Repo
}

func (s *Service) CreateUser(name string) (models.User, error) {
	return s.Repo.SaveUser(name)
}

func (s *Service) CreateExpense(exp models.Expense) error {
	return s.Repo.SaveExpense(exp)
}
