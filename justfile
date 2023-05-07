target := "wint"

export CGO_ENABLED := "0"

default: build

build:
    go build

clean:
    rm -f {{target}}

install destdir bindir="/usr/bin":
    install -D -m 0755 -t "{{destdir}}{{bindir}}" {{target}}
