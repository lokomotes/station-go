#!/bin/sh

# THIS SCRIPT IS INTENDED TO BE RUN IN A CONTAINER.
# DO NOT RUN ON THE HOST

#
# USER DEFINE VARIABLES
#
pkg_station='github.com/lokomotes/station-go'

#
# DO NOT TOUCH
#
go_path=$(go env GOPATH)

deps=$(go list -f '{{.Imports}}' $pkg_station | tr "[" " " | tr "]" " " | xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}')
deps="$deps $pkg_station"

echo $deps | while read -r element; do
    go install $element
done

mv /go/bin/station-go /go/bin/main