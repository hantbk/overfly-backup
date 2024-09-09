curl -Lo vts-backup-linux-arm64.tar.gz https://github.com/hant/vts-backup/releases/download/0.0.1/vts-backup-linux-arm64.tar.gz
 tar zxf vts-backup-linux-arm64.tar.gz && sudo mv vts-backup /usr/local/bin/vts-backup && rm vts-backup-linux-arm64.tar.gz
 mkdir -p ~/.vts-backup && touch ~/.vts-backup/vts-backup.yml