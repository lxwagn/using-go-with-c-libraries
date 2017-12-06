
all:
	cd src; make dynamic 
	go build -o bin/cgo cgo.go
	bin/cgo
