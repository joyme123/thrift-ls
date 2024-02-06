#!/bin/bash

echo "run e2e test case: enums"

IFS="."
for f in ./tests/e2e/enums/*.expect
do
	if test -f "$f"; then
		echo "======================================================"
		substr=${f##*/}
		read -ra options <<<"$substr"
		indent=${options[0]}
		align=${options[1]}
		field_line_comma=${options[2]}
		echo "indent: ${indent}, align: ${align}, field_line_comma: ${field_line_comma}"
		got=$(./bin/thriftls -format -indent "${indent}" -align "${align}" -fieldLineComma "${field_line_comma}" -f tests/e2e/enums/enums.thrift)
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
