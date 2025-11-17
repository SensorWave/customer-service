package company

import (
    "context"
    "customer-service/internal/event"
)

type Service struct {
    Repo      *Repository
    Publisher *event.Publisher
}

func NewService(repo *Repository, pub *event.Publisher) *Service {
    return &Service{Repo: repo, Publisher: pub}
}

func (s *Service) CreateCompany(ctx context.Context, c *Company) error {
    if err := s.Repo.Create(ctx, c); err != nil {
        return err
    }

    s.Publisher.Publish("CompanyCreated", c)
    return nil
}

func (s *Service) GetCompanyByID(ctx context.Context, id string) (*Company, error) {
    return s.Repo.GetByID(ctx, id)
}

func (s *Service) GetAllCompanies(ctx context.Context) ([]*Company, error) {
    return s.Repo.GetAll(ctx)
}

func (s *Service) UpdateCompany(ctx context.Context, c *Company) error {
    if err := s.Repo.Update(ctx, c); err != nil {
        return err
    }

    s.Publisher.Publish("CompanyUpdated", c)
    return nil
}

func (s *Service) DeleteCompany(ctx context.Context, id string) error {
    if err := s.Repo.Delete(ctx, id); err != nil {
        return err
    }

    s.Publisher.Publish("CompanyDeleted", map[string]string{"id": id})
    return nil
}
