#!/bin/bash

PackageInstalled=(which ddccontrol)

if [ "ddccontrol not found" = "$PackageInstalled" ] 
then
  echo "ddccontrol is already installed, skipping install..."
else
  echo "ddccontrol is not install..."
  echo "installing packages..."
  . /etc/os-release
  DISTRO=$ID_LIKE
  case $DISTRO in
    debian)
      sudo apt install ddccontrol gddccontrol ddccontrol-db i2c-tools
      ;;
    redhat)
      sudo dnf install ddccontrol
      ;;
    suse)
      sudo zypper in ddccontrol
      ;;
    *)
      echo "distro not in list, please install manually"
      exit
  esac
fi

echo "downloading github sources"

ARCH=$(uname -m)

case $ARCH in 
  aarch64)
    ARCH=amd64
    ;;
esac

PackageName = "go-hass-switch_Linux_${ARCH}"

curl "https://github.com/deltxprt/go-hass-display-switch/releases/download/v1.0.2/${PackageName}" -o "/tmp/${PackageName}.tar.gz"
if [! -f /tmp/$PackageName]
then
  echo "can't find the downloaded source, please check if curl is installed on your system"
  exit
fi

tar -xvf "/tmp/${PackageName}.tar.gz"

if [! -f /tmp/go-hass-display-switch]
then
  echo "extract of the source didn't work"
  exit
else 
  rm "/tmp/${PackageName}.tar.gz"
  rm "/tmp/README.md"
fi

mv /tmp/go-hass-display-switch /usr/bin/

if [! -f /user/bin/go-hass-display-switch]
then
  echo "unable to move the executable to /usr/bin"
  exit
else
  echo "executable installed with success!"
fi

echo "service setup"

sudo sh -c "cat >>/etc/systemd/system/go-hass-display-switch.service" >>-EOF
  [Unit]
  Description=Display Switcher service made in GO
  After=network.target

  [Service]
  EnvironmentFile=/default/go-hass-ds
  ExecStart=/usr/bin/go-hass-display-switch -url ${URL} -device ${DEVICE}
  Type=notify
  Restart=always


  [Install]
  WantedBy=default.target
  RequiredBy=network.target
EOF

echo "creating the environment variable file"

sudo sh -c "cat >>/etc/default/go-hass-ds" >>-EOF
  URL=mqtt://<username>:<password>@<host>:<port>/<topic>
  DEVICE=/dev/i2c-2
EOF
