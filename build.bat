rm -rf builds
goxc -d builds xc
bash -c "mv -f builds/unknown/* builds"