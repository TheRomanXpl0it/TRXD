#!/bin/bash

set -e

files=(
	usermode_hot_switch.py
	register_team_race.py
	submit_race.py
	instance_types.py
	instance_edge_cases.py
	instance_create_race.py
	instance_create_race_team.py
	instance_lifetimes.py
	discord_webhook.py
)

nginx_proxy_files=(
	instance_types.py
	instance_edge_cases.py
	instance_lifetimes.py
)

for file in "${files[@]}"; do
	echo "Running Test: $file"

	cd ../../backend/
	if [[ $file == *"lifetime"* ]]; then
		echo "Using short reclaim interval for lifetime test"
		RECLAIM_INSTANCE_INTERVAL=1 ALLOW_REGISTER=true ./trxd -test-data-WARNING-DO-NOT-USE-IN-PRODUCTION
	else
		ALLOW_REGISTER=true ./trxd -test-data-WARNING-DO-NOT-USE-IN-PRODUCTION
	fi
	./trxd &
	PID=$!

	cd -
	python3 $file

	cd -
	kill $PID
	wait $PID || true
	cd -
done

docker compose -f ../../backend/compose.yml down traefik
docker compose -f ../../backend/compose.yml up -d nginx

for file in "${nginx_proxy_files[@]}"; do
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

# docker compose -f ../../backend/compose.yml down nginx
# docker compose -f ../../backend/compose.yml up -d traefik
