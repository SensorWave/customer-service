package company

import "customer-service/internal/event"

type Service struct {
    Repo      *Repository
    Publisher *event.Publisher
}

func NewService(repo *Repository, pub *event.Publisher) *Service {
    return &Service{Repo: repo, Publisher: pub}
}

func (s *Service) CreateCompany(c *Company) error {
    if err := s.Repo.Create(c); err != nil {
        return err
    }
    s.Publisher.Publish("CompanyCreated", c)
    return nil
}
