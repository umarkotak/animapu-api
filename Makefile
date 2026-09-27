run:
	go run cmd/web/main.go

migrate-up:
	go run cmd/web/main.go migrate up

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
	tail -f animapu-api.log

install-service:
	@test -x animapu-api || { echo "Build the binary with make bin first" >&2; exit 1; }
	@plist=$$(mktemp); trap 'rm -f "$$plist"' EXIT; \
		sed 's|@APP_DIR@|$(CURDIR)|g' com.animapu-api.plist > "$$plist" && \
		plutil -lint "$$plist" && \
		sudo install -m 644 "$$plist" /Library/LaunchDaemons/com.animapu-api.plist
	sudo launchctl bootstrap system /Library/LaunchDaemons/com.animapu-api.plist

uninstall-service:
	@if sudo launchctl print system/com.animapu-api >/dev/null 2>&1; then \
		sudo launchctl bootout system/com.animapu-api; \
	fi
	sudo rm -f /Library/LaunchDaemons/com.animapu-api.plist

start:
	sudo launchctl start com.animapu-api

stop:
	sudo launchctl stop com.animapu-api

deploy:
	git pull --rebase origin master
	go mod tidy
	go mod vendor
	$(MAKE) bin
	$(MAKE) uninstall-service
	sudo rm -f animapu-api.log animapu-api.error.log
	$(MAKE) install-service

status:
	sudo lsof -i :33000

db-tunnel:
	cloudflared access tcp --hostname pg.cabocil.com --url localhost:54322
