mkdir initrd
gcc -o initrd/init -static kvm_hw.c
( cd initrd/ && find . | cpio -o -H newc ) | gzip > initrd.gz
sudo kvm -kernel /boot/vmlinuz -initrd initrd.gz -append 'console=hvc0 quiet loglevel=0' -chardev stdio,id=stdio,mux=on -device virtio-serial-pci -device virtconsole,chardev=stdio -mon chardev=stdio -display none
