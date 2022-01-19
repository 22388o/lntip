FROM golang:1.17 AS builder

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o lntip

FROM alpine
WORKDIR /lntip
COPY --from=builder /app/lntip /lntip/lntip
ENTRYPOINT ["/lntip/lntip"]