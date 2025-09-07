FROM golang:1.24.7-bookworm AS base

FROM base AS dev

RUN go install github.com/air-verse/air@latest

WORKDIR /usr/src/app

COPY . .

CMD ["air"]