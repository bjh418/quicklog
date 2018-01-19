build:
	go build -o quicklog main.go

start: stop
	./quicklog &

stop:
	killall quicklog || echo "."
	[ ! -e quicklogger.sock ] || rm quicklogger.sock
	[ ! -e quicktailer.sock ] || rm quicktailer.sock

tail:
	nc -U quicktailer.sock
