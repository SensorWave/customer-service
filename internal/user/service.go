package user

import (
    "customer-service/internal/event"
    company "customer-service/internal/company"
)

type Service struct {
    Repo        *Repository
    CompanyRepo *company.Repository
    Publisher   *event.Publisher
}

func NewService(repo *Repository, companyRepo *company.Repository, pub *event.Publisher) *Service {
    return &Service{Repo: repo, CompanyRepo: companyRepo, Publisher: pub}
}

func (s *Service) CreateUser(u *User) error {
    // Vérifier que l'entreprise existe
    if _, err := s.CompanyRepo.GetByID(u.CompanyID); err != nil {
        return err
    }

    if err := s.Repo.Create(u); err != nil {
        return err
    }

    s.Publisher.Publish("UserCreated", u)
    return nil
}
