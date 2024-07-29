.PHONY: run clean

clean:
	rm -rf ./.bin
	rm *.json

run:
	@go build -o ./.bin/ main.go doc.go && ./.bin/main

record:
	@vhs demo.tape
