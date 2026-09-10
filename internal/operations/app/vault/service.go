package vault

import (
	"errors"

	"github.com/tacenva/database"
	"github.com/tacenva/tacpass-core/auth"
	"github.com/tacenva/tacpass-core/config"
	coreEntity "github.com/tacenva/tacpass-core/entity"
	vaultCore "github.com/tacenva/tacpass-core/vault"
	"github.com/tacenva/tacpass-tui/internal/operations/app"
	"github.com/tacenva/tacpass-tui/internal/operations/app/sourceoftruth"
	"github.com/tacenva/tacpass-tui/internal/operations/entity"
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
	sotService   *sourceoftruth.Service
}

func NewService(
	appDeps *app.Deps,
	context *app.Context,
	masterKey string,
	vaultService *vaultCore.Service,
	authService *auth.Service,
	sotService *sourceoftruth.Service,
) *Service {
	return &Service{
		appDeps:      appDeps,
		context:      context,
		masterKey:    masterKey,
		authService:  authService,
		vaultService: vaultService,
		sotService:   sotService,
	}
}

func (s *Service) List() ([]coreEntity.VaultAccess, bool, error) {
	if s.context == nil {
		return nil, false, errors.New("context is nil")
	}

	if s.context.SelectedSoT == nil {
		return nil, false, errors.New("selected source of truth is nil")
	}

	var vaultAccessList []coreEntity.VaultAccess

	if s.context.IsRemote {
		switch s.context.SelectedSoT.SyncMode {
		case entity.SyncModeManual:
			needSync, err := s.needSync()
			return vaultAccessList, needSync, err

		case entity.SyncModeAuto:
			vaultAccessList, err := s.Sync()
			if err != nil {
				return nil, false, err
			}

			return vaultAccessList, false, nil
		}
	} else {
		authUser, err := s.authService.GetUserData(
			s.context.SelectedSoT.AuthToken,
		)
		if err != nil {
			return nil, false, err
		}

		vaultAccessList, err = s.vaultService.VaultAccessList(authUser)
		if err != nil {
			return nil, false, err
		}
	}

	return vaultAccessList, false, nil
}

func (s *Service) needSync() (bool, error) {
	var result struct {
		NeedSync bool `json:"need_sync"`
	}

	err := s.appDeps.Client.CheckVaultSync(
		s.context.SelectedSoT.Address,
		s.context.SelectedSoT.HashSync,
		&result,
	)
	if err != nil {
		return false, err
	}

	return result.NeedSync, nil
}

func (s *Service) Sync() ([]coreEntity.VaultAccess, error) {
	var result struct {
		NeedSync        bool                     `json:"need_sync"`
		ServerVaultHash string                   `json:"server_vault_hash"`
		VaultAccessList []coreEntity.VaultAccess `json:"vault_access_list,omitempty"`
	}

	err := s.appDeps.Client.VaultSync(
		s.context.SelectedSoT.Address,
		s.context.SelectedSoT.HashSync,
		&result,
	)
	if err != nil {
		return nil, err
	}

	db := database.New(
		s.appDeps.Config.Path(
			config.NodeDirName,
			s.context.SelectedSoT.ID,
			"vault",
		),
	)

	vaultFile, err := db.File(
		"vault",
		s.masterKey,
		database.FileModeOpenOrCreate,
	)
	if err != nil {
		return nil, err
	}

	if result.NeedSync {
		vaultAccessPointers := make(
			[]*coreEntity.VaultAccess,
			0,
			len(result.VaultAccessList),
		)

		for i := range result.VaultAccessList {
			vaultAccessPointers = append(
				vaultAccessPointers,
				&result.VaultAccessList[i],
			)
		}

		err = vaultFile.Sync(
			vaultAccessPointers,
		)
		if err != nil {
			return nil, err
		}

		s.context.SelectedSoT.HashSync = result.ServerVaultHash

		if err = s.sotService.Update(
			s.context.SelectedSoT,
		); err != nil {
			return nil, err
		}
	}

	var vaultAccessList []coreEntity.VaultAccess

	err = vaultFile.FindAll(
		&vaultAccessList,
	)
	if err != nil {
		return nil, err
	}

	return vaultAccessList, nil
}

func (s *Service) createVaultRemotely(
	vaultName string,
) (*coreEntity.VaultAccess, error) {
	var vaultAccess coreEntity.VaultAccess

	body := struct {
		Name string `json:"name"`
	}{
		Name: vaultName,
	}

	err := s.appDeps.Client.CreateVault(
		s.context.SelectedSoT.Address,
		body,
		&vaultAccess,
	)
	if err != nil {
		return nil, err
	}

	return &vaultAccess, nil
}

