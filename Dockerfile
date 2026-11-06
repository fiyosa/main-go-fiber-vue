FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY app .

COPY lang/ ./lang/
COPY public/ ./public/

EXPOSE 8080

CMD ["./app"]
