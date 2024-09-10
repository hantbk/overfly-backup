# Storages

- Local
- FTP
- SCP - Upload via SSH copy 
- S3 - Amazon S3

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