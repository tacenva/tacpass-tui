package vault

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

	"github.com/tacenva/database"
	"github.com/tacenva/tacpass-core/auth"
	coreEntity "github.com/tacenva/tacpass-core/entity"
	vaultCore "github.com/tacenva/tacpass-core/vault"
	"github.com/tacenva/tacpass-tui/internal/app"
	"github.com/tacenva/tacpass-tui/internal/app/sourceoftruth"
	"github.com/tacenva/tacpass-tui/internal/config"
	"github.com/tacenva/tacpass-tui/internal/entity"
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

	vaultPath := s.appDeps.Config.Path(
		config.NodeDirName,
		s.context.SelectedSoT.ID,
		"vault",
	)

	db := database.New(
		vaultPath,
	)

	vaultFile, err := db.File(
		"vault",
		s.masterKey,
	)
	if err != nil {
		return nil, false, err
	}

	var vaultAccessList []coreEntity.VaultAccess

	err = vaultFile.FindAll(
		&vaultAccessList,
	)
	if err != nil {
		return nil, false, err
	}

	if s.context.IsRemote {
		switch s.context.SelectedSoT.SyncMode {
		case entity.SyncModeManual:
			needSync, err := s.needSync()
			if err != nil {
				return nil, false, err
			}

			return vaultAccessList, needSync, nil

		case entity.SyncModeAuto:
			vaultAccessList, err := s.sync()
			if err != nil {
				return nil, false, err
			}

			return vaultAccessList, false, nil
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

func (s *Service) sync() ([]coreEntity.VaultAccess, error) {
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

// func (s *Service) Sync() ([]coreEntity.VaultAccess, error) {
// 	var vaultAccessList []coreEntity.VaultAccess

// 	err := s.appDeps.Client.ListVaults(
// 		s.context.SelectedSoT.Address,
// 		&vaultAccessList,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	db := database.New(
// 		s.appDeps.Config.Path(
// 			config.NodeDirName,
// 			s.context.SelectedSoT.ID,
// 			"vault",
// 		),
// 	)

// 	vaultFile, err := db.File(
// 		"vault",
// 		s.masterKey,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	vaultAccessPointers := make(
// 		[]*coreEntity.VaultAccess,
// 		0,
// 		len(vaultAccessList),
// 	)

// 	for i := range vaultAccessList {
// 		vaultAccessPointers = append(
// 			vaultAccessPointers,
// 			&vaultAccessList[i],
// 		)
// 	}

// 	err = vaultFile.UpdateOrCreateBulk(
// 		vaultAccessPointers,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return vaultAccessList, nil
// }

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

	if s.context.SelectedSoT.Address != "localhost" {
		vaultAccess, err = s.createVaultRemotely(
			vaultName,
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

	if s.context.SelectedSoT.Address != "localhost" {
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
	} else {
		authUser, err := s.authService.GetUserData(
			s.context.SelectedSoT.AuthToken,
		)
		if err != nil {
			return err
		}

		updated, err := s.vaultService.Update(
			authUser,
			vault.ID,
			vault.Name,
		)
		if err != nil {
			return err
		}

		*vault = *updated
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

	if s.context.SelectedSoT.Address != "localhost" {
		err := s.appDeps.Client.DeleteVault(
			s.context.SelectedSoT.Address,
			vault.ID,
		)
		if err != nil {
			return err
		}
	} else {
		authUser, err := s.authService.GetUserData(
			s.context.SelectedSoT.AuthToken,
		)
		if err != nil {
			return err
		}

		err = s.vaultService.Delete(
			authUser,
			vault.ID,
		)
		if err != nil {
			return err
		}
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
) ([]coreEntity.VaultRecord, error) {
	if vaultAccess == nil {
		return nil, errors.New(
			"vault access cannot be nil",
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

	if s.context.SelectedSoT.Address == "localhost" {
		return s.listRecordsLocally(
			vaultAccess,
		)
	}

	return s.listRecordsRemotely(
		vaultAccess,
	)
}

func (s *Service) listRecordsLocally(
	vaultAccess *coreEntity.VaultAccess,
) ([]coreEntity.VaultRecord, error) {
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
	)
	if err != nil {
		return nil, err
	}

	var records []coreEntity.VaultRecord

	err = vaultFile.FindAll(
		&records,
	)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (s *Service) listRecordsRemotely(
	vaultAccess *coreEntity.VaultAccess,
) ([]coreEntity.VaultRecord, error) {
	var encryptedRecords map[string][]byte

	err := s.appDeps.Client.ListRecords(
		s.context.SelectedSoT.Address,
		vaultAccess.VaultID,
		&encryptedRecords,
	)
	if err != nil {
		return nil, err
	}

	vaultKeyEncoded, err := s.context.SelectedSoT.KeyPair.Open(
		vaultAccess.VaultKey,
	)
	if err != nil {
		return nil, err
	}

	vaultKey, err := base64.RawURLEncoding.DecodeString(
		string(vaultKeyEncoded),
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

	records := make(
		[]coreEntity.VaultRecord,
		0,
		len(encryptedRecords),
	)

	for recordID, encryptedData := range encryptedRecords {
		decrypted, err := db.Decrypt(
			string(encryptedData),
			vaultKey,
		)
		if err != nil {
			return nil, err
		}

		var record coreEntity.VaultRecord

		err = json.Unmarshal(
			decrypted,
			&record,
		)
		if err != nil {
			return nil, err
		}

		record.ID = recordID

		records = append(
			records,
			record,
		)
	}

	return records, nil
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

	if s.context.SelectedSoT.Address == "localhost" {
		return s.appendRecordLocally(
			vaultAccess,
			record,
		)
	}

	return s.appendRecordRemotely(
		vaultAccess,
		record,
	)
}

func (s *Service) appendRecordLocally(
	vaultAccess *coreEntity.VaultAccess,
	record *coreEntity.VaultRecord,
) (*coreEntity.VaultRecord, error) {
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
	)
	if err != nil {
		return nil, err
	}

	_, err = vaultFile.Insert(
		record,
	)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (s *Service) appendRecordRemotely(
	vaultAccess *coreEntity.VaultAccess,
	record *coreEntity.VaultRecord,
) (*coreEntity.VaultRecord, error) {
	raw, err := s.EncryptRecord(
		vaultAccess,
		record,
	)
	if err != nil {
		return nil, err
	}

	recordID, err := s.appDeps.Client.CreateRecordRaw(
		s.context.SelectedSoT.Address,
		vaultAccess.VaultID,
		raw,
	)
	if err != nil {
		return nil, err
	}

	recordID = strings.TrimSpace(
		recordID,
	)

	record.ID = recordID

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

	if s.context.SelectedSoT.Address == "localhost" {
		return s.updateRecordLocally(
			vaultAccess,
			record,
		)
	}

	return s.updateRecordRemotely(
		vaultAccess,
		record,
	)
}

func (s *Service) updateRecordLocally(
	vaultAccess *coreEntity.VaultAccess,
	record *coreEntity.VaultRecord,
) (*coreEntity.VaultRecord, error) {
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
	)
	if err != nil {
		return nil, err
	}

	err = vaultFile.Update(
		record,
	)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (s *Service) updateRecordRemotely(
	vaultAccess *coreEntity.VaultAccess,
	record *coreEntity.VaultRecord,
) (*coreEntity.VaultRecord, error) {
	raw, err := s.EncryptRecord(
		vaultAccess,
		record,
	)
	if err != nil {
		return nil, err
	}

	err = s.appDeps.Client.UpdateRecordRaw(
		s.context.SelectedSoT.Address,
		vaultAccess.VaultID,
		record.ID,
		raw,
	)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (s *Service) EncryptRecord(
	vaultAccess *coreEntity.VaultAccess,
	record *coreEntity.VaultRecord,
) ([]byte, error) {
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

	vaultKeyEncoded, err := s.context.SelectedSoT.KeyPair.Open(
		vaultAccess.VaultKey,
	)
	if err != nil {
		return nil, err
	}

	vaultKey, err := base64.RawURLEncoding.DecodeString(
		string(vaultKeyEncoded),
	)
	if err != nil {
		return nil, err
	}

	raw, err := json.Marshal(
		record,
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

	encrypted, err := db.Encrypt(
		raw,
		vaultKey,
	)
	if err != nil {
		return nil, err
	}

	return []byte(encrypted), nil
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

	if s.context.SelectedSoT.Address == "localhost" {
		return s.deleteRecordLocally(
			vaultAccess,
			record,
		)
	}

	return s.deleteRecordRemotely(
		vaultAccess,
		record,
	)
}

func (s *Service) deleteRecordLocally(
	vaultAccess *coreEntity.VaultAccess,
	record *coreEntity.VaultRecord,
) error {
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
	)
	if err != nil {
		return err
	}

	_, err = vaultFile.Delete(
		record.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) deleteRecordRemotely(
	vaultAccess *coreEntity.VaultAccess,
	record *coreEntity.VaultRecord,
) error {
	return s.appDeps.Client.DeleteRecord(
		s.context.SelectedSoT.Address,
		vaultAccess.VaultID,
		record.ID,
	)
}
