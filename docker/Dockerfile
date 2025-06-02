FROM golang:1.23

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /profile-store ./cmd

EXPOSE 8080

CMD ["/profile-store"]