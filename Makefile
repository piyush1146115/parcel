redis-run:
	docker run --name redis --net=host  -d redis:7-alpine

redis-pull:
	docker pull redis:7-alpine

redis-ping:
	docker exec -it redis redis-cli ping

# Define the name of the image and the Docker Hub repository
IMAGE_NAME := parcel-image
DOCKER_REPO := piyush1146115/parcel

build:
    # Build the Docker image with the specified name
	docker build -t $(IMAGE_NAME) .

push:
    # Tag the Docker image with the Docker Hub repository name
	docker tag $(IMAGE_NAME) $(DOCKER_REPO):latest

    # Push the Docker image to Docker Hub
	docker push $(DOCKER_REPO):latest

install-dependencies: redis-pull redis-run

install:
	# Pull the latest image from Docker hub
	docker pull $(DOCKER_REPO):latest
	# Run the image locally
	docker run --name parcel-simulator --net=host -d $(DOCKER_REPO):latest

.PHONY: redis-pull redis-run redis-ping build push install install-dependencies