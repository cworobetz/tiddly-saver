
FROM golang:stretch AS build-env

WORKDIR $GOPATH/src/tiddly-saver

RUN apt-get update
RUN apt-get install gcc libgtk-3-dev libappindicator3-dev -y

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN GO111MODULE=on go build -o /go/bin/tiddly-saver

FROM scratch
COPY --from=build-env /go/bin/tiddly-saver /srv/tiddly-saver

WORKDIR /srv
ENTRYPOINT [ "./tiddly-saver" ]