package user

import (
	"context"
	company "customer-service/internal/company"
	"strconv"
)

type Publisher interface {
	Publish(routingKey string, payload interface{})
}

type Service struct {
	Repo        *Repository
	CompanyRepo *company.Repository
	Publisher   Publisher
}

func NewService(repo *Repository, companyRepo *company.Repository, pub Publisher) *Service {
	return &Service{Repo: repo, CompanyRepo: companyRepo, Publisher: pub}
}

func (s *Service) CreateUser(ctx context.Context, u *User) error {
	if _, err := s.CompanyRepo.GetByID(ctx, u.CompanyID); err != nil {
		return err
	}

	if err := s.Repo.Create(ctx, u); err != nil {
		return err
	}

	if s.Publisher != nil {
		s.Publisher.Publish("UserCreated", u)
	}
	return nil
}

func (s *Service) GetUserByID(ctx context.Context, id int) (*User, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *Service) GetUsersByCompany(ctx context.Context, companyID int) ([]User, error) {
	return s.Repo.GetByCompany(ctx, companyID)
}

func (s *Service) UpdateUser(ctx context.Context, u *User) error {
	if err := s.Repo.Update(ctx, u); err != nil {
		return err
	}

	if s.Publisher != nil {
		s.Publisher.Publish("UserUpdated", u)
	}
	return nil
}

func (s *Service) DeleteUser(ctx context.Context, id int) error {
	if err := s.Repo.Delete(ctx, id); err != nil {
		return err
	}

	if s.Publisher != nil {
		s.Publisher.Publish("UserDeleted", map[string]string{"id": strconv.Itoa(id)})
	}
	return nil
}
