#!/bin/bash
# Make sure only root can run our script
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

echo ">> Installing TFTP server and Expect"
apt-get update
apt-get install -y xinetd tftpd tftp tar expect
echo ">>>> Creating service"
cat ./service-tftp.conf > /etc/xinetd.d/tftp
echo ">>>> Setting up directory"
useradd -r -s /bin/false tftpuser
mkdir /tftpboot
chmod -R 777 /tftpboot
chown -R nobody /tftpboot
echo ">>>> Starting server"
/etc/init.d/xinetd restart
echo ">>>> TFTP is now installed"
echo ""
echo ">> Creating icauser"
useradd icauser
echo ">>>> User created"
exit
