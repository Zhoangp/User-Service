FROM golang:1.19-alpine as builder

COPY .  /app/
WORKDIR /app/
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o courses-user-service .

FROM alpine
WORKDIR /app/
COPY --from=builder /app/courses-user-service .
COPY config/*.yml ./config/
COPY wait-for .
CMD [ "./courses-user-service" ]
