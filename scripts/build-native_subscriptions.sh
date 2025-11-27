#!/usr/bin/env bash

set -e

GIT_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "unknown")

if [ "$GIT_TAG" == "unknown" ]; then
  echo "Ошибка: Не удалось определить git тег. Создайте тег в репозитории."
  exit 1
fi

echo "Нативная сборка приложения с версией: $GIT_TAG"

go build -o build/package/subscriptions-"$GIT_TAG" -ldflags="-X main.version=$GIT_TAG" cmd/sub/main.go

cd build/package
ln -sf subscriptions-"$GIT_TAG" subscriptions-latest
cd ../..

echo "Приложение успешно собрано: build/package/subscriptions-$GIT_TAG"
echo "Создан симлинк: build/package/subscriptions-latest -> subscriptions-$GIT_TAG"
