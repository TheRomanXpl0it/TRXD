#!/bin/bash

SUBNET="172.137.0.0/16"
PROXY_IP="172.137.0.2"

set -e

rule_exists() {
	sudo iptables -C "$@" 2>/dev/null
}

create_rule() {
	if ! rule_exists "$@"; then
		sudo iptables -I "$@"
	fi
}

delete_rule() {
	if rule_exists "$@"; then
		sudo iptables -D "$@"
	fi
}

list_all_rules() {
	sudo iptables -L
}

list_rules() {
	sudo iptables -C "$@" 2>/dev/null || echo "Rule does not exist: $@"
}

rules=(
	"DOCKER-USER -s $SUBNET -d $SUBNET -j DROP"
	"DOCKER-USER -s $PROXY_IP -d $SUBNET -j ACCEPT"
	"DOCKER-USER -s $SUBNET -d $SUBNET -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT"
)

ACTION="$1"
case "$ACTION" in
	list)
		for rule in "${rules[@]}"; do
			list_rules $rule
		done
		;;
	list-all)
		list_all_rules
		;;
	delete)
		for rule in "${rules[@]}"; do
			delete_rule $rule
		done
		docker network rm trxd-shared-internal
		;;
	*)
		for rule in "${rules[@]}"; do
			create_rule $rule
		done
		docker network create --subnet=$SUBNET trxd-shared-internal
		;;
esac
