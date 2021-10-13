#!/bin/bash
set -euo pipefail

cd test
cp .env-test .env
go test
