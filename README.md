# Storages

- Local
- FTP
- SCP - Upload via SSH copy 
- S3 - Amazon S3

### Compressor

- Tgz - `.tar.gz`

### Encryptor

- OpenSSL - `aes-256-cbc` encrypt

# Usage
```bash
go build
./vts-backup perform
```

## Install
```bash
curl -sSL https://raw.githubusercontent.com/hantbk/vts-backup/master/install.sh | bash
```

## Schedule run

You may want run backup in scheduly, you need Crontab:

 ```bash
 $ crontab -l
 0 0 * * * /usr/local/bin/vts-backup perform >> ~/.vts-backup/vts-backup.log
 ```

> `0 0 * * *` means run at 0:00 AM, every day.
And after a day, you can check up the execute status by `~/.vts-backup/vts-backup.log`.
>

```bash
go run main.go -- perform -m demo -c ./vtsbackup_test.yml
```

```
2024/09/11 01:04:14 ======== demo ========
2024/09/11 01:04:14 WorkDir: /var/folders/vf/f0kp2vw524z_2h9yfq6pdwnw0000gn/T/vtsbackup/1725991454175656000/demo

2024/09/11 01:04:14 ------------- Archives -------------
2024/09/11 01:04:14 ------------- Archives -------------
2024/09/11 01:04:14 => includes 4 rules
2024/09/11 01:04:14 [debug] /usr/bin/tar -cPf /var/folders/vf/f0kp2vw524z_2h9yfq6pdwnw0000gn/T/vtsbackup/1725991454175656000/demo/archive.tar --exclude=/home/ubuntu/.ssh/known_hosts --exclude=/etc/logrotate.d/syslog /home/ubuntu/.ssh /etc/nginx/nginx.conf /etc/redis/redis.conf /etc/logrotate.d
2024/09/11 01:04:14 ------------- Archives -------------

2024/09/11 01:04:14 ------------- Archives -------------

2024/09/11 01:04:14 ------------ Compressor -------------
2024/09/11 01:04:14 => Compress | tgz
2024/09/11 01:04:14 -> /var/folders/vf/f0kp2vw524z_2h9yfq6pdwnw0000gn/T/vtsbackup/1725991454175656000/2024.09.11.01.04.14.tar.gz
2024/09/11 01:04:14 ------------ Compressor -------------

2024/09/11 01:04:14 ------------ Encryptor -------------
2024/09/11 01:04:14 => Encrypt | openssl
2024/09/11 01:04:14 -> /var/folders/vf/f0kp2vw524z_2h9yfq6pdwnw0000gn/T/vtsbackup/1725991454175656000/2024.09.11.01.04.14.tar.gz.enc
2024/09/11 01:04:14 ------------ Encryptor -------------

2024/09/11 01:04:14 ------------- Storage --------------
2024/09/11 01:04:14 ------------- Storage --------------
2024/09/11 01:04:14 => Storage | local
2024/09/11 01:04:14 Store successed /Users/hant/Downloads/backup1
2024/09/11 01:04:14 ------------- Storage --------------

2024/09/11 01:04:14 ------------- Storage --------------

2024/09/11 01:04:14 Cleanup temp: /var/folders/vf/f0kp2vw524z_2h9yfq6pdwnw0000gn/T/vtsbackup/1725991454175656000/

2024/09/11 01:04:14 ======= End demo =======
```