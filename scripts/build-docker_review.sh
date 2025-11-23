#!/usr/bin/env bash

set -e

GIT_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "unknown")
TAG=${1:-$GIT_TAG}
REPOSITORY_NAME="$APP_REPOSITORY/$APP_NAME"

if [ "$TAG" == "unknown" ] || [ -z "$TAG" ]; then
  echo "Ошибка: Не удалось определить версию. Укажите версию явно или создайте git тег."
  echo "Использование: $0 [версия]"
  exit 1
fi

echo "Сборка Docker образа с версией: $TAG"

docker buildx build \
  --tag "$REPOSITORY_NAME:$TAG" \
  --build-arg TAG="$TAG" \
  --file build/package/review.Dockerfile .

echo "Docker образ успешно собран: $REPOSITORY_NAME:$TAG"
