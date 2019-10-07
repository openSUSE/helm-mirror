FROM opensuse/leap:15.1

LABEL Maintainer="SUSE Containers Team <containers@suse.com>"

RUN zypper -n in \
		git \
		go1.12 \
		golang-github-cpuguy83-go-md2man \
		make \
		tar \
		gzip

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH
RUN go get -u golang.org/x/lint/golint && \
	go get -u github.com/vbatts/git-validation && type git-validation

VOLUME ["/go/src/github.com/openSUSE/helm-mirror"]
WORKDIR /go/src/github.com/openSUSE/helm-mirror
COPY . /go/src/github.com/openSUSE/helm-mirror