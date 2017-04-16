all:
	make -i clean
	make install-deps
	make build

get-glide:
	curl https://glide.sh/get | sh

install-deps:
	cd src; \
	glide update && glide install

clean:
	rm API
build:
	go build -o API ./src
run:
	./API

docker-build:
	docker build -t naexpire-api .

docker-rm:
	docker rm -f naexpire-instance

docker-run:
	docker run -d -p 80:8000 --name naexpire-instance naexpire-api
