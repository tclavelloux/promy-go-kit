// Package example is a minimal, end-to-end reference domain seeded by the
// template so a new service scaffolded from it starts with a working
// handler -> service -> repository -> MySQL pattern to copy, instead of at
// 0% coverage with nothing to model new domains on. Delete it once the
// service has its own real domains.
package example

import (
	"context"
	"fmt"
	"time"

	"github.com/ankorstore/yokai/log"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/tclavelloux/promy-template-go/internal/domain/example/model"
	"github.com/tclavelloux/promy-template-go/internal/domain/example/repository"
)

// Service is the example domain's use-case contract.
type Service interface {
	Create(ctx context.Context, req *model.CreateExampleRequest) (*model.Example, error)
	GetByID(ctx context.Context, id string) (*model.Example, error)
}

// DefaultService is the default implementation of Service.
type DefaultService struct {
	repository repository.Repository
	validator  *validator.Validate
}

// NewService creates a new DefaultService.
func NewService(repository repository.Repository, validator *validator.Validate) *DefaultService {
	return &DefaultService{
		repository: repository,
		validator:  validator,
	}
}

var _ Service = (*DefaultService)(nil)

// Create validates and persists a new example, assigning its ID and creation timestamp.
func (s *DefaultService) Create(ctx context.Context, req *model.CreateExampleRequest) (*model.Example, error) {
	logger := log.CtxLogger(ctx)

	example := &model.Example{
		Name: req.Name,
	}
	example.SetID(uuid.New().String())
	example.SetCreatedAt(time.Now())

	if err := s.validator.Struct(example); err != nil {
		return nil, fmt.Errorf("invalid example: %w", err)
	}

	if err := s.repository.Create(ctx, example); err != nil {
		return nil, fmt.Errorf("failed to create example: %w", err)
	}

	logger.Info().Str("exampleId", example.ID).Msg("example created")

	return example, nil
}

// GetByID returns an example by its identifier.
func (s *DefaultService) GetByID(ctx context.Context, id string) (*model.Example, error) {
	example, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get example: %w", err)
	}

	return example, nil
}
