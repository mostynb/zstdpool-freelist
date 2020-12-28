#!/bin/bash

rv=0

for i in $(find . -type f -name '*.go')
do
	output=$(gofmt -s -l "$i")
	if [ -n "$output" ]
	then
		echo "Error: please run gofmt -w -s \"$i\""
		rv=1
	fi
done

exit $rv
