FROM golang:alpine as gobuild
  RUN apk update
  RUN apk add build-base gcc bash
  RUN rm -rf /var/cache/apk/*

  RUN go install github.com/amonks/run/cmd/run@latest
  WORKDIR /app
  COPY . .
  RUN run build

FROM alpine
  RUN apk update
  RUN apk add ca-certificates iptables ip6tables bash
  RUN rm -rf /var/cache/apk/*

  WORKDIR /app
  COPY . .
  COPY --from=gobuild /app/serve /app/serve

  CMD ["./serve"]
