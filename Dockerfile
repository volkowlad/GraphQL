FROM golang:1.24
ENV GOPATH=/

ARG STORAGE=memory
ENV STORAGE=${STORAGE}

COPY ./ ./

RUN go mod download
RUN go build -o post-app ./app/main.go

CMD ["./post-app"]