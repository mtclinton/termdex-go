FROM golang:alpine AS builder

WORKDIR /pokedex

RUN apk update && apk upgrade && apk add --no-cache ca-certificates gcc alpine-sdk
RUN update-ca-certificates

# Download Go modules
COPY * .
RUN go mod download

# Build
RUN CGO_ENABLED=1 GOOS=linux go build -o /termdex -a -ldflags '-linkmode external -extldflags "-static"' .

FROM debian:buster-slim
COPY sprites/ /sprites
COPY --from=builder /termdex /termdex
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/


#ENTRYPOINT ["/termdex-go"]
RUN chmod +x /termdex
# Run
CMD ["/termdex"]
