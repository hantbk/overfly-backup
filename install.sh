version='0.0.3'
if [[ `uname` == 'Darwin' ]]; then
   platform='darwin'
else
   platform='linux'
fi
curl -Lo vtsbackup.tar.gz https://github.com/hantbk/vtsbackup/releases/download/$version/vtsbackup-$platform-arm64.tar.gz
tar zxf vtsbackup.tar.gz

if [[ `whoami` == 'root' ]]; then
   mv vtsbackup /usr/local/bin/vtsbackup
else
   sudo mv vtsbackup /usr/local/bin/vtsbackup
fi
mkdir -p ~/.vtsbackup && touch ~/.vtsbackup/vtsbackup.yml
rm vtsbackup.tar.gz
