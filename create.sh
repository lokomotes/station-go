#!/bin/sh

# THIS SCRIPT IS INTENDED TO BE RUN IN A CONTAINER.
# DO NOT RUN ON THE HOST

#
# USER DEFINE VARIABLES
#
pkg_usr="$1"

#
# DO NOT TOUCH
#
go_path=$(go env GOPATH)
pkg_path=$go_path/pkg/$(go env GOHOSTOS)_$(go env GOHOSTARCH)

go get $pkg_usr/...
go install $pkg_usr

mv $pkg_path/$pkg_usr.a $pkg_path/app.a

# rm -rf $go_path/src/app/*

echo "//go:binary-only-package

package app" > $go_path/src/app/app.go

go install github.com/lokomotes/station-go

mv /go/bin/station-go /go/bin/main
