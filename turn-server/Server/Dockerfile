FROM golang:1.15.2-alpine

WORKDIR $GOPATH/src/test-server

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY ./Server .
COPY ./webrtc ./webrtc
RUN go mod tidy
RUN cd $GOPATH/src/test-server/ && go build main.go

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["./main"]