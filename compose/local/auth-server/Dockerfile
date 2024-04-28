FROM golang:1.20.5 as base

FROM base as dev

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

WORKDIR /usr/src/auth-server

RUN git config --global --add safe.directory /usr/src/auth-server

COPY . .

CMD ["air"]