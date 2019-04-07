all: build

build:
	go build

clean:
	rm ./aacautoupdate

install: aacautoupdate
	cp aacautoupdate /usr/local/bin/aacautoupdate
	mkdir -p /var/www/cell.bdavidson.dev/data

remove:
	rm /usr/local/bin/aacautoupdate

purge: remove
	rm -r /var/www/cell.bdavidson.dev/data