#!/bin/bash

echo "run e2e test case: field line comma"

IFS="."
for f in ./tests/e2e/field_line_comma/*.expect
do
	if test -f "$f"; then
		echo "======================================================"
		substr=${f##*/}
		read -ra options <<<"$substr"
		field_line_comma=${options[0]}
		echo "fieldLineComma: ${field_line_comma}"
		got=$(./bin/thriftls -format -fieldLineComma "${field_line_comma}" -f tests/e2e/field_line_comma/fields.thrift)
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
