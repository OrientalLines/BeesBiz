FROM alpine:latest

# Install Go 1.23.1
RUN apk add --no-cache go=1.23.2-r0 --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community

RUN apk --no-cache add ca-certificates

ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV PATH="${GOPATH}/bin:${PATH}"

WORKDIR /app

COPY . .

RUN go mod download

ENV PATH="/usr/local/lib:${PATH}"

RUN GOOS=linux CGO_ENABLED=0 go build -o /usr/local/bin/main ./cmd/api/main.go

ENTRYPOINT ["main"]
