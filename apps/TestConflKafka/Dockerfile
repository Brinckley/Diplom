FROM golang:alpine
WORKDIR /app_c
COPY go.mod go.sum ./
RUN go mod download

RUN mkdir "/app_c/cmd"
COPY /cmd /app_c/cmd

WORKDIR "/app_c/cmd"
RUN go build -o main .

CMD ["./main"]