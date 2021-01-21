FROM golang:1.14-alpine as dev

WORKDIR /go/src/Geeksbot
COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .
RUN go install github.com/dustinpianalto/geeksbot/...

CMD [ "go", "run", "cmd/geeksbot/main.go"]

from alpine

WORKDIR /bin

COPY --from=dev /go/bin/geeksbot ./geeksbot
COPY --from=dev /go/src/Geeksbot/internal/database/migrations/* ./internal/database/migrations/

CMD [ "geeksbot" ]
