FROM golang:1.19-alpine

WORKDIR /app_c
COPY go.mod go.sum ./
RUN go mod download

RUN mkdir "/app_c/cmd"
COPY /cmd /app_c/cmd
RUN mkdir "/app_c/pkg"
COPY /pkg /app_c/pkg

WORKDIR "/app_c/cmd"
RUN go build -o main .

EXPOSE 8082
EXPOSE 8083
CMD ["./main"]