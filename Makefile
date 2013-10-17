all: check-env clean install

install:
	go install github.com/czertbytes/tierheimdb/piggybank
	go install github.com/czertbytes/tierheimdb/catnip
	go install github.com/czertbytes/tierheimdb/catnip/archenoah
	go install github.com/czertbytes/tierheimdb/catnip/berlin
	go install github.com/czertbytes/tierheimdb/catnip/dellbrueck
	go install github.com/czertbytes/tierheimdb/catnip/dresden
	go install github.com/czertbytes/tierheimdb/catnip/frankfurtmain
	go install github.com/czertbytes/tierheimdb/catnip/franziskushamburg
	go install github.com/czertbytes/tierheimdb/catnip/muenchen
	go install github.com/czertbytes/tierheimdb/catnip/samtpfotenneukoelln
	go install github.com/czertbytes/tierheimdb/watchdog
	go install github.com/czertbytes/tierheimdb/kennel
	go install github.com/czertbytes/tierheimdb/parade

clean: check-env
	rm -rf $(GOPATH)/bin
	rm -rf $(GOPATH)/pkg

check-env:
	if test "$(GOPATH)" = "" ; then \
		echo "GOPATH is not set"; \
		exit 1; \
	fi
