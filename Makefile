build: clean
	go build -o "pvz-cli-app" cmd/main.go

run: build 
	./pvz-cli-app  

clean: 
	rm -rf pvz-cli-app