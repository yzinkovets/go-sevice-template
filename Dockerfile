FROM golang:1.22.5-alpine as builder
LABEL stage=gobuilder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8888
CMD ["./main"]