FROM golang:1.15.2-alpine

WORKDIR $GOPATH/src/golang-docker-deploy

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .
RUN go mod tidy
RUN cd $GOPATH/src/golang-docker-deploy/ && go build main.go

# This container exposes port 8082 to the outside world
EXPOSE 8081

# Run the executable
CMD [ "./main" ]