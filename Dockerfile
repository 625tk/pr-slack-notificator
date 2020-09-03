FROM golang:1.15.1 as builder

WORKDIR /pr-slack-notificator
COPY . .
RUN CGO_ENABLED=0 go build 

FROM alpine:3.7 
COPY --from=builder /pr-slack-notificator/pr-slack-notificator /
RUN apk add git
CMD ["/pr-slack-notificator"]

