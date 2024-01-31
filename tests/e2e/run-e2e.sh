for f in ./tests/e2e/*
do
	if test -d "$f";then
		bash "$f"/test.sh
	fi
done
