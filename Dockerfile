FROM golang:1.22

WORKDIR /app

COPY ./ ./

RUN go mod tidy

RUN apt-get -y update && apt-get -y upgrade && apt-get install -y ffmpeg

CMD ["go",  "run", "main.go"]
