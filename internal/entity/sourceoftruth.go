package entity

import "github.com/tacenva/tacpass-core/util/keyring"

type SourceOfTruth struct {
	ID        string          `json:"id"`
	Hostname  string          `json:"hostname"`
	Address   string          `json:"address"`
	AuthToken string          `json:"auth_token"`
	KeyPair   keyring.KeyPair `json:"key_pair"`
}

// tacenva/
// ├── local/
// │   ├── config.toml
// │   └── keypair.tacenva
// │
// ├── source-of-truth/
// │   ├── tacenva.db
// │   └── vault/
// │       ├── {vault_id}.tacenva
// │       └── ...
// │
// └── replicas/
//     └── {node_ulid}/
//         ├── tacenva.db
//         └── vault/
//             ├── {vault_id}.tacenva
//             └── ...
