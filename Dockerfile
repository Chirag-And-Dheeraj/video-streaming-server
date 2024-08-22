FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

RUN go install github.com/air-verse/air@latest

RUN apt-get -y update && apt-get -y upgrade && apt-get install -y ffmpeg

COPY ./ ./

CMD ["air"]
