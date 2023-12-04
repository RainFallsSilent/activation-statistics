all:
	 go build -o activation main.go

linux:
	GOOS=linux GOARCH=amd64 go build -o activation main.go

format:
	go fmt ./*

clean:
	rm -rf *.8 *.o *.out *.6

