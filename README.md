# VTS Backup Service with Management Agent

[![Go Report Card](https://goreportcard.com/badge/github.com/hantbk/vtsbackup)](https://goreportcard.com/report/github.com/hantbk/vtsbackup)
[![GoDoc](https://godoc.org/github.com/hantbk/vtsbackup?status.svg)](https://godoc.org/github.com/hantbk/vtsbackup)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GitHub release](https://img.shields.io/github/release/hantbk/vtsbackup.svg)](https://github.com/hantbk/vtsbackup/releases/)
[![Go version](https://img.shields.io/github/go-mod/go-version/hantbk/vtsbackup.svg)](https://github.com/hantbk/vtsbackup)
[![CI Status](https://github.com/hantbk/vtsbackup/actions/workflows/ci.yml/badge.svg)](https://github.com/hantbk/vtsbackup/actions/workflows/ci.yml)
[![Release Status](https://github.com/hantbk/vtsbackup/actions/workflows/release.yml/badge.svg)](https://github.com/hantbk/vtsbackup/actions/workflows/release.yml)

## Features
- 🔄 Register agent with control plane (`In Progress`)
- 💾 Backup entire Linux system
- 🔄 Restore from backup 
- 🗄️ Multiple Storage type support (Local, FTP, SFTP, SCP, MinIO, S3)
- 📦 Archive paths or files into a tar
- 🔐 Encrypt backup file with OpenSSL
- 🗜️ Compress backup file with gzip
- 📂 Split large backup file into multiple parts
- ⏰ Run as daemon to backup on schedule
- 🖥️ Web UI for management
- 📱 Telegram Notifier

## Quick Start

### Install on Linux

```shell
export MINIO_ACCESS_KEY_ID=test-user
export MINIO_SECRET_ACCESS_KEY=test-user-secret
curl -sSL https://raw.githubusercontent.com/hantbk/vtsbackup/master/install | sh
```

### Install Linux with specific version
```shell
export MINIO_ACCESS_KEY_ID=test-user
export MINIO_SECRET_ACCESS_KEY=test-user-secret
curl -sSL https://raw.githubusercontent.com/hantbk/vtsbackup/master/install | sh -s v0.0.1
```

## Documentation

- [📘 Usage Guide](./docs/usage.md)
- [🛠️ Development Environment Setup](./docs/minio-setup.md)
- [🚀 Release Process](./docs/release.md)
- [🔍 Agent Health Check](./docs/check-agent.md)
- [🔐 Encrypt and Compress](./docs/encrypt-compress.md)
- [🔧 Control Panel](./docs/control-panel.md)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the [Apache License 2.0](LICENSE).

