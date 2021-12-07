#
# Copyright contributors to the Galasa project 
#
ARG dockerRepository
FROM ${dockerRepository}/dockerhub/library/golang:latest AS builder

WORKDIR $GOPATH/src/github.com/galasa-dev/galasa-kubernetes-operator
COPY . $GOPATH/src/github.com/galasa-dev/galasa-kubernetes-operator

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=auto
ENV GOPATH=/go

RUN go build -o $GOPATH/bin/operator $GOPATH/src/github.com/galasa-dev/galasa-kubernetes-operator/cmd/operator/main.go

FROM ${dockerRepository}/dockerhub/library/alpine:latest
COPY --from=builder /go/bin/operator /go/bin/operator
ENTRYPOINT ["/go/bin/operator"]
