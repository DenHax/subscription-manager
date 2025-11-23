#!/bin/bash

set -e

GIT_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "unknown")

if [ "$GIT_TAG" == "unknown" ]; then
  echo "Ошибка: Не удалось определить git тег. Создайте тег в репозитории."
  exit 1
fi

echo "Нативная сборка приложения с версией: $GIT_TAG"

mkdir -p build/bin

go build -o build/bin/review-"$GIT_TAG" -ldflags="-X main.version=$GIT_TAG" cmd/review/main.go

cd build/bin
ln -sf review-"$GIT_TAG" review-latest
cd ../..

echo "Приложение успешно собрано: build/bin/review-$GIT_TAG"
echo "Создан симлинк: build/bin/review-latest -> review-$GIT_TAG"
