#!/usr/bin/env bash
set -e

# перейти в папку со скриптом
cd "$(dirname "$0")"

# если Go ставили через Homebrew, добавим его в PATH (иначе у Finder пустой PATH)
export PATH="/opt/homebrew/bin:/usr/local/go/bin:$PATH"

# запускаем сервер
go run server.go

# чтобы окно Terminal не закрылось
read -r -p "Нажмите Enter для закрытия..."
