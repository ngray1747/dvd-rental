package dvd

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/google/uuid"
)

var (
	errInvalidDVDName = errors.New("invalid dvd name")
	errInvalidDVDID   = errors.New("invalid DVD id")
)

type Service interface {
	CreateDVD(ctx context.Context, name string) error
	RentDVD(ctx context.Context, id string) error
}

type dvdService struct {
	repo Repository
}

func NewDVDService(dvdRepo Repository) Service {
	return &dvdService{repo: dvdRepo}
}

func NewService(dvdRepo Repository, logger log.Logger, counter metrics.Counter, histogram metrics.Histogram) Service {
	var dvdService Service
	{
		dvdService = NewDVDService(dvdRepo)
		dvdService = NewLoggerMiddleware(logger)(dvdService)
		dvdService = NewMetrictMiddleware(counter, histogram)(dvdService)
	}
	return dvdService
}

func (d *dvdService) CreateDVD(ctx context.Context, name string) error {
	if name == "" {
		return errInvalidDVDName
	}

	dvd, err := NewDVD(name)
	if err != nil {
		return err
	}

	return d.repo.Store(dvd)
}

func (d *dvdService) RentDVD(ctx context.Context, id string) error {
	if id == "" {
		return errInvalidDVDID
	}
	_, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	
	return d.repo.Update(id, NotAvailable)
}
