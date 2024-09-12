# VTS Backup Service with Management Agent

## Features
- No dependencies.
- Multiple Storage type support.
- Archive paths or files into a tar.
- Split large backup file into multiple parts.
- Run as daemon to backup in schedule.

### Storages

- Local
- FTP
- SFTP
- SCP - Upload via SSH copy 
- S3 - Amazon S3
- WebDAV - For Synology NAS 

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

### Notifier
Send notification when backup has success or failed.
- Github 
- Mail - Send email via SMTP
- Postmark - Send email via Postmark API

### Usage
```bash
go build
./vtsbackup perform
```

### Install (macOS / Linux)
```shell
curl -sSL https://raw.githubusercontent.com/hantbk/vts-backup/master/install | sh
```

## Start 
```shell
go run main.go -h
```

```
NAME:
   vtsbackup - Backup agent.

USAGE:
   vtsbackup [global options] command [command options]

VERSION:
   master

COMMANDS:
   perform  
   start    Start as daemon
   run      Run VtsBackup
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
   my_backup:
     schedule:
       # At 04:05 on Sunday.
       cron: "5 4 * * sun"
     storages:
       local:
         type: local
         path: /path/to/backups
   other_backup:
     # At 04:05 on every day.
     schedule:
       every: "1day",
       at: "04:05"
     storages:
       local:
         type: local
         path: /path/to/backups
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
 hant            20443   0.0  0.1 409232800   8912   ??  Ss    7:47PM   0:00.02 vtsbackup run
 # Reload configuration
 $ kill -HUP 20443
 # Exit daemon
 $ kill -QUIT 20443
 ```
