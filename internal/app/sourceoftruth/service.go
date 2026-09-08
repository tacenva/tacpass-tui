package sourceoftruth

import (
	"errors"
	"os"

	"github.com/tacenva/database"
	"github.com/tacenva/tacpass-core/auth"
	"github.com/tacenva/tacpass-core/util/keyring"
	"github.com/tacenva/tacpass-tui/internal/api"
	"github.com/tacenva/tacpass-tui/internal/app"
	"github.com/tacenva/tacpass-tui/internal/entity"
)

var (
	ErrForbidden = errors.New("Forbidden")
)

type Service struct {
	appDeps     *app.Deps
	sotFile     *database.DatabaseFile
	authService *auth.Service
}

func NewService(appDeps *app.Deps, authService *auth.Service) *Service {
	return &Service{
		appDeps:     appDeps,
		authService: authService,
	}
}

func (s *Service) Access(masterPassword string) error {
	sotFile, err := s.appDeps.AppDB.File("source-of-truth", masterPassword)
	if err != nil {
		return err
	}

	s.sotFile = sotFile
	if s.sotFile.Count() == 0 {
		hostname, err := os.Hostname()
		if err != nil {
			return err
		}

		token, keypair, err := s.authService.Initialize(hostname)

		if err != nil {
			return err
		}

		SoTULID, err := s.sotFile.Insert(&entity.SourceOfTruth{
			Hostname:  hostname,
			Address:   "localhost",
			AuthToken: token,
			KeyPair:   *keypair,
		})
		if err != nil {
			return err
		}

		s.appDeps.Config.SetSoTULID(SoTULID)
	}

	return nil
}

func (s *Service) Create(
	address string,
	keypair keyring.KeyPair,
) (string, error) {
	if s.sotFile == nil {
		return "", ErrForbidden
	}

	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	var fingerprint string

	s.appDeps.Client.ConfigureTLS(
		address,
		api.TLSConfig{
			Fingerprint: "",

			OnFirstTrust: func(
				newFingerprint string,
			) error {
				fingerprint = newFingerprint
				return nil
			},
		},
	)

	response, err := s.appDeps.Client.Enroll(
		address,
		api.EnrollRequest{
			Hostname:  hostname,
			PublicKey: keypair.PublicKey,
		},
	)
	if err != nil {
		return "", err
	}

	return s.sotFile.Insert(&entity.SourceOfTruth{
		Hostname:       hostname,
		Address:        address,
		AuthToken:      response.AuthToken,
		KeyPair:        keypair,
		TLSFingerprint: fingerprint,
	})
}

func (s *Service) Update(updatedSoT *entity.SourceOfTruth) error {
	if s.sotFile == nil {
		return ErrForbidden
	}
	return s.sotFile.Update(updatedSoT)
}

func (s *Service) Get(id string) (*entity.SourceOfTruth, error) {
	if s.sotFile == nil {
		return nil, ErrForbidden
	}

	var sotData entity.SourceOfTruth
	err := s.sotFile.Find(id, sotData)
	if err != nil {
		return nil, err
	}
	return &sotData, nil
}

func (s *Service) List() ([]entity.SourceOfTruth, error) {
	var sotData []entity.SourceOfTruth
	err := s.sotFile.FindAll(&sotData)
	if err != nil {
		return nil, err
	}
	return sotData, nil
}
