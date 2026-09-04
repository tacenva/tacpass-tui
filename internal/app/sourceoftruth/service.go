package sourceoftruth

import (
	"errors"

	"github.com/tacenva/database"
	"github.com/tacenva/tacpass-core/auth"
	"github.com/tacenva/tacpass-tui/internal/entity"
)

var (
	ErrForbidden = errors.New("Forbidden")
)

type Service struct {
	db          *database.DB
	sotFile     *database.DatabaseFile
	authService *auth.Service
}

func NewService(db *database.DB, authService *auth.Service) *Service {
	return &Service{
		db:          db,
		authService: authService,
	}
}

func (s *Service) Access(masterPassword string) error {
	sotFile, err := s.db.File("source-of-truth", masterPassword)
	if err != nil {
		return err
	}

	s.sotFile = sotFile
	if s.sotFile.Count() == 0 {
		hostname := "archpc"
		token, keypair, err := s.authService.Initialize(hostname)

		if err != nil {
			return err
		}

		if err := s.sotFile.Insert(&entity.SourceOfTruth{
			Hostname:  hostname,
			Address:   "localhost",
			AuthToken: token,
			KeyPair:   entity.KeyPair(*keypair),
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) Create(hostname string, address string, keypair entity.KeyPair) error {
	if s.sotFile == nil {
		return ErrForbidden
	}
	return s.sotFile.Insert(entity.SourceOfTruth{
		Hostname: hostname,
		Address:  address,
		KeyPair:  keypair,
	})
}

func (s *Service) Get(id string) (*entity.SourceOfTruth, error) {
	if s.sotFile == nil {
		return nil, ErrForbidden
	}

	var sourceOfTruthData entity.SourceOfTruth
	err := s.sotFile.Find(id, sourceOfTruthData)
	if err != nil {
		return nil, err
	}
	return &sourceOfTruthData, nil
}
