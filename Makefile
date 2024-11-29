run:
	go run cmd/web/main.go

buildrun:
	go build -o local cmd/web/main.go && ./local

build_ubuntu:
	GOOS=linux GOARCH=amd64 go build cmd/web/main.go

run_build_ubuntu:
	./main

connect:
	ssh -i "~/.ssh/default.pem" ubuntu@ec2-13-214-123-225.ap-southeast-1.compute.amazonaws.com

deploy_aws:
	GOOS=linux GOARCH=amd64 go build cmd/web/main.go
	scp -i "~/.ssh/default.pem" .env ubuntu@ec2-13-214-123-225.ap-southeast-1.compute.amazonaws.com:/home/ubuntu/app
	scp -i "~/.ssh/default.pem" Makefile ubuntu@ec2-13-214-123-225.ap-southeast-1.compute.amazonaws.com:/home/ubuntu/app
	ssh -i "~/.ssh/default.pem" ubuntu@ec2-13-214-123-225.ap-southeast-1.compute.amazonaws.com "sudo systemctl stop animapu-api"
	scp -i "~/.ssh/default.pem" main ubuntu@ec2-13-214-123-225.ap-southeast-1.compute.amazonaws.com:/home/ubuntu/app
	# ssh -i "~/.ssh/default.pem" ubuntu@ec2-13-214-123-225.ap-southeast-1.compute.amazonaws.com "sudo systemctl start animapu-api"

aws_stop_nohup:
	GOOS=linux GOARCH=amd64 go build cmd/web/main.go
	scp -i "~/.ssh/default.pem" .env ubuntu@ec2-13-214-123-225.ap-southeast-1.compute.amazonaws.com:/home/ubuntu/app
	scp -i "~/.ssh/default.pem" Makefile ubuntu@ec2-13-214-123-225.ap-southeast-1.compute.amazonaws.com:/home/ubuntu/app
	-ssh -i "~/.ssh/default.pem" ubuntu@ec2-13-214-123-225.ap-southeast-1.compute.amazonaws.com "sudo pkill main"
	scp -i "~/.ssh/default.pem" main ubuntu@ec2-13-214-123-225.ap-southeast-1.compute.amazonaws.com:/home/ubuntu/app
	ssh -i "~/.ssh/default.pem" ubuntu@ec2-13-214-123-225.ap-southeast-1.compute.amazonaws.com
	# cd ~/app
	# sudo nohup ./main &
	# rm -rf nohup.out

connect_idcloud:
	ssh umarkotak@103.187.146.246

deploy_idcloud:
	GOOS=linux GOARCH=amd64 go build -o animapu-api cmd/web/main.go
	scp .env umarkotak@103.187.146.246:/home/umarkotak/app/animapu-api
	ssh umarkotak@103.187.146.246 "sudo pkill animapu-api"
	scp animapu-api umarkotak@103.187.146.246:/home/umarkotak/app/animapu-api
	ssh umarkotak@103.187.146.246

default_restart_on_cloud:
	cd ~/app/animapu-api
	sudo nohup ./animapu-api &
	rm -rf nohup.out
