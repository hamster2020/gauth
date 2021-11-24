FROM golang:1.17.2-alpine3.14 AS builder

RUN apk update
RUN apk add --no-cache make g++

COPY . /go/src/github.com/hamster2020/gauth
WORKDIR /go/src/github.com/hamster2020/gauth

ARG gauth_email_verifier_token
ENV GAUTH_EMAIL_VERIFIER_TOKEN $gauth_email_verifier_token

ARG gauth_pwned_passwords_url
ENV GAUTH_PWNED_PASSWORDS_URL $gauth_pwned_passwords_url

RUN make check-docker

RUN go build ./cmd/gauth
RUN go build ./cmd/gauthctl

FROM alpine

ARG gauth_email_verifier_token
ENV GAUTH_EMAIL_VERIFIER_TOKEN $gauth_email_verifier_token

ARG gauth_pwned_passwords_url
ENV GAUTH_PWNED_PASSWORDS_URL $gauth_pwned_passwords_url

COPY --from=builder /go/src/github.com/hamster2020/gauth/gauth /gauth/gauth
COPY --from=builder /go/src/github.com/hamster2020/gauth/gauthctl /gauth/gauthctl
COPY --from=builder /go/src/github.com/hamster2020/gauth/bin/* /usr/bin/

RUN chmod +x /usr/bin/entry

RUN apk add --no-cache shadow
RUN usermod -d /tmp/nobody nobody

USER nobody
WORKDIR /gauth
EXPOSE 3000
ENTRYPOINT ["/usr/bin/entry"]
CMD /gauth/gauth
