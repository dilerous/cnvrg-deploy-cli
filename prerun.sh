apt-get -y install git
git config --global user.email "brad.soper@gmail.com"

apt-get update
apt install -y build-essential libssl-dev libreadline-dev zlib1g-dev
apt install -y rbenv
rbenv init
eval "$(rbenv init -)"
mkdir -p "$(rbenv root)"/plugins
git clone https://github.com/rbenv/ruby-build.git "$(rbenv root)"/plugins/ruby-build
rbenv install 2.6.9
rbenv global 2.6.9
export FREEDESKTOP_MIME_TYPES_PATH=/cnvrg/freedesktop.org.xml
gem install cnvrg --no-document
rbenv rehash
cnvrg version