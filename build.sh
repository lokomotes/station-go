#!/usr/bin/env bash

#
# USER DEFINE VARIABLES
#
imageRef='lokomotes/station-go:latest'
station_pkg='github.com/lokomotes/station-go'

#
# DO NOT TOUCH
#
go_path=$(go env GOPATH)

if [ ! -f  "./Dockerfile" ]; then
    echo "Dockerfile not provided"
    exit 1
fi

cp -R ./app $go_path/src

tmp_path=$(mktemp -d)

deps=$(go list -f '{{.Deps}}' | tr "[" " " | tr "]" " " | xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}')
deps="${deps}
${station_pkg}"

while read -r element; do
    mkdir -p $tmp_path/src/$element/
    rsync -a $go_path/src/$element/ $tmp_path/src/$element/
    echo $element" done"
done <<< "$deps"

cp ./Dockerfile $tmp_path/Dockerfile
cp ./install.sh $tmp_path/install.sh
cp ./create.sh $tmp_path/create.sh

docker build \
    -t $imageRef \
    $tmp_path

rm -rf $tmp_path