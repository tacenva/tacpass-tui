package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/tacenva/database"
	"github.com/tacenva/tacpass-core/auth"
	"github.com/tacenva/tacpass-core/entity"
	"github.com/tacenva/tacpass-core/permission"
	"github.com/tacenva/tacpass-core/user"
	vaultCore "github.com/tacenva/tacpass-core/vault"
	"github.com/tacenva/tacpass-core/vaultaccess"
	"github.com/tacenva/tacpass-tui/internal/app"
	"github.com/tacenva/tacpass-tui/internal/config"
)

var (
	ErrForbidden = errors.New("Forbidden")
)

type Service struct {
	appDeps      *app.Deps
	context      *app.Context
	masterKey    string
	vaultService *vaultCore.Service
	authService  *auth.Service
}

func NewService(
	appDeps *app.Deps,
	context *app.Context,
	masterKey string,
) *Service {
	userRepository := user.NewRepository(appDeps.SqliteDB)
	userService := user.NewService(userRepository)
	permissionRepository := permission.NewRepository(appDeps.SqliteDB)
	permissionService := permission.NewService(permissionRepository)
	authService := auth.NewService(userService, permissionService)
	vaultaccesRepository := vaultaccess.NewRepository(appDeps.SqliteDB)
	vaultaccessService := vaultaccess.NewService(vaultaccesRepository)
	vaultRepository := vaultCore.NewRepository(appDeps.SqliteDB)
	vaultService := vaultCore.NewService(vaultRepository, appDeps.AppDB, vaultaccessService)

	return &Service{
		appDeps:      appDeps,
		context:      context,
		authService:  authService,
		masterKey:    masterKey,
		vaultService: vaultService,
	}
}

func (s *Service) List() ([]entity.VaultAccess, error) {
	debugFile, err := os.OpenFile(
		"debug.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, err
	}
	defer debugFile.Close()

	fmt.Fprintf(debugFile, "=== List ===\n")
	fmt.Fprintf(debugFile, "SoT ID: %s\n", s.context.SelectedSoT.ID)
	fmt.Fprintf(debugFile, "SoT Address: %s\n", s.context.SelectedSoT.Address)
	fmt.Fprintf(debugFile, "MasterKey exists: %t\n", s.masterKey != "")

	vaultPath := s.appDeps.Config.Path(
		config.NodeDirName,
		s.context.SelectedSoT.ID,
		"vault",
	)

	db := database.New(vaultPath)

	fmt.Fprintf(debugFile, "Vault path: %s\n", vaultPath)

	vaultFile, err := db.File("vault", s.masterKey)
	if err != nil {
		fmt.Fprintf(debugFile, "db.File error: %v\n", err)
		return nil, err
	}

	var vaultAccessList []entity.VaultAccess

	err = vaultFile.FindAll(&vaultAccessList)
	if err != nil {
		fmt.Fprintf(debugFile, "FindAll error: %v\n", err)
		return nil, err
	}

	fmt.Fprintf(debugFile, "Count: %d\n", len(vaultAccessList))

	jsonData, err := json.MarshalIndent(vaultAccessList, "", "  ")
	if err != nil {
		fmt.Fprintf(debugFile, "JSON marshal error: %v\n", err)
		return nil, err
	}

	fmt.Fprintln(debugFile, string(jsonData))
	fmt.Fprintln(debugFile, "============")

	return vaultAccessList, nil
}

func (s *Service) Sync() ([]entity.VaultAccess, error) {
	var vaultAccessList []entity.VaultAccess

	err := s.appDeps.Client.ListVaults(
		s.context.SelectedSoT.Address,
		&vaultAccessList,
	)
	if err != nil {
		return nil, err
	}

	db := database.New(s.appDeps.Config.Path(
		config.NodeDirName,
		s.context.SelectedSoT.ID,
		"vault",
	))

	vaultFile, err := db.File("vault", s.masterKey)
	if err != nil {
		return nil, err
	}

	err = vaultFile.UpdateOrCreateBulk(vaultAccessList)
	if err != nil {
		return nil, err
	}

	return vaultAccessList, nil
}

func (s *Service) CreateVault(vaultName string) (*entity.VaultAccess, error) {
	if s.context.SelectedSoT.Address != "localhost" {
		return nil, nil
	}

	authUser, err := s.authService.GetUserData(
		s.context.SelectedSoT.AuthToken,
	)
	if err != nil {
		return nil, err
	}

	vaultAccess, err := s.vaultService.Create(
		authUser,
		vaultName,
		&s.context.SelectedSoT.KeyPair,
	)
	if err != nil {
		return nil, err
	}

	db := database.New(s.appDeps.Config.Path(
		config.NodeDirName,
		s.context.SelectedSoT.ID,
		"vault",
	))

	vaultFile, err := db.File("vault", s.masterKey)
	if err != nil {
		return nil, err
	}

	_, err = vaultFile.Insert(vaultAccess)
	if err != nil {
		return nil, err
	}
	// vaultAccess.ID

	return vaultAccess, nil
}

func (s *Service) DecryptedRecordList(vaultId string) ([]entity.VaultRecord, error) {
	if s.context.SelectedSoT.Address != "localhost" {
		return nil, nil
	}

	authUser, err := s.authService.GetUserData(
		s.context.SelectedSoT.AuthToken,
	)
	if err != nil {
		return nil, err
	}

	return s.vaultService.DecryptedRecordList(
		authUser,
		vaultId,
		&s.context.SelectedSoT.KeyPair,
	)
}

func (s *Service) AppendRecord(
	vaultId string,
	record *entity.VaultRecord,
) (*entity.VaultRecord, error) {
	if s.context.SelectedSoT.Address != "localhost" {
		return nil, nil
	}

	authUser, err := s.authService.GetUserData(
		s.context.SelectedSoT.AuthToken,
	)
	if err != nil {
		return nil, err
	}

	vaultRecord, err := s.vaultService.AppendRecord(
		authUser,
		vaultId,
		record,
		&s.context.SelectedSoT.KeyPair,
	)
	if err != nil {
		return nil, err
	}

	return vaultRecord, nil
}

func (s *Service) UpdateRecord(
	vaultId string,
	record *entity.VaultRecord,
) (*entity.VaultRecord, error) {
	if s.context.SelectedSoT.Address != "localhost" {
		return nil, nil
	}

	authUser, err := s.authService.GetUserData(
		s.context.SelectedSoT.AuthToken,
	)
	if err != nil {
		return nil, err
	}

	vaultRecord, err := s.vaultService.UpdateRecord(
		authUser,
		vaultId,
		record,
		&s.context.SelectedSoT.KeyPair,
	)
	if err != nil {
		return nil, err
	}

	return vaultRecord, nil
}
