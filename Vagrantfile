Vagrant.configure(2) do |config|
  config.vm.box = "ogarcia/archlinux-x64"
  config.vm.synced_folder ".", "/vagrant"
  config.vm.synced_folder ".", "/home/vagrant/src/github.com/elwinar/rambler"
  config.vm.provision "shell", path: "etc/provision.sh"
  config.vm.provision "shell", path: "etc/start.sh", run: "always"
end
