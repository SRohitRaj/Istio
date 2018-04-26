# Setting up kubectl on host to talk to kubernetest cluster on Vagrant VM.
echo "your old ~/.kube/config file can be found at ~/.kube/config_old"
cp ~/.kube/config ~/.kube/config_old
vagrant ssh -c "cat ~/.kube/config" > ~/.kube/config
sed -i 's/server: http:\/\/localhost:8080/server: http:\/\/localhost:'"$IstioKport"'/' ~/.kube/config
