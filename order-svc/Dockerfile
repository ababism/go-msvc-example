FROM golang:1.20-bullseye

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

## make wait-for-postgres.sh executable
#RUN chmod +x wait-for-postgres.sh

# build go app
RUN go mod download
RUN go build -o order-svc ./cmd/app/main.go

CMD ["./order-svc"]