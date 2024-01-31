#!/bin/bash

echo "run e2e test case: fields"

IFS="."
for f in ./tests/e2e/fields/*.expect
do
	if test -f "$f"; then
		echo "======================================================"
		substr=${f##*/}
		read -ra options <<<"$substr"
		indent=${options[0]}
		align=${options[1]}
		echo "indent: ${indent}, align: ${align}"
		got=$(./bin/thriftls -format -indent "${indent}" -align "${align}" -f tests/e2e/fields/fields.thrift)
		expected=$(cat "$f")
		if [ "$got" ==  "$expected" ];then
			echo "pass"
		else :
			echo "failed"
			printf 'got: \n%s\n' "${got}"
			printf 'expected: \n%s\n' "${expected}"
			diff <(echo "$got") <(echo "$expected")
		fi
		echo "======================================================"
	fi

done
