package sourceoftruth

import (
	"errors"
	"os"

	"github.com/tacenva/database"
	"github.com/tacenva/tacpass-core/auth"
	coreEntity "github.com/tacenva/tacpass-core/entity"
	"github.com/tacenva/tacpass-core/permission"
	"github.com/tacenva/tacpass-core/util/keyring"
	"github.com/tacenva/tacpass-tui/internal/api"
	"github.com/tacenva/tacpass-tui/internal/app"
	"github.com/tacenva/tacpass-tui/internal/entity"
)

var (
	ErrForbidden = errors.New("Forbidden")
)

type Service struct {
	appDeps           *app.Deps
	sotFile           *database.DatabaseFile
	authService       *auth.Service
	permissionService *permission.Service
}

func NewService(
	appDeps *app.Deps,
	authService *auth.Service,
	permissionService *permission.Service,
) *Service {
	return &Service{
		appDeps:           appDeps,
		authService:       authService,
		permissionService: permissionService,
	}
}

func (s *Service) Initialize(
	name string,
	hostname string,
) (string, *keyring.KeyPair, error) {
	_, keyPair, err := s.permissionService.Create(
		name,
		coreEntity.PrivilegeAdmin,
	)
	if err != nil {
		return "", nil, err
	}

	token, err := s.authService.Enroll(
		hostname,
		keyPair.PublicKey,
		coreEntity.UserStatusApproved,
	)
	if err != nil {
		return "", nil, err
	}

	return token, keyPair, nil
}

func (s *Service) Access(
	masterPassword string,
) error {
	sotFile, err := s.appDeps.AppDB.File(
		"source-of-truth",
		masterPassword,
	)
	if err != nil {
		return err
	}

	s.sotFile = sotFile

	if s.sotFile.Count() == 0 {
		hostname, err := os.Hostname()
		if err != nil {
			return err
		}

		token, keyPair, err := s.Initialize(
			string(coreEntity.PrivilegeAdmin),
			hostname,
		)
		if err != nil {
			return err
		}

		sotULID, err := s.sotFile.Insert(
			&entity.SourceOfTruth{
				Hostname:  hostname,
				Address:   "localhost",
				AuthToken: token,
				KeyPair:   *keyPair,
			},
		)
		if err != nil {
			return err
		}

		if err := s.appDeps.Config.SetSoTULID(sotULID); err != nil {
			return err
		}
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

	return s.sotFile.Insert(
		&entity.SourceOfTruth{
			Hostname:       hostname,
			Address:        address,
			AuthToken:      response.AuthToken,
			KeyPair:        keypair,
			TLSFingerprint: fingerprint,
		},
	)
}

func (s *Service) Update(
	updatedSoT *entity.SourceOfTruth,
) error {
	if s.sotFile == nil {
		return ErrForbidden
	}

	return s.sotFile.Update(updatedSoT)
}

func (s *Service) Get(
	id string,
) (*entity.SourceOfTruth, error) {
	if s.sotFile == nil {
		return nil, ErrForbidden
	}

	var sotData entity.SourceOfTruth

	err := s.sotFile.Find(
		id,
		&sotData,
	)
	if err != nil {
		return nil, err
	}

	return &sotData, nil
}

func (s *Service) List() ([]entity.SourceOfTruth, error) {
	if s.sotFile == nil {
		return nil, ErrForbidden
	}

	var sotData []entity.SourceOfTruth

	err := s.sotFile.FindAll(&sotData)
	if err != nil {
		return nil, err
	}

	return sotData, nil
}
