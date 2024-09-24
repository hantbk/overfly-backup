# VTS Backup Service with Management Agent

## Features
- Multiple Storage type support (Local, FTP, SFTP, SCP, MinIO, S3)
- Archive paths or files into a tar.
- Encrypt backup file with openssl.
- Compress backup file with gzip.
- Split large backup file into multiple parts.
- Run as daemon to backup in schedule.
- Web UI for management.
- Telegram Notifier.

### Install Linux
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

## [Usage](./doc/usage.md)
## [Set up Development Environment](./doc/minio-setup.md)
## [Release](./doc/release.md)
## [Check Agent](./doc/check-agent.md)

