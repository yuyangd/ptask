FROM golang:1.15.3-alpine3.12 AS builder

WORKDIR /src

COPY . .

RUN CGO_ENABLED=0 go build -o /dist/ptask .

FROM scratch AS export-stage

COPY --from=builder /dist/ptask /

ENTRYPOINT [ "/ptask" ]
