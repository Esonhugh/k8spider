FROM golang:1.21 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLE=0 go build -ldflags="-s -w" -tags netgo -o k8spider main.go

RUN 	./k8spider all --help && \
    	./k8spider axfr --help  && \
    	./k8spider dns --help && \
    	./k8spider dnssd --help && \
    	./k8spider dnssd ptr --help && \
    	./k8spider dnssd srv --help && \
    	./k8spider metrics --help && \
    	./k8spider neighbor --help && \
    	./k8spider neighbor pod --help && \
    	./k8spider neighbor svc --help && \
    	./k8spider whereisdns --help && \
    	./k8spider wild --help && \
        echo "all command flag is successfully processed with no conflict"

FROM alpine:latest
LABEL maintainer="Esonhugh"
LABEL name="k8spider"
LABEL url="https://github.com/esonhugh/k8spider"
LABEL description="A Powerful+Fast+Low Privilege Kubernetes discovery tools"
COPY --from=builder /app/k8spider /k8spider
ENTRYPOINT ["/k8spider"]
CMD ["all"]