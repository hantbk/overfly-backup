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

