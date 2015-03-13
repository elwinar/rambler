Vagrant.configure(2) do |config|
	config.vm.box = "elwinar/golang"
	config.vm.synced_folder ".", "/home/vagrant/src/github.com/elwinar/rambler"
	config.vm.provision "shell", inline: <<EOS
		chown vagrant:users -R /home/vagrant
EOS
end
