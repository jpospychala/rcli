all: build test

build:
	go build R.go

test:
	./test.sh

install:
	cp ./R $(DESTDIR)/usr/bin/R

uninstall:
	rm -f $(DESTDIR)/usr/bin/R

clean:
	rm -f ./R

deb_package:
	rm -f ../rcli-0.1.tar.gz
	tar -czvf rcli-0.1.tar.gz *
	mkdir rcli-0.1
	cd rcli-0.1 && tar -zxvf ../rcli-0.1.tar.gz
	cd rcli-0.1 && dh_make -c bsd -e 'jacek.pospychala@gmail.com' -f ../rcli-0.1.tar.gz -s
