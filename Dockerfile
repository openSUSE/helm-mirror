FROM opensuse/amd64:42.3

LABEL Maintainer="SUSE Containers Team <containers@suse.com>"

RUN zypper -n up
RUN zypper -n in \
		git \
		go \
		golang-github-cpuguy83-go-md2man \
		make \
		tar

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH
RUN go get -u github.com/golang/dep/cmd/dep && \
	go get -u github.com/golang/lint/golint && \
	go get -u github.com/vbatts/git-validation && type git-validation

VOLUME ["/go/src/github.com/openSUSE/helm-mirror"]
WORKDIR /go/src/github.com/openSUSE/helm-mirror
COPY . /go/src/github.com/openSUSE/helm-mirror