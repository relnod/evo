#!/bin/sh

main () {
    local addr=$1
    local target=$2

    echo $1
    echo $2
    sed -i -- "s/{addr}/$addr/g" $target
}

main $1 $2
