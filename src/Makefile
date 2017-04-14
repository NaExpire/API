all:
	make -i clean
	make install-deps
	make build

get-glide:
	curl https://glide.sh/get | sh

install-deps:
	glide update && glide install

clean:
	rm API
build:
	go build -o API 
run:
	./API

docker-build:
	make install-deps
	docker build -t naexpire-api .

docker-run:
	docker run -i -t -p 8000:8000 naexpire-api
