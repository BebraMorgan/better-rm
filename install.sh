#!/bin/bash

APP_NAME="brm"
BINARY_NAME="brm"
SRC_DIR="$(pwd)"
GO_CMD=$(command -v go 2>/dev/null)

echo "Запуск установки $APP_NAME..."

if [ -z "$GO_CMD" ]; then
	echo "Установите Go перед продолжением." >&2
	exit 1
fi

echo "Найден Go: $GO_CMD"

cd "$SRC_DIR" || {
	echo "Не могу перейти в директорию $SRC_DIR" >&2
	exit 1
}

echo "Собираю проект..."
go build -o "$BINARY_NAME" || {
	echo "❌ Ошибка при сборке проекта" >&2
	exit 1
}

DEST="/usr/local/bin/$BINARY_NAME"

echo "Копирую $BINARY_NAME в $DEST..."
sudo mv "$BINARY_NAME" "$DEST" || {
	echo "Не могу скопировать файл в $DEST — проверьте права" >&2
	exit 1
}

sudo chmod +x "$DEST"

if command -v "$BINARY_NAME" &>/dev/null; then
	echo "Установка завершена! Теперь вы можете запустить $BINARY_NAME из терминала."
else
	echo "Что-то пошло не так — команда $BINARY_NAME не найдена в PATH"
	exit 1
fi
