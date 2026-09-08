package vault

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
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

func debugVault(
	format string,
	args ...any,
) {
	file, err := os.OpenFile(
		"debug.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return
	}

	defer file.Close()

	fmt.Fprintf(
		file,
		format,
		args...,
	)

	fmt.Fprintln(file)
}

func (s *Service) List() ([]entity.VaultAccess, error) {
	debugVault("=== List Vault ===")

	if s.context == nil {
		debugVault("Context is nil")
		return nil, errors.New("context is nil")
	}

	if s.context.SelectedSoT == nil {
		debugVault("SelectedSoT is nil")
		return nil, errors.New("selected source of truth is nil")
	}

	debugVault(
		"SoT ID: %s",
		s.context.SelectedSoT.ID,
	)

	debugVault(
		"SoT Address: %s",
		s.context.SelectedSoT.Address,
	)

	debugVault(
		"MasterKey exists: %t",
		s.masterKey != "",
	)

	debugVault(
		"Token exists: %t",
		s.appDeps.Client.Token != "",
	)

	vaultPath := s.appDeps.Config.Path(
		config.NodeDirName,
		s.context.SelectedSoT.ID,
		"vault",
	)

	debugVault(
		"Vault path: %s",
		vaultPath,
	)

	db := database.New(
		vaultPath,
	)

	vaultFile, err := db.File(
		"vault",
		s.masterKey,
	)
	if err != nil {
		debugVault(
			"db.File error: %v",
			err,
		)

		return nil, err
	}

	var vaultAccessList []entity.VaultAccess

	err = vaultFile.FindAll(
		&vaultAccessList,
	)
	if err != nil {
		debugVault(
			"FindAll error: %v",
			err,
		)

		return nil, err
	}

	debugVault(
		"Local vault count: %d",
		len(vaultAccessList),
	)

	if s.context.IsRemote {
		debugVault(
			"No local vault data, syncing from daemon",
		)

		outOfSync, err := s.OutOfSync()
		if err != nil {
			return nil, err
		}

		if outOfSync {
			return s.Sync()
		}
	}

	return vaultAccessList, nil
}

func (s *Service) OutOfSync() (bool, error) {
	vaultPath := s.appDeps.Config.Path(
		config.NodeDirName,
		s.context.SelectedSoT.ID,
		"vault",
	)

	debugVault(
		"Vault path: %s",
		vaultPath,
	)

	db := database.New(
		vaultPath,
	)

	vaultFile, err := db.RawFile(
		"vault",
	)
	if err != nil {
		debugVault(
			"db.RawFile error: %v",
			err,
		)

		return false, err
	}

	hashVal, err := vaultFile.Hash()
	if err != nil {
		debugVault(
			"vaultFile.Hash error: %v",
			err,
		)

		return false, err
	}

	if !s.context.IsRemote {
		return false, nil
	}

	result := struct {
		OutOfSync bool `json:"out_of_sync"`
	}{}

	err = s.appDeps.Client.OutOfSync(
		s.context.SelectedSoT.Address,
		s.context.SelectedSoT.ID,
		hashVal,
		&result,
	)
	if err != nil {
		debugVault(
			"client.OutOfSync error: %v",
			err,
		)

		return false, err
	}

	return result.OutOfSync, nil
}

func (s *Service) Sync() ([]entity.VaultAccess, error) {
	debugVault("=== Sync Vault ===")

	var vaultAccessList []entity.VaultAccess

	err := s.appDeps.Client.ListVaults(
		s.context.SelectedSoT.Address,
		&vaultAccessList,
	)
	if err != nil {
		debugVault(
			"ListVaults error: %v",
			err,
		)

		return nil, err
	}

	debugVault(
		"Remote vault count: %d",
		len(vaultAccessList),
	)

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
		debugVault(
			"db.File error: %v",
			err,
		)

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
		debugVault(
			"UpdateOrCreateBulk error: %v",
			err,
		)

		return nil, err
	}

	debugVault(
		"Sync success",
	)

	return vaultAccessList, nil
}

func (s *Service) createVaultRemotely(
	vaultName string,
) (*entity.VaultAccess, error) {
	debugVault(
		"=== Create Vault Remote ===",
	)

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
		debugVault(
			"CreateVault remote error: %v",
			err,
		)

		return nil, err
	}

	debugVault(
		"CreateVault remote success | vaultID=%s",
		vaultAccess.VaultID,
	)

	return &vaultAccess, nil
}

func (s *Service) CreateVault(
	vaultName string,
) (*entity.VaultAccess, error) {
	debugVault(
		"=== Create Vault ===",
	)

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
			debugVault(
				"GetUserData error: %v",
				err,
			)

			return nil, err
		}

		vaultAccess, err = s.vaultService.Create(
			authUser,
			vaultName,
		)
		if err != nil {
			debugVault(
				"Create local vault error: %v",
				err,
			)

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

	debugVault(
		"CreateVault success | vaultID=%s",
		vaultAccess.VaultID,
	)

	return vaultAccess, nil
}

