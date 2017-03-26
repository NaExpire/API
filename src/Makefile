all:
	make -i clean
	make install-deps
	make build

get-glide:
	curl https://glide.sh/get | sh

install-deps:
	glide up

clean:
	rm API
build:
	go build -o API
run:
	./API

docker:
	docker build -t naexpire-api .