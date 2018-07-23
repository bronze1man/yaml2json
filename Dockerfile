FROM golang as build-xc

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

ENV GOPATH=/go

WORKDIR /go/src/app

COPY ["./Gopkg.*", "./build-xc.sh", "./"]

RUN dep ensure --vendor-only \
    && chmod a+x ./build-xc.sh

ARG APP_NAME=yaml2json
ARG APP_VERSION

ENV APP_NAME=${APP_NAME}
ENV APP_VERSION=${APP_VERSION}

COPY ["./main.go", "./"]

RUN ./build-xc.sh

RUN ls -al /go/src/app/builds

FROM scratch AS binaries

COPY --from=build-xc /go/src/app/builds /builds

# Now for a clean image after the build
FROM alpine

ARG APP_NAME=yaml2json
ARG GOOS=linux
ARG GOARCH=amd64

COPY --from=binaries /builds/${APP_NAME}-${GOOS}-${GOARCH} /bin/${APP_NAME}

ENTRYPOINT [ "/bin/yaml2json" ]