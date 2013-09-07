all: check-env clean get-apps

get-apps:
	go get github.com/czertbytes/tierheimdb/catnip
	go get github.com/czertbytes/tierheimdb/catnip/franziskushamburg
	go get github.com/czertbytes/tierheimdb/catnip/samtpfotenneukoelln
	go get github.com/czertbytes/tierheimdb/catnip/tierheimberlin
	go get github.com/czertbytes/tierheimdb/catnip/tierheimdresden
	go get github.com/czertbytes/tierheimdb/catnip/tierheimmuenchen
	go get github.com/czertbytes/tierheimdb/kennel
	go get github.com/czertbytes/tierheimdb/parade
	go get github.com/czertbytes/tierheimdb/piggybank
	go get github.com/czertbytes/tierheimdb/prism
	go get github.com/czertbytes/tierheimdb/watchdog

clean: check-env
	rm -rf $(GOPATH)/bin
	rm -rf $(GOPATH)/pkg

check-env:
	if test "$(GOPATH)" = "" ; then \
		echo "GOPATH is not set"; \
		exit 1; \
	fi
