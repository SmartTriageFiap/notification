FROM golang:latest AS builder

WORKDIR /app
COPY . .

ENV USER=appuser
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/app" \
    --shell "/sbin/nologin" \
    --no-create-home \
    "${USER}"

RUN go mod download
RUN CGO_ENABLED=0 go build -o /app/bin/main

FROM builder AS test
RUN go test ./... -coverprofile=cover.html

FROM scratch
COPY --from=test /app/cover.html /
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /app/bin/main /app/bin/main
EXPOSE 8080
USER appuser:appuser

ENTRYPOINT ["/app/bin/main"]