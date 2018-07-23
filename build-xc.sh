#!/bin/sh

APP_NAME=${APP_NAME:-app}

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
    output="./builds/${APP_NAME}-${GOOS}-${GOARCH}"

    # deal with windows filename
    bin_ext=""
    test "$GOOS" = "windows" && bin_ext=".exe"
    
    echo "Building for $GOOS/$GOARCH"
    go build -ldflags "-X \"main.VersionInfo=${APP_VERSION} $(date "+(%b %Y)")\"" -o "${output}${bin_ext}"
done