func (s *Service) ListRecords(
	vaultAccess *entity.VaultAccess,
) ([]entity.VaultRecord, error) {
	debugVault("=== List Records ===")

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

	debugVault(
		"VaultID: %s",
		vaultAccess.VaultID,
	)

	debugVault(
		"SoT Address: %s",
		s.context.SelectedSoT.Address,
	)

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
	debugVault(
		"=== List Records Local ===",
	)

	// IMPORTANT:
	// Local database.File() menggunakan vault key encoded
	// sebagai password. Jangan Base64 Decode di sini.
	vaultKey, err := s.context.SelectedSoT.KeyPair.Open(
		vaultAccess.VaultKey,
	)
	if err != nil {
		debugVault(
			"Open VaultKey error: %v",
			err,
		)

		return nil, err
	}

	debugVault(
		"VaultKey opened successfully | length=%d",
		len(vaultKey),
	)

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
		debugVault(
			"db.File error: %v",
			err,
		)

		return nil, err
	}

	var records []entity.VaultRecord

	err = vaultFile.FindAll(
		&records,
	)
	if err != nil {
		debugVault(
			"FindAll records error: %v",
			err,
		)

		return nil, err
	}

	debugVault(
		"Local records count: %d",
		len(records),
	)

	return records, nil
}

func (s *Service) listRecordsRemotely(
	vaultAccess *entity.VaultAccess,
) ([]entity.VaultRecord, error) {
	debugVault(
		"=== List Records Remote ===",
	)

	var encryptedRecords map[string][]byte

	err := s.appDeps.Client.ListRecords(
		s.context.SelectedSoT.Address,
		vaultAccess.VaultID,
		&encryptedRecords,
	)
	if err != nil {
		debugVault(
			"ListRecords API error: %v",
			err,
		)

		return nil, err
	}

	debugVault(
		"Encrypted records count: %d",
		len(encryptedRecords),
	)

	// HPKE Open menghasilkan encoded vault key 43 karakter.
	vaultKeyEncoded, err := s.context.SelectedSoT.KeyPair.Open(
		vaultAccess.VaultKey,
	)
	if err != nil {
		debugVault(
			"Open VaultKey error: %v",
			err,
		)

		return nil, err
	}

	debugVault(
		"VaultKey opened | encoded length=%d",
		len(vaultKeyEncoded),
	)

	// Remote record menggunakan AES langsung,
	// sehingga kita butuh raw 32-byte key.
	vaultKey, err := base64.RawURLEncoding.DecodeString(
		string(vaultKeyEncoded),
	)
	if err != nil {
		debugVault(
			"Decode VaultKey error: %v",
			err,
		)

		return nil, err
	}

	debugVault(
		"VaultKey decoded | raw length=%d",
		len(vaultKey),
	)

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
		debugVault(
			"Decrypt record | recordID=%s | encrypted length=%d | key length=%d",
			recordID,
			len(encryptedData),
			len(vaultKey),
		)

		decrypted, err := db.Decrypt(
			string(encryptedData),
			vaultKey,
		)
		if err != nil {
			debugVault(
				"Decrypt error | recordID=%s | error=%v",
				recordID,
				err,
			)

			return nil, err
		}

		debugVault(
			"Decrypt success | recordID=%s | decrypted length=%d",
			recordID,
			len(decrypted),
		)

		var record entity.VaultRecord

		err = json.Unmarshal(
			decrypted,
			&record,
		)
		if err != nil {
			debugVault(
				"Unmarshal error | recordID=%s | error=%v",
				recordID,
				err,
			)

			return nil, err
		}

		record.ID = recordID

		records = append(
			records,
			record,
		)
	}

	debugVault(
		"Remote records loaded: %d",
		len(records),
	)

	return records, nil
}

func (s *Service) AppendRecord(
	vaultAccess *entity.VaultAccess,
	record *entity.VaultRecord,
) (*entity.VaultRecord, error) {
	debugVault("=== Append Record ===")

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

	debugVault(
		"VaultID: %s | RecordID: %s",
		vaultAccess.VaultID,
		record.ID,
	)

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
	debugVault(
		"=== Append Record Local ===",
	)

	// Local memakai encoded key sebagai password.
	vaultKey, err := s.context.SelectedSoT.KeyPair.Open(
		vaultAccess.VaultKey,
	)
	if err != nil {
		debugVault(
			"Open VaultKey error: %v",
			err,
		)

		return nil, err
	}

	debugVault(
		"VaultKey opened | length=%d",
		len(vaultKey),
	)

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
		debugVault(
			"db.File error: %v",
			err,
		)

		return nil, err
	}

	_, err = vaultFile.Insert(
		record,
	)
	if err != nil {
		debugVault(
			"Insert record error: %v",
			err,
		)

		return nil, err
	}

	debugVault(
		"Local record inserted",
	)

	return record, nil
}

