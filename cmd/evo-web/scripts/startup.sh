#!/bin/sh

main () {
    local target=$1
    local addr=$2

    sed -i -- "s/{addr}/$addr/g" $target
}

main $1 $2

exec "nginx" "-g" "daemon off;"
