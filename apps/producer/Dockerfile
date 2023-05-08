FROM golang:1.19-alpine

WORKDIR /app_p
COPY go.mod go.sum ./
RUN go mod download

RUN mkdir "/app_p/cmd"
COPY /cmd /app_p/cmd
RUN mkdir "/app_p/internal"
COPY /internal /app_p/internal
RUN mkdir "/app_p/test"
COPY /test /app_p/test
COPY /input.txt /app_p/input.txt

COPY go.mod .
COPY go.sum .
WORKDIR "/app_p/cmd/producer"
RUN go build -o main .

CMD ["./main"]