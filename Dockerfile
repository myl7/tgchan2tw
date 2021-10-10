FROM golang:alpine
WORKDIR /app
COPY . /app
RUN apk add gcc musl-dev wget tar gzip \
 && wget -O twitter-text-parse-go.tar.gz https://github.com/myl7/twitter-text-parse-go/archive/4ef36e65e0b6ad532d8aa77413453310659f141d.tar.gz \
 && tar xzf twitter-text-parse-go.tar.gz \
 && mv twitter-text-parse-go-4ef36e65e0b6ad532d8aa77413453310659f141d twitter-text-parse-go \
 && mkdir twitter-text-parse-go/lib \
 && wget -O twitter-text-parse-go/lib/libtwitter_text_parse_go.a twitter-text-parse-go https://github.com/myl7/twitter-text-parse-go/releases/download/prebuilt/libtwitter_text_parse_go-x86_64-unknown-linux-musl.a \
 && mv twitter-text-parse-go .. \
 && go build -o tgchan2tw cmd/tgchan2tw/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=0 /app/tgchan2tw /app/tgchan2tw
RUN mkdir /db
ENTRYPOINT ["./tgchan2tw"]
