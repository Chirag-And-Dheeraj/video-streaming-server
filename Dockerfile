FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN apt-get -y update && apt-get -y upgrade && apt-get install -y ffmpeg

COPY ./ ./

EXPOSE 8000

CMD [ "go", "run", "main.go" ]