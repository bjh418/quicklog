build:
	go build -o quicklog main.go

start: stop clean
	./quicklog &

stop:
	killall quicklog || echo "."

clean:
	[ ! -e quicklogger.sock ] || rm quicklogger.sock
	[ ! -e quicktailer.sock ] || rm quicktailer.sock

tail:
	nc -U quicktailer.sock
