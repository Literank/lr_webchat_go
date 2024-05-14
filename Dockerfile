# Alpine linux https://www.alpinelinux.org/
FROM alpine:3.19

ENV APP_BIN=lr_webchat
ARG SERVER_DIR=/home/.server
WORKDIR $SERVER_DIR
COPY ./bin/${APP_BIN} .

ENV GIN_MODE=release

CMD ./${APP_BIN}
