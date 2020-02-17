NS = rq
IMAGEUSER = jemccorm
IMAGE = my-custom-test
IMAGE_VERSION = v0.0.1
build-custom-test:
	go build -o bin/$(IMAGE) ./pkg/custom
	sudo podman build --tag quay.io/$(IMAGEUSER)/$(IMAGE):$(IMAGE_VERSION) -f ./Dockerfile
	sudo podman push --authfile /home/jeffmc/.docker/config.json $(IMAGEUSER)/$(IMAGE):$(IMAGE_VERSION) docker://quay.io/$(IMAGEUSER)/$(IMAGE):$(IMAGE_VERSION)
clean:   
	rm ./bin/$(IMAGE)
run:   
	sudo podman run --name my-custom-test --rm quay.io/jemccorm/my-custom-test:v0.0.1

