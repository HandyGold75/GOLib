get:
	go get go@latest
	go get -u
	go mod tidy
	for indirect in $(tail +3 go.mod | grep "// indirect" | awk '{if ($1 =="require") print $2; else print $1;}'); do go get -u "${indirect}"; done
	go get -u
	go mod tidy

	for subdir in *"/"; do \
		if [ ! -f "$(shell pwd)/$$subdir"go.mod ]; then \
			continue ; \
		fi ; \
		cd "$(shell pwd)/$$subdir" ; \
		go get go@latest ; \
		go get -u ; \
		go mod tidy ; \
		for indirect in $(tail +3 go.mod | grep "// indirect" | awk '{if ($1 =="require") print $2; else print $1;}'); do go get -u "${indirect}"; done ; \
		go get -u ; \
		go mod tidy ; \
		echo $$subdir ; \
	done
