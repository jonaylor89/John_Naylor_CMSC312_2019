
APP=OS

build:
	go build -o ${APP} .

run:
	go run -race main.go

clean:
	rm -rf ${APP}