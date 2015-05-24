## golang 1.5 image.
# 1.5 needs 1.4 bootstrap
FROM            golang:1.4.2

# Clone go and checkout to pined 1.5 revision (waiting for stable)
RUN             git clone https://github.com/golang/go /goroot && cd /goroot && git checkout 8017ace496f5a21bcd55377e250e325f8ba11d45
# Set the 1.4 bootstrap goroot
ENV             GOROOT_BOOTSTRAP /usr/src/go
# Compile 1.5
RUN             cd /goroot/src && ./make.bash
# Set Envs & Path
ENV             GOROOT /goroot
ENV             GOPATH /gopath
ENV             GOBIN  $GOROOT/bin
ENV             PATH   $GOBIN:$PATH
## End og 1.5 base image.

RUN             go get github.com/tools/godep
MAINTAINER      Guillaume J. Charmes <guillaume@charmes.net

ENV             APP_DIR $GOPATH/src/gotest
WORKDIR         $APP_DIR

ENTRYPOINT      ["/gotest"]
EXPOSE          9090

ADD             .       $APP_DIR
RUN             cd $APP_DIR && godep go build -o /gotest -tags netgo -ldflags '-w -s -linkmode external -extldflags -static'
