build:
	go build -ldflags="-s -w" -o bin/

clean:
	rm -rf ./bin
