package entity

import "github.com/tacenva/tacpass-core/util/keyring"

type SourceOfTruth struct {
	ID             string          `json:"id"`
	Hostname       string          `json:"hostname"`
	Address        string          `json:"address"`
	AuthToken      string          `json:"auth_token"`
	KeyPair        keyring.KeyPair `json:"key_pair"`
	TLSFingerprint string          `json:"tls_fingerprint"`
}
