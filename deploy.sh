# if ! [[ "$ENV" == "stg" || "$ENV" == "prod" ]]; then
#   echo "Usage ENV=prod|stg ./build.sh"
#   exit 1
# fi

go build -o main cmd/main.go
scp -i ~/ui-network/ui-network.pem main ubuntu@ui-network.com://home/ubuntu/
ssh -i ~/ui-network/ui-network.pem ubuntu@ui-network.com "sudo systemctl stop ${ENV}-ui-web-server"
ssh -i ~/ui-network/ui-network.pem ubuntu@ui-network.com "sudo cp /home/ubuntu/main /opt/webserver/${ENV}/"
ssh -i ~/ui-network/ui-network.pem ubuntu@ui-network.com "sudo systemctl start ${ENV}-ui-web-server"
