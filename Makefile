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
	tail -f animapu-api.error.log

install-service:
	sudo chmod +x animapu-api
	sudo cp com.animapu-api.plist /Library/LaunchDaemons
	sudo chmod +x /Library/LaunchDaemons/com.animapu-api.plist
	sudo launchctl bootstrap system /Library/LaunchDaemons/com.animapu-api.plist

uninstall-service:
	sudo launchctl unload /Library/LaunchDaemons/com.animapu-api.plist
	sudo rm /Library/LaunchDaemons/com.animapu-api.plist

start:
	sudo launchctl start com.animapu-api

stop:
	sudo launchctl stop com.animapu-api

status:
	sudo lsof -i :6001
