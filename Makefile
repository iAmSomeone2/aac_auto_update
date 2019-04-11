name = aacautoupdate
ver = 1.1

all: build

build:
	go build
	@mkdir build && mv ./aacautoupdate build/aacautoupdate
	@echo aacautoupdate binary can be found in ./build.

pack: build 
	@cp install.sh build/ && cp uninstall.sh build/
	@cp service/aacautoupdate.service build/
	@tar -zcvf $(name)-$(ver).tar.gz ./build/*
	@echo Created $(name)-$(ver).tar.gz from ./build


clean:
	rm -r ./build
	rm $(name)-$(ver).tar.gz

install: aacautoupdate
	cp aacautoupdate /usr/local/bin/
	cp ./aacautoupdate.service /etc/systemd/system/
	mkdir -p /var/www/.cache && chown -hR www-data:www-data /var/www/.cache
	mkdir -p /var/www/cell.bdavidson.dev/html/data/
	chown -R www-data:www-data /var/www/cell.bdavidson.dev/html/data/
	chmod -R 764 /var/www/cell.bdavidson.dev/html/data/
	systemctl enable aacautoupdate && systemctl start aacautoupdate
	@echo aacautoupdate installed and active!

uninstall:
	rm /usr/local/bin/aacautoupdate
	systemctl disable aacautoupdate
	rm -r /etc/systemd/system/aacautoupdate.service
	@echo aacautoupdate disabled and removed.

purge: uninstall
	rm -r /var/www/cell.bdavidson.dev/html/data/
	rm -r /var/www/.cache/adopt-a-cell/