func (s *Service) CreateVault(
	vaultName string,
) (*coreEntity.VaultAccess, error) {
	var (
		vaultAccess *coreEntity.VaultAccess
		err         error
	)

	if s.context.IsRemote {
		vaultAccess, err = s.createVaultRemotely(
			vaultName,
		)
		if err != nil {
			return nil, err
		}

		db := database.New(
			s.appDeps.Config.Path(
				config.NodeDirName,
				s.context.SelectedSoT.ID,
				"vault",
			),
		)

		vaultFile, err := db.File(
			"vault",
			s.masterKey,
			database.FileModeOpenOrCreate,
		)
		if err != nil {
			return nil, err
		}

		_, err = vaultFile.Insert(
			vaultAccess,
		)
		if err != nil {
			return nil, err
		}
	} else {
		authUser, err := s.authService.GetUserData(
			s.context.SelectedSoT.AuthToken,
		)
		if err != nil {
			return nil, err
		}

		vaultAccess, err = s.vaultService.Create(
			authUser,
			vaultName,
		)
		if err != nil {
			return nil, err
		}

		db := database.New(
			s.appDeps.Config.Path(
				config.NodeDirName,
				s.context.SelectedSoT.ID,
				"vault",
			),
		)

		vaultKey, err := s.context.SelectedSoT.KeyPair.Open(vaultAccess.VaultKey)
		if err != nil {
			return nil, err
		}

		if _, err := db.File(
			vaultAccess.VaultID,
			string(vaultKey),
			database.FileModeOpenOrCreate,
		); err != nil {
			return nil, err
		}
	}

	return vaultAccess, nil
}

