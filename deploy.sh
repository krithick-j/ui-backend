go build -o main cmd/main.go
scp -i ~/ui-network/ui-network.pem main ubuntu@ui-network.com://home/ubuntu/
ssh -i ~/ui-network/ui-network.pem ubuntu@ui-network.com "sudo systemctl stop prod-ui-web-server"
ssh -i ~/ui-network/ui-network.pem ubuntu@ui-network.com "sudo cp /home/ubuntu/main /opt/webserver/prod/"
ssh -i ~/ui-network/ui-network.pem ubuntu@ui-network.com "sudo systemctl start prod-ui-web-server"
