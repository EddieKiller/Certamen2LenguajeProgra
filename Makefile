build:
	go build -o simulator main.go dispatcher.go process.go

run:
	./simulator 5 10 order.txt

clean:
	rm -f simulator
