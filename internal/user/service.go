package user

import (
    "context"
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

func (s *Service) CreateUser(ctx context.Context, u *User) error {
    // Check if the company exists
    if _, err := s.CompanyRepo.GetByID(ctx, u.CompanyID); err != nil {
        return err
    }

    if err := s.Repo.Create(ctx, u); err != nil {
        return err
    }

    s.Publisher.Publish("UserCreated", u)
    return nil
}

func (s *Service) GetUserByID(ctx context.Context, id string) (*User, error) {
    return s.Repo.GetByID(ctx, id)
}

func (s *Service) GetUsersByCompany(ctx context.Context, companyID string) ([]User, error) {
    return s.Repo.GetByCompany(ctx, companyID)
}

func (s *Service) UpdateUser(ctx context.Context, u *User) error {
    if err := s.Repo.Update(ctx, u); err != nil {
        return err
    }

    s.Publisher.Publish("UserUpdated", u)
    return nil
}

func (s *Service) DeleteUser(ctx context.Context, id string) error {
    if err := s.Repo.Delete(ctx, id); err != nil {
        return err
    }

    s.Publisher.Publish("UserDeleted", map[string]string{"id": id})
    return nil
}
