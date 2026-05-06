run:
	go run .

migrate-up:
	go run . migrate up

bin:
	go build -o animapu-api cmd/web/main.go

bin_run:
	./animapu-api

nohup_run:
	nohup ./animapu-api &

stopd:
	pkill animapu-api

statusd:
	ps aux | grep animapu-api

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

deploy:
	git pull --rebase origin master
	go mod tidy
	go mod vendor
	make bin
	make stop
	make start

status:
	sudo lsof -i :33000

db-tunnel:
	cloudflared access tcp --hostname pg.cabocil.com --url localhost:54322