#!/bin/bash -ex
for i in examples/*; do
    dir=${i#"examples/"}
    go run ${i}/main.go > ${dir}.gcode
done
