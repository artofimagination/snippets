FROM golang:1.15.2-alpine

WORKDIR $GOPATH/src/load-tester

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

RUN apk add --update g++ git
RUN git clone https://github.com/artofimagination/mysql-user-db-go-interface $GOPATH/src/user-db-mysql && cp -r $GOPATH/src/user-db-mysql/db $GOPATH/src/load-tester/ && rm -fr $GOPATH/src/user-db-mysql
RUN go mod tidy
RUN apk --no-cache add curl && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.31.0
RUN $GOPATH/bin/golangci-lint run -v
RUN cd $GOPATH/src/load-tester/ && go build main.go

# This container exposes port 8082 to the outside world
EXPOSE 8081

# Run the executable
CMD [ "./main" ]