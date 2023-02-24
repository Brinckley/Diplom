FROM golang:1.19-alpine

WORKDIR /app_ef
COPY go.mod go.sum ./
RUN go mod download

RUN mkdir "/app_ef/cmd"
COPY /cmd /app_ef/cmd
RUN mkdir "/app_ef/internal"
COPY /internal /app_ef/internal

WORKDIR "/app_ef/cmd/publisher"
RUN go build -o main .

CMD ["./main"]