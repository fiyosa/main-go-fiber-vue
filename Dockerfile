FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY app .

COPY lang/ ./lang/
COPY public/ ./public/

# copy file .env "scripts/deploy.sh:20"
COPY .env .

EXPOSE 8080

CMD ["./app"]
