protoc-setup:
	cd meshes
	wget https://raw.githubusercontent.com/layer5io/meshery/master/meshes/meshops.proto

proto:	
	protoc -I meshes/ meshes/meshops.proto --go_out=plugins=grpc:./meshes/

docker:
	docker build -t layer5/meshery-osm .

docker-run:
	(docker rm -f meshery-osm) || true
	docker run --name meshery-osm -d \
	-p <port>:<port> \
	-e DEBUG=true \
	layer5/meshery-osm

run:
	DEBUG=true go run main.go

# setup-adapter sets up a new adapter with the given name & port
setup-adapter:
	mv "osm" ${ADAPTER}
	find . -type f -exec sed -i '' -e 's/osm/${ADAPTER}/g' {} +
	find . -type f -exec sed -i '' -e 's/<port>/${PORT}/g' {} +
	find . -type f -exec sed -i '' -e 's/<go_version>/${GO_VERSION}/g' {} +
	