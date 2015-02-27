Vagrant.configure(2) do |config|
	config.vm.box = "elwinar/golang"
	config.vm.synced_folder ".", "/home/vagrant/go/src/github.com/elwinar/rambler"
	config.vm.provision "shell", inline: <<EOS
		pacman -S git --noconfirm
		chown vagrant:users -R /home/vagrant/go
EOS
end
