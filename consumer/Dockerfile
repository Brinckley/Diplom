FROM golang:1.19-alpine

WORKDIR /app_c
COPY go.mod go.sum ./
RUN go mod download

RUN mkdir "/app_c/structs-lastfm"
COPY /parsers/structs-lastfm /app_c/structs-lastfm
RUN mkdir "/app_c/structs-discogs"
COPY /parsers/structs-discogs /app_c/structs-discogs
RUN mkdir "/app_c/services"
COPY /services /app_c/services

COPY main.go .
RUN go build -o main .

CMD ["./main"]