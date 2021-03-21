FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN groupadd --gid 1001 seidhr \
&& useradd -g seidhr --uid 1001 seidhr

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o seidhr ./cmd/main.go



FROM alpine:latest

COPY --from=builder /etc/passwd /etc/passwd

USER seidhr

COPY --from=builder /app/seidhr /app/

EXPOSE 8443

ENTRYPOINT ["/app/seidhr"] 