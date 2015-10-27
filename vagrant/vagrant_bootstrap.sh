#!/usr/bin/env bash

## Install software
# Install python software tools, add key for mariadb repo
apt-get update
apt-get install -y python-software-properties
apt-key adv --recv-keys --keyserver hkp://keyserver.ubuntu.com:80 0xcbcb082a1bb943db

# Add repositories for MariaDB 10.0
add-apt-repository -y 'deb http://nyc2.mirrors.digitalocean.com/mariadb/repo/10.0/ubuntu trusty main'

# Set the root password for MariaDB install
export DEBIAN_FRONTEND=noninteractive
sudo debconf-set-selections <<< 'mariadb-server-10.0 mysql-server/root_password password a'
sudo debconf-set-selections <<< 'mariadb-server-10.0 mysql-server/root_password_again password a'

# Update apt-get, install software
apt-get update
apt-get install -y mariadb-server npm nodejs-legacy git

# Install Go
wget https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.5.1.linux-amd64.tar.gz

echo 'export GOPATH=/srv/go' >> "/home/vagrant/.bashrc"
echo 'export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin' >> "/home/vagrant/.bashrc"
