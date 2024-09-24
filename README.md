# VTS Backup Service with Management Agent

## Features
- Multiple Storage type support.
- Archive paths or files into a tar.
- Split large backup file into multiple parts.
- Run as daemon to backup in schedule.
- Web UI for management.

### Storages

- Local
- FTP (`Error in testing`)
- SFTP (`Error in testing`)
- SCP - Upload via SSH copy (`Error in testing`)
- S3 - Amazon S3
- MinIO - S3 compatible object storage server

### Install Linux
```shell
curl -sSL https://raw.githubusercontent.com/hantbk/vtsbackup/master/install | sh
```

### Install Linux with specific version
```shell
curl -sSL https://raw.githubusercontent.com/hantbk/vtsbackup/master/install | sh -s v0.0.1
```

## Usage
```shell
vtsbackup -h
```

```
NAME:
   vtsbackup - Backup agent.

USAGE:
   vtsbackup [global options] command [command options]

VERSION:
   master

COMMANDS:
   perform          Perform backup pipeline using config file
   start            Start Backup agent as daemon
   run              Run Backup agent without daemon
   stop             Stop the running Backup agent
   reload           Reload the running Backup agent
   list-model       List all configured backup models
   list-backup      List backup files for a specific model
   download-backup  Download a backup file for a specific model
   help, h          Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## Config

```bash
nano ~/.vtsbackup/vtsbackup.yml
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
 Starting API server on port http://127.0.0.1:1201

## Signal handling

### List running Backup agents

 ```bash
 $ vtsbackup list
 ```

 ```
Running Backup agents PIDs:
67078
 ```

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

 Or you can use `vtsbackup reload` to reload the configuration.

 ```bash
 $ vtsbackup reload
 ```

 Or you can use `vtsbackup stop` to stop the running Backup agent.

 ```bash
 $ vtsbackup stop
 ```

 ```
Stopping Backup agent...
Running Backup agents PIDs:
67078
Backup agent stopped successfully
 ```

## Install the MinIO Server and Client
Use can use [MinIO](https://min.io) for local development. It is a self-hosted S3-compatible object storage server.
```bash
brew install minio/stable/minio
brew install minio/stable/mc
```
Start MinIO server:
```bash
minio server /tmp/minio
```
And then visit http://localhost:9000 to see the MinIO browser.
The Admin user:
- username: `minioadmin`
- password: `minioadmin`
## Initialize a MinIO bucket
Now we need to create a bucket for testing, we will use the following credentials:
- Bucket: `vtsbackup-test`
- AccessKeyId: `test-user`
- SecretAccessKey: `test-user-secret`
### Configure MinIO Client
Config MinIO Client with a default alias: `minio`
```bash
mc config host add minio http://localhost:9000 minioadmin minioadmin
```
Create a Bucket
```bash
mc mb minio/vtsbackup-test
```
Add Test AccessKeyId and SecretAccessKey.
With
- access_key_id: `test-user`
- secret_access_key: `test-user-secret`
```bash
 mc admin user add minio test-user test-user-secret
 mc admin policy attach minio readwrite --user test-user
 ```

 ## Start Backup in local for MinIO

 ```bash
GO_ENV=dev go run main.go -- perform --config ./tests/minio.yml
 ```

 # Guide for Release new version

 Just create a new tag and push, the GitHub Actions will to the rest.

 ```bash
 git tag v*.*.*
 git push origin v*.*.*
 ```

 After the GitHub Actions finished, the new version will be released to GitHub Releases.


## Check API server 

### Check status

 ```bash
curl http://0.0.0.0:1201/status
 ```

 ```json
 {
 {"message":"Backup is running.","version":"0.0.14"}
 }
 ```

### Get config

 ```bash
curl http://0.0.0.0:1201/api/config
 ```

 ```json
 {"models":{"test-minio":{"description":"test backup with minio storage","schedule":{},"schedule_info":"disabled"}}}
 ```
 
### List backup files for a specific model:

 ```bash
curl "http://0.0.0.0:1201/api/list?model=test-minio"
 ```

 ```json
{"files":[{"filename":"backups/2024.09.21.22.49.41.tar.gz","size":422,"last_modified":"2024-09-21T15:49:41.339Z"},{"filename":"backups/2024.09.21.22.19.30.tar","size":10240,"last_modified":"2024-09-21T15:19:30.494Z"}]}
 ```

### Download a backup file:

 ```bash
curl -O -J -L "http://0.0.0.0:1201/api/download?model=test-minio&file=backups/2024.09.21.22.49.41.tar.gz"
 ```

 ```
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100    41  100    41    0     0  33198      0 --:--:-- --:--:-- --:--:-- 41000
 ```

### Perform a backup:

 ```bash
curl -X POST -H "Content-Type: application/json" -d '{"model":"test-minio"}' http://0.0.0.0:1201/api/perform
 ```

 ```json
{"message":"Backup: test-minio performed in background."}
 ```

### Get log stream

 ```bash
curl http://0.0.0.0:1201/api/log
 ```

 ```
2024/09/21 22:51:16 [Config] Load config from default path.
2024/09/21 22:51:16 [Config] Other users are able to access /root/.vtsbackup/vtsbackup.yml with mode -rw-r--r--
2024/09/21 22:51:16 [Config] Config file: /root/.vtsbackup/vtsbackup.yml
2024/09/21 22:51:16 [Config] Config loaded, found 1 models.
2024/09/21 22:51:16 [API] You are running with insecure API server. Please don't forget setup `web.password` in config file for more safety.
2024/09/21 22:51:16 [API] Starting API server on port http://0.0.0.0:1201
2024/09/21 22:59:37 [Model: test-minio] WorkDir: /tmp/backup1412014860/1726933876879263691/test-minio
2024/09/21 22:59:37 [Archive] => includes 1 rules
2024/09/21 22:59:37 [Compressor] => Compress: tgz
2024/09/21 22:59:37 [Compressor] -> /tmp/backup1412014860/1726933876879263691/2024.09.21.22.59.37.tar.gz
2024/09/21 22:59:37 [Storage] => Storage: minio
2024/09/21 22:59:37 [MinIO] -> Uploading (422 B)...
2024/09/21 22:59:37 [MinIO] Uploaded: http://127.0.0.1:9000/vtsbackup-test/backups/2024.09.21.22.59.37.tar.gz (Duration 51 milliseconds 496 microseconds)
2024/09/21 22:59:37 [Model] Cleanup temp: /tmp/backup1412014860/
2024/09/21 23:00:28 [Model] Failed to seek log file: EOF
 ```

