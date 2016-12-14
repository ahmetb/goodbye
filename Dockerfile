FROM golang:1.7-onbuild
VOLUME /etc/goodbye/config.json

# override the entrypoint as go-wrapper logs some unstructured text
ENTRYPOINT /go/bin/app
