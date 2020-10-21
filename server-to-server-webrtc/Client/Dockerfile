FROM golang:1.15.2-alpine

WORKDIR $GOPATH/src/server-to-server-webrtc

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY ./Client .
COPY ./webrtc ./webrtc

RUN apk add --update g++
RUN go mod tidy
RUN apk --no-cache add curl && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.31.0
RUN $GOPATH/bin/golangci-lint run -v
RUN cd $GOPATH/src/server-to-server-webrtc/ && go build main.go

# This container exposes port 8082 to the outside world
EXPOSE 8082

# Run the executable
CMD ["./main"]