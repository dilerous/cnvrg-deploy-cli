apt-get -y install git
git config --global user.email "brad.soper@gmail.com"

apt-get update
apt install -y build-essential libssl-dev libreadline-dev zlib1g-dev shared-mime-info
apt install -y rbenv
rbenv init
eval "$(rbenv init -)"
mkdir -p "$(rbenv root)"/plugins
git clone https://github.com/rbenv/ruby-build.git "$(rbenv root)"/plugins/ruby-build
rbenv install 2.6.9
rbenv global 2.6.9
gem install cnvrg --no-document
rbenv rehash
cnvrg version
echo 'eval "$(rbenv init -)"' >> ~/.bashrc
