## Install the MinIO Server and Client
Use can use [MinIO](https://min.io) for local development. It is a self-hosted S3-compatible object storage server.

## Install in MacOS
```bash
brew install minio/stable/minio
brew install minio/stable/mc
```

## Install in Linux AMD64
### Install MinIO Server
```bash
wget https://dl.min.io/server/minio/release/linux-amd64/archive/minio_20240913202602.0.0_amd64.deb -O minio.deb
sudo dpkg -i minio.deb
```

### Install MinIO Client
```bash
wget https://dl.min.io/client/mc/release/linux-amd64/mc
chmod +x mc
sudo mv mc /usr/local/bin/mc
```

## Install in Linux ARM64
### Install MinIO Server
```bash
wget https://dl.min.io/server/minio/release/linux-arm64/archive/minio_20240913202602.0.0_arm64.deb -O minio.deb
sudo dpkg -i minio.deb
```

### Install MinIO Client
```bash
wget https://dl.min.io/client/mc/release/linux-arm64/mc
chmod +x mc
sudo mv mc /usr/local/bin/mc
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