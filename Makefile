

build :
	go build -o bin/adpApi ./cmd/api
	go build -o bin/adpSocket ./cmd/socket


clean:
	rm -f bin/adpApi bin/adpSocket
	