# VTS Backup Service with Management Agent

## Features
- Multiple Storage type support.
- Archive paths or files into a tar.
- Split large backup file into multiple parts.
- Run as daemon to backup in schedule.
- Web UI for management.

### Storages

- Local
- FTP
- SFTP
- SCP - Upload via SSH copy 
- S3 - Amazon S3
- MinIO - S3 compatible object storage server

### Compressor

| Type                            | Ext         | Parallel Support |
 |---------------------------------|-------------|------------------|
| `gz`, `tgz`, `taz`, `tar.gz`    | `.tar.gz`   | pigz             |
| `Z`, `taZ`, `tar.Z`             | `.tar.Z`    |                  |
| `bz2`, `tbz`, `tbz2`, `tar.bz2` | `.tar.bz2`  | pbzip2           |
| `lz`, `tar.lz`                  | `.tar.lz`   |                  |
| `lzma`, `tlz`, `tar.lzma`       | `.tar.lzma` |                  |
| `lzo`, `tar.lzo`                | `.tar.lzo`  |                  |
| `xz`, `txz`, `tar.xz`           | `.tar.xz`   | pixz             |
| `zst`, `tzst`, `tar.zst`        | `.tar.zst`  |                  |
| `tar`                           | `.tar`      |                  |
| default                         | `.tar`      |                  |

### Encryptor

- OpenSSL - `aes-256-cbc` encrypt

### Install (macOS / Linux)
```shell
curl -sSL https://raw.githubusercontent.com/hantbk/vtsbackup/master/install | sh
```

## Start 
```shell
vtsbackup -h
```

```
NAME:
   vtsbackup - Backup agent.

USAGE:
   vtsbackup [global options] command [command options]

VERSION:
   0.0.6

COMMANDS:
   perform  Perform backup using config file
   start    Start Backup agent as daemon
   run      Run Backup agent without daemon
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```


## Backup schedule

VtsBackup built in a daemon mode, you can use `vtsbackup start` to start it.

You can configure the `schedule` for each models, it will run backup task at the time you set.

### For example

Configure your schedule in `vtsbackup.yml`

 ```yml
models:
  test-local:
    description: "test backup with local storage"
    schedule:
      cron: "0 0 * * *" # every day at midnight
    archive:
      includes:
        - /Users/hant/Documents
      excludes:
        - /Users/hant/Documents/backup.txt
    compress_with:
      type: tgz
    storages:
      local:
        type: local
        keep: 10
        path: /Users/hant/Downloads/backup1
  test-minio:
    description: "test backup with minio storage"
    schedule:
      every: "1day"
      at: "00:00"
    archive:
      includes:
        - /Users/hant/Documents
    compress_with:
      type: tgz
    encrypt_with:
      type: openssl
      password: 123
      salt: false
      openssl: true
    storages:
      minio:
        type: minio
        bucket: vtsbackup-test
        endpoint: http://127.0.0.1:9000
        path: backups
        access_key_id:
        secret_access_key:
  test-s3:
    description: "test backup with s3 storage"
    schedule:
      every: "180s"
    archive:
      includes:
        - /Users/hant/Documents
    compress_with:
      type: tgz
    storages:
      s3:
        type: s3
        bucket: vts-backup-test
        regions: us-east-1
        path: backups
        access_key_id:
        secret_access_key:
  test-scp:
    description: "test backup with scp storage"
    archive:
      includes:
        - /Users/hant/Documents
    compress_with:
      type: tgz
    storages:
      scp:
        type: scp
        host: 192.168.103.129
        port: 22
        path: ~/backups
        username: hant
        private_key: ~/.ssh/id_rsa

 ```

And then start daemon:
 ```bash
 vtsbackup start
 ```

> NOTE: If you wants start without daemon, use `vtsbackup run` instead.

### Start Daemon & Web UI
 Backup built a HTTP Server for Web UI, you can start it by `vtsbackup start`.

 It also will handle the backup schedule.

 ```bash
 $ vtsbackup start
 Starting API server on port http://127.0.0.1:2703

## Signal handling

VtsBackup will handle the following signals:

- `HUP` - Hot reload configuration.
- `QUIT` - Graceful shutdown.

 ```bash
 $ ps aux | grep vtsbackup
hant             48966   0.3  0.2 411599488  30880   ??  Ss    1:52AM   0:01.41 vtsbackup run
hant             49182   0.0  0.0 410200752   1184 s023  S+    1:56AM   0:00.00 grep --color=auto --exclude-dir=.bzr --exclude-dir=CVS --exclude-dir=.git --exclude-dir=.hg --exclude-dir=.svn --exclude-dir=.idea --exclude-dir=.tox vtsbackup
```

```bash
 # Reload configuration
 $ kill -HUP 48966
 # Exit daemon
 $ kill -QUIT 48966
 ```
