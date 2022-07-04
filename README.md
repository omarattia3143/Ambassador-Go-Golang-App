# GoAndNextProject
**A simple CMS for products built using Go, Fiber, React and Next
This project is a backend for** https://github.com/omarattia3143/react-admin .

To run project using **docker compose**:

```
version: '3.9'
services:
  backend:
    build: .
    ports:
      - "8000:3000"
    volumes:
      - .:/app
    depends_on:
      - db
      - redis

  db:
    image: mysql
    restart: always
    environment:
      MYSQL_DATABASE: ambassador
      MYSQL_ROOT_PASSWORD: root
    volumes:
      - .dbdata:/var/lib/mysql
    ports:
      - "3306:3306"
  redis:
    image: redis:latest
    ports:
      - "6379:6379"
```

also **dockerfile** for mailhog (smtp server) but in dockerfile to avoid some complications:

```
ARG GOTAG=1.18-alpine
FROM golang:${GOTAG} as builder
MAINTAINER CD2Team <codesign2@icloud.com>

RUN set -x \
  && buildDeps='git musl-dev gcc' \
  && apk add --update $buildDeps \
  && GOPATH=/tmp/gocode go install github.com/mailhog/MailHog@latest

FROM alpine:latest
WORKDIR /bin
COPY --from=builder tmp/gocode/bin/MailHog /bin/MailHog
EXPOSE 1025 8025
ENTRYPOINT ["MailHog"]
```

also future refactoring will be done to this project (moving repo and bussniess logic to service and repository layers)
