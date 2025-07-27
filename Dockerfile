FROM golang:1.21

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -v -o /otp-service main.go

EXPOSE 8080

CMD ["/otp-service"]
