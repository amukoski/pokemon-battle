FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /pokemon-battle ./main.go

FROM alpine:3.19
RUN apk --no-cache add ca-certificates wget
COPY --from=builder /pokemon-battle /pokemon-battle
COPY migrations/ /migrations/
EXPOSE 8080
CMD ["/pokemon-battle"]
