# tacpass-tui

Tacenva Password Manager terminal user interface.

`tacpass-tui` is a standalone terminal application for managing password vaults locally. It works offline and does not require `tacpassd` for local password management.

For online functionality and remote vault access, `tacpass-tui` can connect to [`tacpassd`](https://github.com/tacenva/tacpassd).

## Download

Pre-built binaries are available in the GitHub Releases.

**[Download the latest release](https://github.com/tacenva/tacpass-tui/releases/latest)**

Choose the binary matching your operating system and architecture.

## Features

* Local password management
* Password vaults
* Add, edit, and view password records
* Password expiration
* Password generation
* Source of Truth management
* Access control
* User management
* Terminal user interface
* Local encrypted storage
* Offline-first operation
* Optional online functionality through `tacpassd`

## Standalone

`tacpass-tui` can be used completely offline.

```text id="q3x7kx"
┌─────────────────┐
│   tacpass-tui   │
│  Terminal App   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  tacpass-core   │
│ Core Services   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│    database     │
│ Encrypted DB    │
└─────────────────┘
```

A running `tacpassd` instance is not required for local vault management.

## Online Mode

When remote or synchronization functionality is required, `tacpass-tui` connects to `tacpassd` over HTTPS.

```text id="gk5d8w"
┌─────────────────┐
│   tacpass-tui   │
│  Terminal App   │
└────────┬────────┘
         │
         │ HTTPS
         ▼
┌─────────────────┐
│    tacpassd     │
│     Daemon      │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  tacpass-core   │
│ Core Services   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│    database     │
│ Encrypted DB    │
└─────────────────┘
```

## Requirements

* Linux, macOS, or Windows
* Terminal emulator

No Go installation is required when using a pre-built release binary.

## Local Data

`tacpass-tui` stores its configuration and local application data under:

```text
~/.tacenva/
```

The main configuration file is:

```text
~/.tacenva/config.toml
```

Local vault data is stored using the encrypted database layer provided by [`database`](https://github.com/tacenva/database).

## Configuration

The application configuration includes information required for local operation and optional remote connectivity.

The Source of Truth identity is generated when the configuration is created and remains associated with the local configuration.

## Online Authentication

When connected to `tacpassd`, the application supports:

* User enrollment
* User authentication
* Bearer token authentication
* Access-controlled vault access
* Remote vault operations
* Vault synchronization

## Development

For development, the project requires Go 1.26.5 or later.

Run tests:

```bash
go test ./...
```

Format the source:

```bash
go fmt ./...
```

Build from source:

```bash
go build ./cmd
```

## Ecosystem

* [`tacpass-tui`](https://github.com/tacenva/tacpass-tui) - Standalone terminal application
* [`tacpassd`](https://github.com/tacenva/tacpassd) - Daemon for online functionality
* [`tacpass-core`](https://github.com/tacenva/tacpass-core) - Core functionality
* [`database`](https://github.com/tacenva/database) - Encrypted database layer

## License

This project is free to use.

See the repository license for details.
