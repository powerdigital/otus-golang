# Собираем в гошке
FROM golang:1.17.8 as build

ENV BIN_FILE /opt/project/project-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o /opt/project/project-app ./cmd/main.go

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="project"
LABEL MAINTAINERS="powerdigital2010@gmail.ru"

ENV BIN_FILE "opt/project/project-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENTRYPOINT /opt/project
