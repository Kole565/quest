#!/usr/bin/env bash
set -e

# Собираем под Linux amd64
GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o dist/locator-linux-amd64 .

# Опционально — под Windows и macOS
# GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o dist/locator-windows-amd64.exe .
# GOOS=darwin  GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o dist/locator-darwin-amd64 .
# GOOS=darwin  GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o dist/locator-darwin-arm64 .

# Копируем картинку рядом
cp map.png dist/

echo "Готово: dist/"
ls -la dist/
