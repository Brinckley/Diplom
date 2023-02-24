FROM golang:1.19-alpine

WORKDIR /app_c
COPY go.mod go.sum ./
RUN go mod download

RUN mkdir "/app_c/cmd"
COPY /cmd /app_c/cmd
RUN mkdir "/app_c/internal"
COPY /internal /app_c/internal

WORKDIR "/app_c/cmd/consumer"
RUN go build -o main .

CMD ["./main"]