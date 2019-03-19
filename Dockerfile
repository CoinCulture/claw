FROM golang:1.12

RUN apt-get update && apt-get upgrade -y && apt-get install -y pandoc texlive-xetex

ADD https://github.com/golang/dep/releases/download/v0.5.1/dep-linux-amd64 /usr/bin/dep

RUN chmod +x /usr/bin/dep

WORKDIR $GOPATH/src/github.com/CoinCulture/claw

COPY Gopkg.toml Gopkg.lock ./

RUN dep ensure --vendor-only

COPY . ./

RUN go install
