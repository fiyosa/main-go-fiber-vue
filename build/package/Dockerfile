FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY app .
COPY public/ ./public/
EXPOSE 8000
CMD ["./app"]