func (s *Service) UpdateVault(
	vault *coreEntity.Vault,
) error {
	if vault == nil {
		return errors.New(
			"vault cannot be nil",
		)
	}

	if vault.ID == "" {
		return errors.New(
			"vault id cannot be empty",
		)
	}

	if vault.Name == "" {
		return errors.New(
			"vault name cannot be empty",
		)
	}

	if s.context == nil {
		return errors.New(
			"context is nil",
		)
	}

	if s.context.SelectedSoT == nil {
		return errors.New(
			"selected source of truth is nil",
		)
	}

	if s.context.IsRemote {
		var updated coreEntity.Vault

		err := s.appDeps.Client.UpdateVault(
			s.context.SelectedSoT.Address,
			vault.ID,
			struct {
				Name string `json:"name"`
			}{
				Name: vault.Name,
			},
			&updated,
		)
		if err != nil {
			return err
		}

		*vault = updated

		db := database.New(
			s.appDeps.Config.Path(
				config.NodeDirName,
				s.context.SelectedSoT.ID,
				"vault",
			),
		)

		vaultFile, err := db.File(
			"vault",
			s.masterKey,
			database.FileModeOpenOrCreate,
		)
		if err != nil {
			return err
		}

		var vaultAccessList []coreEntity.VaultAccess

		err = vaultFile.FindWhere(
			&vaultAccessList,
			func(data map[string]any) bool {
				vaultID, ok := data["vault_id"]

				return ok && vaultID == vault.ID
			},
		)
		if err != nil {
			return err
		}

		if len(vaultAccessList) == 0 {
			return errors.New("vault access not found")
		}

		vaultAccess := vaultAccessList[0]

		vaultAccess.Vault = *vault

		return vaultFile.Update(
			&vaultAccess,
		)
	} else {
		authUser, err := s.authService.GetUserData(
			s.context.SelectedSoT.AuthToken,
		)
		if err != nil {
			return err
		}

		if _, err = s.vaultService.Update(
			authUser,
			vault.ID,
			vault.Name,
		); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) DeleteVault(
	vault *coreEntity.Vault,
) error {
	if vault == nil {
		return errors.New(
			"vault cannot be nil",
		)
	}

	if vault.ID == "" {
		return errors.New(
			"vault id cannot be empty",
		)
	}

	if s.context == nil {
		return errors.New(
			"context is nil",
		)
	}

	if s.context.SelectedSoT == nil {
		return errors.New(
			"selected source of truth is nil",
		)
	}

	if s.context.IsRemote {
		err := s.appDeps.Client.DeleteVault(
			s.context.SelectedSoT.Address,
			vault.ID,
		)
		if err != nil {
			return err
		}

		db := database.New(
			s.appDeps.Config.Path(
				config.NodeDirName,
				s.context.SelectedSoT.ID,
				"vault",
			),
		)

		vaultFile, err := db.File(
			"vault",
			s.masterKey,
			database.FileModeOpenOrCreate,
		)
		if err != nil {
			return err
		}

		var vaultAccessList []coreEntity.VaultAccess

		err = vaultFile.FindWhere(
			&vaultAccessList,
			func(data map[string]any) bool {
				vaultID, ok := data["vault_id"]

				return ok && vaultID == vault.ID
			},
		)
		if err != nil {
			return err
		}

		vaultAccess := vaultAccessList[0]
		if _, err = vaultFile.Delete(
			vaultAccess.ID,
		); err != nil {
			return err
		}

		return s.deleteVaultFile(vault.ID)
	} else {
		authUser, err := s.authService.GetUserData(
			s.context.SelectedSoT.AuthToken,
		)
		if err != nil {
			return err
		}

		if err = s.vaultService.Delete(
			authUser,
			vault.ID,
		); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) deleteVaultFile(vaultId string) error {
	db := database.New(
		s.appDeps.Config.Path(
			config.NodeDirName,
			s.context.SelectedSoT.ID,
			"vault",
		),
	)

	return db.Delete(vaultId)
}

func (s *Service) ListRecords(
	vaultAccess *coreEntity.VaultAccess,
) ([]coreEntity.VaultRecord, bool, error) {
	if vaultAccess == nil {
		return nil, false, errors.New(
			"vault access cannot be nil",
		)
	}

	if vaultAccess.VaultID == "" {
		return nil, false, errors.New(
			"vault id cannot be empty",
		)
	}

	if s.context == nil {
		return nil, false, errors.New(
			"context is nil",
		)
	}

	if s.context.SelectedSoT == nil {
		return nil, false, errors.New(
			"selected source of truth is nil",
		)
	}

	vaultKey, err := s.context.SelectedSoT.KeyPair.Open(
		vaultAccess.VaultKey,
	)
	if err != nil {
		return nil, false, err
	}

	db := database.New(
		s.appDeps.Config.Path(
			config.NodeDirName,
			s.context.SelectedSoT.ID,
			"vault",
		),
	)

	address := s.context.SelectedSoT.Address

	if address == "localhost" {
		vaultFile, err := db.File(
			vaultAccess.VaultID,
			string(vaultKey),
			database.FileModeOpenOrCreate,
		)
		if err != nil {
			return nil, false, err
		}

		var vaultRecords []coreEntity.VaultRecord

		if err := vaultFile.FindAll(&vaultRecords); err != nil {
			return nil, false, err
		}

		return vaultRecords, false, nil
	}

	vaultFile, err := db.File(
		vaultAccess.VaultID,
		string(vaultKey),
		database.FileModeOpen,
	)

	if err != nil {
		if !errors.Is(err, database.ErrFileNotFound) {
			return nil, false, err
		}

		blob, err := s.appDeps.Client.RecordBlob(
			address,
			vaultAccess.VaultID,
		)
		if err != nil {
			return nil, false, err
		}

		if err := db.Write(
			vaultAccess.VaultID,
			blob,
		); err != nil {
			return nil, false, err
		}

		vaultFile, err = db.File(
			vaultAccess.VaultID,
			string(vaultKey),
			database.FileModeOpen,
		)
		if err != nil {
			return nil, false, err
		}
	}

	var vaultRecords []coreEntity.VaultRecord

	if err := vaultFile.FindAll(&vaultRecords); err != nil {
		return nil, false, err
	}

	version, err := db.GetVersion(
		vaultAccess.VaultID,
	)
	if err != nil {
		return nil, false, err
	}

	var result struct {
		NeedSync bool   `json:"need_sync"`
		Version  uint64 `json:"version"`
	}

	err = s.appDeps.Client.CheckRecordBlob(
		address,
		vaultAccess.VaultID,
		version,
		&result,
	)
	if err != nil {
		return nil, false, err
	}

	switch s.context.SelectedSoT.SyncMode {
	case entity.SyncModeManual:
		return vaultRecords, result.NeedSync, nil

	case entity.SyncModeAuto:
		if !result.NeedSync {
			return vaultRecords, false, nil
		}

		blob, err := s.appDeps.Client.RecordBlob(
			address,
			vaultAccess.VaultID,
		)
		if err != nil {
			return nil, false, err
		}

		if err := db.Write(
			vaultAccess.VaultID,
			blob,
		); err != nil {
			return nil, false, err
		}

		vaultFile, err = db.File(
			vaultAccess.VaultID,
			string(vaultKey),
			database.FileModeOpen,
		)
		if err != nil {
			return nil, false, err
		}

		vaultRecords = nil

		if err := vaultFile.FindAll(&vaultRecords); err != nil {
			return nil, false, err
		}

		return vaultRecords, false, nil

	default:
		return vaultRecords, false, nil
	}
}

func (s *Service) AppendRecord(
	vaultAccess *coreEntity.VaultAccess,
	record *coreEntity.VaultRecord,
) (*coreEntity.VaultRecord, error) {
	if vaultAccess == nil {
		return nil, errors.New(
			"vault access cannot be nil",
		)
	}

	if record == nil {
		return nil, errors.New(
			"record cannot be nil",
		)
	}

	if vaultAccess.VaultID == "" {
		return nil, errors.New(
			"vault id cannot be empty",
		)
	}

	if s.context == nil {
		return nil, errors.New(
			"context is nil",
		)
	}

	if s.context.SelectedSoT == nil {
		return nil, errors.New(
			"selected source of truth is nil",
		)
	}

	db := database.New(
		s.appDeps.Config.Path(
			config.NodeDirName,
			s.context.SelectedSoT.ID,
			"vault",
		),
	)

	vaultKey, err := s.context.SelectedSoT.KeyPair.Open(
		vaultAccess.VaultKey,
	)

	if err != nil {
		return nil, err
	}

	vaultFile, err := db.File(
		vaultAccess.VaultID,
		string(vaultKey),
		database.FileModeOpen,
	)
	if err != nil {
		return nil, err
	}

	encrypted, err := vaultFile.SimulateEncryption(record)
	if err != nil {
		return nil, err
	}

	var recordID string
	if s.context.IsRemote {
		recordID, err = s.appDeps.Client.CreateRecordRaw(
			s.context.SelectedSoT.Address,
			vaultAccess.VaultID,
			[]byte(encrypted),
		)
		if err != nil {
			return nil, err
		}

		record.ID = recordID
	} else {
		recordID, err := vaultFile.Insert(record)
		if err != nil {
			return nil, err
		}

		record.ID = recordID
	}

	return record, nil
}

func (s *Service) UpdateRecord(
	vaultAccess *coreEntity.VaultAccess,
	record *coreEntity.VaultRecord,
) (*coreEntity.VaultRecord, error) {
	if vaultAccess == nil {
		return nil, errors.New(
			"vault access cannot be nil",
		)
	}

	if record == nil {
		return nil, errors.New(
			"record cannot be nil",
		)
	}

	if vaultAccess.VaultID == "" {
		return nil, errors.New(
			"vault id cannot be empty",
		)
	}

	if record.ID == "" {
		return nil, errors.New(
			"record id cannot be empty",
		)
	}

	if s.context == nil {
		return nil, errors.New(
			"context is nil",
		)
	}

	if s.context.SelectedSoT == nil {
		return nil, errors.New(
			"selected source of truth is nil",
		)
	}

	vaultKey, err := s.context.SelectedSoT.KeyPair.Open(
		vaultAccess.VaultKey,
	)

	if err != nil {
		return nil, err
	}

	db := database.New(
		s.appDeps.Config.Path(
			config.NodeDirName,
			s.context.SelectedSoT.ID,
			"vault",
		),
	)

	vaultFile, err := db.File(
		vaultAccess.VaultID,
		string(vaultKey),
		database.FileModeOpen,
	)
	if err != nil {
		return nil, err
	}

	encrypted, err := vaultFile.SimulateEncryption(record)
	if err != nil {
		return nil, err
	}

	if s.context.IsRemote {
		if err := s.appDeps.Client.UpdateRecordRaw(
			s.context.SelectedSoT.Address,
			vaultAccess.VaultID,
			record.ID,
			[]byte(encrypted),
		); err != nil {
			return nil, err
		}
	} else {
		if err := vaultFile.Update(record); err != nil {
			return nil, err
		}
	}

	return record, nil
}

func (s *Service) DeleteRecord(
	vaultAccess *coreEntity.VaultAccess,
	record *coreEntity.VaultRecord,
) error {
	if vaultAccess == nil {
		return errors.New(
			"vault access cannot be nil",
		)
	}

	if record == nil {
		return errors.New(
			"record cannot be nil",
		)
	}

	if vaultAccess.VaultID == "" {
		return errors.New(
			"vault id cannot be empty",
		)
	}

	if record.ID == "" {
		return errors.New(
			"record id cannot be empty",
		)
	}

	if s.context == nil {
		return errors.New(
			"context is nil",
		)
	}

	if s.context.SelectedSoT == nil {
		return errors.New(
			"selected source of truth is nil",
		)
	}

	if s.context.IsRemote {
		if err := s.appDeps.Client.DeleteRecord(
			s.context.SelectedSoT.Address,
			vaultAccess.VaultID,
			record.ID,
		); err != nil {
			return err
		}
	}

	vaultKey, err := s.context.SelectedSoT.KeyPair.Open(
		vaultAccess.VaultKey,
	)
	if err != nil {
		return err
	}

	db := database.New(
		s.appDeps.Config.Path(
			config.NodeDirName,
			s.context.SelectedSoT.ID,
			"vault",
		),
	)

	vaultFile, err := db.File(
		vaultAccess.VaultID,
		string(vaultKey),
		database.FileModeOpen,
	)
	if err != nil {
		return err
	}

	_, err = vaultFile.Delete(record.ID)
	return err
}
