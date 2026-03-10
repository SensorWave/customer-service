package company

import (
	"context"
	"strconv"
)

type Publisher interface {
	Publish(routingKey string, payload interface{})
}

type Service struct {
	Repo      *Repository
	Publisher Publisher
}

func NewService(repo *Repository, pub Publisher) *Service {
	return &Service{Repo: repo, Publisher: pub}
}

func (s *Service) CreateCompany(ctx context.Context, c *Company) error {
	if err := s.Repo.Create(ctx, c); err != nil {
		return err
	}

	if s.Publisher != nil {
		s.Publisher.Publish("CompanyCreated", c)
	}
	return nil
}

func (s *Service) GetCompanyByID(ctx context.Context, id int) (*Company, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *Service) GetAllCompanies(ctx context.Context) ([]*Company, error) {
	return s.Repo.GetAll(ctx)
}

func (s *Service) UpdateCompany(ctx context.Context, c *Company) error {
	if err := s.Repo.Update(ctx, c); err != nil {
		return err
	}

	if s.Publisher != nil {
		s.Publisher.Publish("CompanyUpdated", c)
	}
	return nil
}

func (s *Service) DeleteCompany(ctx context.Context, id int) error {
	if err := s.Repo.Delete(ctx, id); err != nil {
		return err
	}

	if s.Publisher != nil {
		s.Publisher.Publish("CompanyDeleted", map[string]string{"id": strconv.Itoa(id)})
	}
	return nil
}
