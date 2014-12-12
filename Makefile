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

docs:
	head -`grep -A 1 -n '^Functions$$' README.md | tail -1 | cut -f1 -d'-'` < README.md > README.md.tmp
	./R help | \
		sed -n '4,$$p' | \
		cut -f '1' -d' ' | \
		sort | \
		while read F; do \
			echo; echo $$F; \
			echo $$F|sed 's/./-/g'; \
			Z=`./R help $$F | grep -v "Usage:"`; \
			echo "$$Z" | tr '\n' '!' | sed 's/Example:\(.*\)/Example:!```bash\1```/' | tr '!' '\n'; \
		done >> README.md.tmp
	mv README.md.tmp README.md
