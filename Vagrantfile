# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure(2) do |config|
  config.vm.box = "ubuntu/trusty64"
  config.vm.network "forwarded_port", guest: 8080, host: 8085
  config.vm.provision :shell, path: "vagrant/vagrant_bootstrap.sh"
  config.vm.synced_folder "../../../../.", "/srv/go"
end
