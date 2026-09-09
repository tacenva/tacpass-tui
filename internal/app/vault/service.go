package vault

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

	"github.com/tacenva/database"
	"github.com/tacenva/tacpass-core/auth"
	"github.com/tacenva/tacpass-core/entity"
	vaultCore "github.com/tacenva/tacpass-core/vault"
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
	vaultService *vaultCore.Service,
	authService *auth.Service,
) *Service {
	return &Service{
		appDeps:      appDeps,
		context:      context,
		masterKey:    masterKey,
		authService:  authService,
		vaultService: vaultService,
	}
}

func (s *Service) List() ([]entity.VaultAccess, error) {
	if s.context == nil {
		return nil, errors.New("context is nil")
	}

	if s.context.SelectedSoT == nil {
		return nil, errors.New("selected source of truth is nil")
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
		return nil, err
	}

	var vaultAccessList []entity.VaultAccess

	err = vaultFile.FindAll(
		&vaultAccessList,
	)
	if err != nil {
		return nil, err
	}

	if s.context.IsRemote {
		return s.Sync()
	}

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

	vaultAccessPointers := make(
		[]*entity.VaultAccess,
		0,
		len(vaultAccessList),
	)

	for i := range vaultAccessList {
		vaultAccessPointers = append(
			vaultAccessPointers,
			&vaultAccessList[i],
		)
	}

	err = vaultFile.UpdateOrCreateBulk(
		vaultAccessPointers,
	)
	if err != nil {
		return nil, err
	}

	return vaultAccessList, nil
}

func (s *Service) createVaultRemotely(
	vaultName string,
) (*entity.VaultAccess, error) {
	var vaultAccess entity.VaultAccess

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
) (*entity.VaultAccess, error) {
	var (
		vaultAccess *entity.VaultAccess
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
	vault *entity.Vault,
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
		var updated entity.Vault

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

	var vaultAccessList []entity.VaultAccess

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
	vault *entity.Vault,
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

	var vaultAccessList []entity.VaultAccess

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
	vaultAccess *entity.VaultAccess,
) ([]entity.VaultRecord, error) {
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
	vaultAccess *entity.VaultAccess,
) ([]entity.VaultRecord, error) {
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

	var records []entity.VaultRecord

	err = vaultFile.FindAll(
		&records,
	)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (s *Service) listRecordsRemotely(
	vaultAccess *entity.VaultAccess,
) ([]entity.VaultRecord, error) {
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
		[]entity.VaultRecord,
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

		var record entity.VaultRecord

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
	vaultAccess *entity.VaultAccess,
	record *entity.VaultRecord,
) (*entity.VaultRecord, error) {
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
	vaultAccess *entity.VaultAccess,
	record *entity.VaultRecord,
) (*entity.VaultRecord, error) {
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
	vaultAccess *entity.VaultAccess,
	record *entity.VaultRecord,
) (*entity.VaultRecord, error) {
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
	vaultAccess *entity.VaultAccess,
	record *entity.VaultRecord,
) (*entity.VaultRecord, error) {
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
	vaultAccess *entity.VaultAccess,
	record *entity.VaultRecord,
) (*entity.VaultRecord, error) {
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
	vaultAccess *entity.VaultAccess,
	record *entity.VaultRecord,
) (*entity.VaultRecord, error) {
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
	vaultAccess *entity.VaultAccess,
	record *entity.VaultRecord,
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
	vaultAccess *entity.VaultAccess,
	record *entity.VaultRecord,
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
	vaultAccess *entity.VaultAccess,
	record *entity.VaultRecord,
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
	vaultAccess *entity.VaultAccess,
	record *entity.VaultRecord,
) error {
	return s.appDeps.Client.DeleteRecord(
		s.context.SelectedSoT.Address,
		vaultAccess.VaultID,
		record.ID,
	)
}
