FROM golang:1.15.2-alpine

WORKDIR $GOPATH/src/test-client

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY ./Client .
COPY ./webrtc ./webrtc

RUN go mod tidy
RUN cd $GOPATH/src/test-client/ && go build main.go

# This container exposes port 8082 to the outside world
EXPOSE 8082

# Run the executable
CMD ["./main"]