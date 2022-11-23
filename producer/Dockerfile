FROM golang:1.19-alpine

WORKDIR /app_p
COPY go.mod go.sum ./
RUN go mod download

RUN mkdir "/app_p/collector"
COPY /collector /app_p/collector
RUN mkdir "/app_p/parsers"
COPY /parsers /app_p/parsers
RUN mkdir "/app_p/services"
COPY /services /app_p/services

COPY go.mod .
COPY go.sum .
COPY main.go .
RUN go build -o main .

CMD ["./main"]