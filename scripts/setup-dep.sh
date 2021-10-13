#!/bin/bash
set -euo pipefail

# $1 should be gnu/musl

wget -qO twitter-text-parse-go-1.0.0.tar.gz https://github.com/myl7/twitter-text-parse-go/archive/v1.0.0.tar.gz
tar xzf twitter-text-parse-go-1.0.0.tar.gz
mkdir twitter-text-parse-go-1.0.0/lib
wget -qO twitter-text-parse-go-1.0.0/lib/libtwitter_text_parse_go.a https://github.com/myl7/twitter-text-parse-go/releases/download/v1.0.0/libtwitter_text_parse_go-x86_64-unknown-linux-$1.a
mkdir -p third_party
mv twitter-text-parse-go-1.0.0 third_party/twitter-text-parse-go
