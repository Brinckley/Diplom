FROM golang:1.19-alpine

WORKDIR /app_ch
COPY go.mod go.sum ./
RUN go mod download

RUN mkdir "/app_ch/cmd"
COPY /cmd /app_ch/cmd
RUN mkdir "/app_ch/pkg"
COPY /pkg /app_ch/pkg

WORKDIR "/app_ch/cmd/"
RUN go build -o main .

CMD ["./main"]