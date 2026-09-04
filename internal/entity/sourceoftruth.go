package entity

type KeyPair struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

type SourceOfTruth struct {
	ID        string  `json:"id"`
	Hostname  string  `json:"hostname"`
	Address   string  `json:"address"`
	AuthToken string  `json:"auth_token"`
	KeyPair   KeyPair `json:"key_pair"`
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
