#!/bin/bash -ex
DIRS=$(cd examples && find . -mindepth 1 -type d -printf '%f ')
for dir in ${DIRS}; do
    go run examples/${dir}/main.go > ${dir}.gcode
done
