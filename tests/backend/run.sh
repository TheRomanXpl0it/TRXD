#!/bin/bash

set -e

files=(
	register_team_race.py
	submit_race.py
	instance_types.py
	instance_edge_cases.py
	instance_create_race.py
	instance_create_race_team.py
	instance_lifetimes.py
	discord_webhook.py
)

for file in "${files[@]}"; do
	echo "Running Test: $file"

	cd ../../backend/
	if [[ $file == *"lifetime"* ]]; then
		echo "Using short reclaim interval for lifetime test"
		RECLAIM_INSTANCE_INTERVAL=1 ./trxd -test-data-WARNING-DO-NOT-USE-IN-PRODUCTION
	else
		./trxd -test-data-WARNING-DO-NOT-USE-IN-PRODUCTION
	fi
	./trxd -t
	./trxd &
	PID=$!

	cd -
	python3 $file

	cd -
	kill $PID
	wait $PID || true
	cd -
done
