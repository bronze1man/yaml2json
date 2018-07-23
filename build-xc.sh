#!/bin/sh

APPLICATION_NAME="yaml2json"

SUPPORTED_PLATFORMS=$(cat <<'EOF'
darwin/386
darwin/amd64
freebsd/386
freebsd/amd64
freebsd/arm
linux/386
linux/amd64
linux/arm
netbsd/386
netbsd/amd64
netbsd/arm
openbsd/386
openbsd/amd64
plan9/386
windows/386
windows/amd64
EOF
)

mkdir -p ./builds

for platform in $SUPPORTED_PLATFORMS; do
    export GOOS=`echo $platform | cut -d '/' -f 1`
    export GOARCH=`echo $platform | cut -d '/' -f 2`
    output_folder="./builds/${GOOS}_${GOARCH}"

    # deal with windows filename
    bin_ext=""
    test "$GOOS" = "windows" && bin_ext=".exe"
    
    mkdir -p "$output_folder"

    echo "Building for $GOOS/$GOARCH"
    go build -o "$output_folder/${APPLICATION_NAME}${bin_ext}"
done