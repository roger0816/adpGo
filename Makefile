

build :
	go build -o bin/adpApi ./cmd/api
	go build -o bin/adpGo ./cmd/socket


clean:
	rm -f bin/adpApi bin/adpGo
	