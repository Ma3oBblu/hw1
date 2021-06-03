# docker build -t mailgo_hw1 .
FROM golang:1.15.12
COPY . .
RUN go test -v