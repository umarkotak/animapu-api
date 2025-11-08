run:
	go run cmd/web/main.go

bin:
	go build -o animapu-api cmd/web/main.go

bin_run:
./animapu-api

nohup_run:
	nohup ./animapu-api &

stopd:
	pkill animapu-api

statusd:
	ps aux | grep animapu

logs:
	log stream --predicate 'processImagePath contains "animapu-api"' --level debug
