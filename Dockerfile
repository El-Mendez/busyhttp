FROM golang:1.23.5-alpine AS build

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /go/bin/app -v .


FROM busybox
COPY --from=build /go/bin/app /app/bin

RUN adduser busyhttp -u 4000 -D
USER busyhttp

EXPOSE 8080
ENTRYPOINT [ "/app/bin" ]

