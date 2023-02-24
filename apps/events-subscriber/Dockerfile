FROM golang:alpine

WORKDIR /app_es
COPY go.mod go.sum ./
RUN go mod download

RUN mkdir "/app_es/cmd"
COPY /cmd /app_es/cmd
RUN mkdir "/app_es/internal"
COPY /internal /app_es/internal

WORKDIR "/app_es/cmd/subscriber"
RUN go build -o main .

CMD ["./main"]