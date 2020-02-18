NS = rq
IMAGEUSER = jemccorm
IMAGE = my-custom-test
IMAGE_VERSION = v0.0.1
compile:
	go build -o bin/$(IMAGE) ./pkg/custom
build-image: compile
	sudo podman build --tag quay.io/$(IMAGEUSER)/$(IMAGE):$(IMAGE_VERSION) -f ./Dockerfile
push-image: 
	sudo podman push --authfile /home/jeffmc/.docker/config.json $(IMAGEUSER)/$(IMAGE):$(IMAGE_VERSION) docker://quay.io/$(IMAGEUSER)/$(IMAGE):$(IMAGE_VERSION)
clean:   
	rm ./bin/$(IMAGE)
test:   
	kubectl -n rq delete pod my-custom-test-pod
	kubectl -n rq create -f manifests/pod.yaml
run:   
	sudo podman run --name my-custom-test --rm quay.io/jemccorm/my-custom-test:v0.0.1

