run:
	go run cmd/web/main.go

build_ubuntu:
	GOOS=linux GOARCH=amd64 go build cmd/web/main.go

run_build_ubuntu:
	./main

connect:
	ssh -i "~/.ssh/default.pem" ubuntu@ec2-13-214-123-225.ap-southeast-1.compute.amazonaws.com

deploy_aws:
	scp -i "~/.ssh/default.pem" .env ubuntu@ec2-13-214-123-225.ap-southeast-1.compute.amazonaws.com:/home/ubuntu/app
	scp -i "~/.ssh/default.pem" main ubuntu@ec2-13-214-123-225.ap-southeast-1.compute.amazonaws.com:/home/ubuntu/app
