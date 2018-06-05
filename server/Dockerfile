FROM golang:1.10-alpine AS go-wasm
RUN apk add --no-cache git bash
RUN git clone --depth 1 --branch wasm-wip https://github.com/neelance/go.git /go-wasm
WORKDIR /go-wasm/src
#ENV GOOS=linux
#ENV GOARCH=amd64
ENV CGO_ENABLED=0
RUN ./make.bash
# Need to skip os/exec/TestCredentialNoSetGroups, os/user/TestCurrent, os/user/TestLookup, os/user/TestLookupId
# https://github.com/golang/go/issues/18243
# RUN ./run.bash --no-rebuild
ENV PATH="/go-wasm/bin:${PATH}"
ENV GOROOT=/go-wasm
ENV GOOS=js
ENV GOARCH=wasm

FROM go-wasm AS clientbuilder
COPY lib /go/src/github.com/vsekhar/go-e2e/lib
COPY client /client/
WORKDIR /client
RUN go get -d ./...
RUN go build -o app.wasm .

FROM golang:1.10-alpine AS serverbuilder
COPY server /go/src/server
RUN go install server

FROM alpine:latest
COPY --from=serverbuilder /go/bin/server .
COPY static /static/
COPY --from=clientbuilder /client/app.wasm /static/app.wasm
ENV PORT 8080
CMD ["./server", "-path=/static"]