func (s *Service) appendRecordRemotely(
	vaultAccess *entity.VaultAccess,
	record *entity.VaultRecord,
) (*entity.VaultRecord, error) {
	debugVault(
		"=== Append Record Remote ===",
	)

	raw, err := s.EncryptRecord(
		vaultAccess,
		record,
	)
	if err != nil {
		return nil, err
	}

	debugVault(
		"Encrypted record length: %d",
		len(raw),
	)

	recordID, err := s.appDeps.Client.CreateRecordRaw(
		s.context.SelectedSoT.Address,
		vaultAccess.VaultID,
		raw,
	)
	if err != nil {
		debugVault(
			"CreateRecordRaw error: %v",
			err,
		)

		return nil, err
	}

	recordID = strings.TrimSpace(
		recordID,
	)

	record.ID = recordID

	debugVault(
		"Remote record created | recordID=%s",
		record.ID,
	)

	return record, nil
}

func (s *Service) UpdateRecord(
	vaultAccess *entity.VaultAccess,
	record *entity.VaultRecord,
) (*entity.VaultRecord, error) {
	debugVault("=== Update Record ===")

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

	debugVault(
		"VaultID: %s | RecordID: %s",
		vaultAccess.VaultID,
		record.ID,
	)

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
	debugVault(
		"=== Update Record Local ===",
	)

	// Local memakai encoded key sebagai password.
	vaultKey, err := s.context.SelectedSoT.KeyPair.Open(
		vaultAccess.VaultKey,
	)
	if err != nil {
		debugVault(
			"Open VaultKey error: %v",
			err,
		)

		return nil, err
	}

	debugVault(
		"VaultKey opened | length=%d",
		len(vaultKey),
	)

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
		debugVault(
			"db.File error: %v",
			err,
		)

		return nil, err
	}

	err = vaultFile.Update(
		record,
	)
	if err != nil {
		debugVault(
			"Update record error: %v",
			err,
		)

		return nil, err
	}

	debugVault(
		"Local record updated",
	)

	return record, nil
}

func (s *Service) updateRecordRemotely(
	vaultAccess *entity.VaultAccess,
	record *entity.VaultRecord,
) (*entity.VaultRecord, error) {
	debugVault(
		"=== Update Record Remote ===",
	)

	raw, err := s.EncryptRecord(
		vaultAccess,
		record,
	)
	if err != nil {
		return nil, err
	}

	debugVault(
		"Encrypted record length: %d",
		len(raw),
	)

	err = s.appDeps.Client.UpdateRecordRaw(
		s.context.SelectedSoT.Address,
		vaultAccess.VaultID,
		record.ID,
		raw,
	)
	if err != nil {
		debugVault(
			"UpdateRecordRaw error: %v",
			err,
		)

		return nil, err
	}

	debugVault(
		"Remote record updated | recordID=%s",
		record.ID,
	)

	return record, nil
}

func (s *Service) EncryptRecord(
	vaultAccess *entity.VaultAccess,
	record *entity.VaultRecord,
) ([]byte, error) {
	debugVault(
		"=== Encrypt Record ===",
	)

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
		debugVault(
			"Open VaultKey error: %v",
			err,
		)

		return nil, err
	}

	debugVault(
		"VaultKey opened | encoded length=%d",
		len(vaultKeyEncoded),
	)

	// credential.Generate(32) menghasilkan:
	//
	// 32 raw bytes
	//        ↓
	// Base64 Raw URL
	//        ↓
	// 43 characters
	//
	// db.Encrypt() membutuhkan raw AES key,
	// jadi decode kembali ke 32 bytes.
	vaultKey, err := base64.RawURLEncoding.DecodeString(
		string(vaultKeyEncoded),
	)
	if err != nil {
		debugVault(
			"Decode VaultKey error: %v",
			err,
		)

		return nil, err
	}

	debugVault(
		"VaultKey decoded | raw length=%d",
		len(vaultKey),
	)

	raw, err := json.Marshal(
		record,
	)
	if err != nil {
		debugVault(
			"Marshal record error: %v",
			err,
		)

		return nil, err
	}

	debugVault(
		"Record marshaled | length=%d",
		len(raw),
	)

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
		debugVault(
			"Encrypt error | key length=%d | error=%v",
			len(vaultKey),
			err,
		)

		return nil, err
	}

	debugVault(
		"Encrypt success | encrypted length=%d",
		len(encrypted),
	)

	return []byte(encrypted), nil
}
