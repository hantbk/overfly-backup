version='0.0.4'

 platform='linux'

 os_uname=`uname`
 if [[ "$os_uname" == 'Linux' ]]; then
    platform='linux'
 elif [[ "$os_uname" == 'Darwin' ]]; then
    platform='darwin'
 fi

curl -Lo vts-backup.tar.gz https://github.com/hantbk/vts-backup/releases/download/$version/vts-backup-$platform-arm64.tar.gz
tar zxf vts-backup.tar.gz && sudo mv vts-backup /usr/local/bin/vts-backup && rm vts-backup.tar.gz
mkdir -p ~/.vts-backup && touch ~/.vts-backup/config.yml