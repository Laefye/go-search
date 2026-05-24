FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o server ./cmd/server
RUN go build -o listener ./cmd/listener
RUN go build -o cleaner ./cmd/cleaner

FROM alpine:latest AS server

WORKDIR /app

COPY --from=builder /app/server .

CMD ["./server"]


FROM alpine:latest AS listener

WORKDIR /app

COPY --from=builder /app/listener .

CMD ["./listener"]


FROM alpine:latest AS cleaner

WORKDIR /app

COPY --from=builder /app/cleaner .

CMD ["./cleaner"]
