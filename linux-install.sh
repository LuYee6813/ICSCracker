# /bin/bash
apt install golang-go -y
go build -o ICSCracker ICSCracker.go
sudo mv ICSCracker /usr/local/bin/
sudo chmod +x /usr/local/bin/ICSCracker