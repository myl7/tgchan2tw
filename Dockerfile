FROM golang:alpine
WORKDIR /app
COPY . /app
RUN apk add gcc musl-dev wget tar gzip
RUN bash scripts/setup-dep.sh musl
RUN go build -o tgchan2tw cmd/tgchan2tw/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=0 /app/tgchan2tw /app/tgchan2tw
RUN mkdir /db
ENTRYPOINT ["./tgchan2tw"]
