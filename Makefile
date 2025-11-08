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

install-service:
	sudo cp com.animapu-api.plist /Library/LaunchDaemons
	sudo chmod 644 /Library/LaunchDaemons/com.animapu-api.plist
	sudo launchctl load /Library/LaunchDaemons/com.animapu-api.plist

uninstall-service:
	sudo launchctl unload /Library/LaunchDaemons/com.animapu-api.plist
	sudo rm /Library/LaunchDaemons/com.animapu-api.plist

start:
	sudo launchctl start com.animapu-api

stop:
	sudo launchctl stop com.animapu-api
