FROM golang as build-xc

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

ENV GOPATH=/go

WORKDIR /go/src/app

COPY ["./Gopkg.*", "./build-xc.sh", "./"]

RUN dep ensure --vendor-only \
    && chmod a+x ./build-xc.sh

COPY ["./main.go", "./"]

RUN ./build-xc.sh

RUN ls -al

FROM scratch AS binaries

COPY --from=build-xc /go/src/app/builds /builds

# Now for a clean image after the build
FROM alpine

COPY --from=binaries /builds/linux_amd64/yaml2json /bin/yaml2json

ENTRYPOINT [ "/bin/yaml2json" ]