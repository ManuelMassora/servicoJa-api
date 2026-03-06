FROM golang:1.24-alpine AS builder

RUN apk update && apk upgrade && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main main.go

FROM alpine:3.21

RUN apk update && apk upgrade && apk add --no-cache ca-certificates tzdata
RUN adduser -D -g '' appuser

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080

ENV PORT=8080
ENV GIN_MODE=release

CMD ["./main"]