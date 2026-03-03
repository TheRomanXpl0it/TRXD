#!/bin/bash

set -e

files=(
	instance_types.py
	instance_edge_cases.py
	instance_lifetimes.py
)

for file in "${files[@]}"; do
	echo "Running Test: $file"

	cd ../../backend/
	if [[ $file == *"lifetime"* ]]; then
		echo "Using short reclaim interval for lifetime test"
		PROXY=nginx RECLAIM_INSTANCE_INTERVAL=1 ALLOW_REGISTER=true ./trxd -test-data-WARNING-DO-NOT-USE-IN-PRODUCTION
	else
		PROXY=nginx ALLOW_REGISTER=true ./trxd -test-data-WARNING-DO-NOT-USE-IN-PRODUCTION
	fi
	./trxd &
	PID=$!

	cd -
	PROXY=nginx python3 $file

	cd -
	kill $PID
	wait $PID || true
	cd -
done

