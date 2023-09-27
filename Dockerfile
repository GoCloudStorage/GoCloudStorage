#FROM golang:latest AS builder
#
#WORKDIR /src
#
#COPY . .
#
#RUN go build -o storageSvc ./logic/local/api

FROM ubuntu:latest

WORKDIR /app

#COPY --from=builder /src/storageSvc .
#COPY --from=builder /src/etc/config.yaml .

COPY ./storageSvc .
COPY ./etc/config.yaml .

CMD ./storageSvc -c ./config.yaml