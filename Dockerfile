FROM golang:1.27.0 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o snip .


FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/snip .

CMD ["./snip"]