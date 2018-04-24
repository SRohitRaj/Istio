#!/bin/bash

# Setup vagrant.
echo "Setup vagrant"
vagrant up --provider virtualbox

# Setting up host to talk to insecure registry on VM.
echo "Adding insecure registry to docker daemon in host vm..."
echo "You old docker daemon file can be found at /lib/systemd/system/docker.service_old"
sudo cp /lib/systemd/system/docker.service /lib/systemd/system/docker.service_old
echo "sudo sed -i 's/ExecStart=\/usr\/bin\/dockerd -H fd:\/\//ExecStart=\/usr\/bin\/dockerd -H fd:\/\/ --insecure-registry 10.10.0.2:5000/' /lib/systemd/system/docker.service"
sudo sed -i 's/ExecStart=\/usr\/bin\/dockerd -H fd:\/\//ExecStart=\/usr\/bin\/dockerd -H fd:\/\/ --insecure-registry 10.10.0.2:5000/' /lib/systemd/system/docker.service
sudo systemctl daemon-reload
sudo systemctl restart docker

# Set route rule to access VM
echo "Set route rule to access VM"
sudo ip route add 10.0.0.0/24 via 10.10.0.2

echo "$(tput setaf 1)Please run docker info and make sure insecure registry address is updated to 10.10.0.2:5000$(tput sgr 0)"
