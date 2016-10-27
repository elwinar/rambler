Vagrant.configure(2) do |config|
	config.vm.box = "kaorimatz/archlinux-x86_64"
	config.vm.synced_folder ".", "/home/vagrant/src/github.com/elwinar/rambler"
	config.vm.provision "shell", inline: <<EOS
		chown vagrant:users -R /home/vagrant

		echo "Server = https://mirror.compojoom.com/archlinux/$repo/os/$arch
Server = http://mirror.archlinux.ikoula.com/archlinux/$repo/os/$arch
Server = http://archlinux.de-labrusse.fr/$repo/os/$arch
Server = http://mir.archlinux.fr/$repo/os/$arch
Server = http://archlinux.mirror.root.lu/$repo/os/$arch
" >> /etc/pacman.d/mirrorlist 
		pacman -Syy --noconfirm

		pacman -S --noconfirm go git subversion mercurial bzr
		echo "export GOPATH=~
		export PATH=\$GOPATH/bin:\$PATH
		" >> /home/vagrant/.bashrc
EOS
end
