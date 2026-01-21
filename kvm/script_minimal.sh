#!/bin/bash
set -e

# Minimal bootable Linux/QEMU hello world system
# Runs entirely in working directory, prints "hello world" on boot

WORK_DIR="$(pwd)/qemu_boot"
#KERNEL_VERSION="6.19-rc5"
KERNEL_VERSION="6.18.5"
KERNEL_SRC_DIR="linux-"${KERNEL_VERSION}
#KERNEL_URL="https://git.kernel.org/torvalds/t/"${KERNEL_SRC_DIR}".tar.xz"
KERNEL_URL="https://cdn.kernel.org/pub/linux/kernel/v6.x/"${KERNEL_SRC_DIR}".tar.xz"
BUSYBOX_URL="https://busybox.net/downloads/busybox-1.36.1.tar.bz2"

mkdir -p "$WORK_DIR"
cd "$WORK_DIR"



# Download and extract kernel if not present
if [ ! -f ${KERNEL_SRC_DIR}/arch/x86_64/boot/bzImage ]; then
    echo "Building kernel..."
    wget -N "$KERNEL_URL"
    tar xf ${KERNEL_SRC_DIR}.tar.xz
    
    cd ${KERNEL_SRC_DIR}
    pwd
    cp ../../kernel_config .config
    
    make olddefconfig -j$(nproc)
    make bzImage -j$(nproc)
    cd ..
fi

if [ ! -d initrd ]; then
   mkdir initrd
fi
gcc -o initrd/init -static ../kvm_hw.c
( cd initrd/ && find . | cpio -o -H newc ) | gzip > initrd.gz
cd ${WORK_DIR}
# Boot with QEMU
echo "Booting Linux system..."
kvm -kernel ${KERNEL_SRC_DIR}/arch/x86_64/boot/bzImage \
	 -initrd initrd.gz \
	 -append 'console=hvc0 quiet loglevel=0' \
	 -chardev stdio,id=stdio,mux=on\
	 -device virtio-serial-pci \
	 -device virtconsole,chardev=stdio \
	 -mon chardev=stdio \
	 -display none
echo "QEMU session ended."

