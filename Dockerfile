FROM golang:1.22

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY ./ ./

RUN go mod tidy

RUN apt-get -y update && apt-get -y upgrade && apt-get install -y ffmpeg

CMD ["air"]
