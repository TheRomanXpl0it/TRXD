#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1}"
ADMIN_EMAIL="${ADMIN_EMAIL:-dsaodijsajio@gmail.com}"
ADMIN_PASSWORD="${ADMIN_PASSWORD:-testpass}"
CHALL_COUNT="${CHALL_COUNT:-15}"
PLAYER_COUNT="${PLAYER_COUNT:-8}"
PREFIX="${PREFIX:-$(date +%y%m%d%H%M%S)}"

CATEGORY="bulk-${PREFIX}"
WORKDIR="/tmp/trxd-seed-${PREFIX}"
mkdir -p "${WORKDIR}"

HTTP_STATUS=""
HTTP_BODY=""

log() {
	printf '[seed] %s\n' "$*"
}

fail() {
	printf '[seed][error] %s\n' "$*" >&2
	exit 1
}

csrf_token() {
	local jar="$1"
	awk '$6=="csrf_" { print $7 }' "${jar}"
}

session_init() {
	local jar="$1"
	: > "${jar}"
	curl -sS -b "${jar}" -c "${jar}" "${BASE_URL}/api/info" >/dev/null
}

api_call() {
	local method="$1"
	local jar="$2"
	local path="$3"
	local payload="${4-}"
	local body_file csrf

	body_file="$(mktemp)"
	csrf="$(csrf_token "${jar}")"
	if [[ -n "${payload}" ]]; then
		HTTP_STATUS="$(curl -sS -X "${method}" -b "${jar}" -c "${jar}" \
			-H "X-CSRF-Token: ${csrf}" \
			-H 'Content-Type: application/json' \
			-d "${payload}" \
			-o "${body_file}" \
			-w '%{http_code}' \
			"${BASE_URL}/api${path}")"
	else
		HTTP_STATUS="$(curl -sS -X "${method}" -b "${jar}" -c "${jar}" \
			-H "X-CSRF-Token: ${csrf}" \
			-o "${body_file}" \
			-w '%{http_code}' \
			"${BASE_URL}/api${path}")"
	fi
	HTTP_BODY="$(cat "${body_file}")"
	rm -f "${body_file}"
}

expect_status() {
	local expected="$1"
	shift
	api_call "$@"
	if [[ "${HTTP_STATUS}" != "${expected}" ]]; then
		fail "request $1 $3 returned ${HTTP_STATUS}: ${HTTP_BODY}"
	fi
}

json_escape() {
	jq -rn --arg value "$1" '$value'
}

login_admin() {
	local jar="$1"
	session_init "${jar}"
	expect_status 200 POST "${jar}" "/login" "$(jq -nc \
		--arg email "${ADMIN_EMAIL}" \
		--arg password "${ADMIN_PASSWORD}" \
		'{email:$email,password:$password}')"
}

create_player() {
	local idx="$1"
	local jar="${WORKDIR}/player-${idx}.cookies"
	local uname="bulk-u$(printf '%02d' "${idx}")-${PREFIX}"
	local email="bulk-u$(printf '%02d' "${idx}")-${PREFIX}@seed.local"
	local password="SeedPass$(printf '%02d' "${idx}")!"
	local team="bulk-t$(printf '%02d' "${idx}")-${PREFIX}"

	session_init "${jar}"
	expect_status 200 POST "${jar}" "/register" "$(jq -nc \
		--arg name "${uname}" \
		--arg email "${email}" \
		--arg password "${password}" \
		'{name:$name,email:$email,password:$password}')"
	expect_status 200 POST "${jar}" "/teams/register" "$(jq -nc \
		--arg name "${team}" \
		--arg password "${password}" \
		'{name:$name,password:$password}')"

	PLAYER_JARS+=("${jar}")
	PLAYER_NAMES+=("${uname}")
	TEAM_NAMES+=("${team}")
}

create_challenge() {
	local admin_jar="$1"
	local idx="$2"
	local cname="bulk-c$(printf '%02d' "${idx}")-${PREFIX}"
	local cflag="TRXD{bulk_${PREFIX}_$(printf '%02d' "${idx}")}"
	local cid

	expect_status 200 POST "${admin_jar}" "/challenges" "$(jq -nc \
		--arg name "${cname}" \
		--arg category "${CATEGORY}" \
		--arg description "Bulk seeded nginx challenge ${idx}" \
		'{name:$name,category:$category,description:$description,type:"Container",max_points:100,score_type:"Static"}')"

	expect_status 200 GET "${admin_jar}" "/challenges"
	cid="$(printf '%s' "${HTTP_BODY}" | jq -r --arg name "${cname}" '.[] | select(.name==$name) | .id')"
	[[ -n "${cid}" && "${cid}" != "null" ]] || fail "could not resolve challenge id for ${cname}"

	expect_status 200 PATCH "${admin_jar}" "/challenges" "$(jq -nc \
		--argjson chall_id "${cid}" \
		--arg host "lvh.me" \
		--arg image "nginx:1.29-alpine" \
		'{
			chall_id:$chall_id,
			hidden:false,
			host:$host,
			port:80,
			conn_type:"HTTP",
			image:$image,
			hash_domain:true,
			lifetime:3600,
			max_memory:64,
			max_cpu:"1"
		}')"

	expect_status 200 POST "${admin_jar}" "/flags" "$(jq -nc \
		--argjson chall_id "${cid}" \
		--arg flag "${cflag}" \
		'{chall_id:$chall_id,flag:$flag,regex:false}')"

	CHALL_IDS+=("${cid}")
	CHALL_NAMES+=("${cname}")
	CHALL_FLAGS+=("${cflag}")
}

submit_flag_checked() {
	local jar="$1"
	local chall_id="$2"
	local flag="$3"
	local expected_status="$4"

	expect_status 200 POST "${jar}" "/submissions" "$(jq -nc \
		--argjson chall_id "${chall_id}" \
		--arg flag "${flag}" \
		'{chall_id:$chall_id,flag:$flag}')"

	local actual_status
	actual_status="$(printf '%s' "${HTTP_BODY}" | jq -r '.status')"
	[[ "${actual_status}" == "${expected_status}" ]] || fail "submission status ${actual_status} != ${expected_status} for challenge ${chall_id}"
}

declare -a PLAYER_JARS=()
declare -a PLAYER_NAMES=()
declare -a TEAM_NAMES=()
declare -a CHALL_IDS=()
declare -a CHALL_NAMES=()
declare -a CHALL_FLAGS=()

admin_jar="${WORKDIR}/admin.cookies"
first_instance_host=""
wrong_count=0
correct_count=0
repeated_count=0
instance_count=0

log "prefix ${PREFIX}"
log "workdir ${WORKDIR}"
log "logging in admin ${ADMIN_EMAIL}"
login_admin "${admin_jar}"

log "creating category ${CATEGORY}"
expect_status 200 POST "${admin_jar}" "/categories" "$(jq -nc --arg name "${CATEGORY}" '{name:$name}')"

for idx in $(seq 1 "${CHALL_COUNT}"); do
	log "creating challenge ${idx}/${CHALL_COUNT}"
	create_challenge "${admin_jar}" "${idx}"
done

for idx in $(seq 1 "${PLAYER_COUNT}"); do
	log "registering player ${idx}/${PLAYER_COUNT}"
	create_player "${idx}"
done

for p_idx in "${!PLAYER_JARS[@]}"; do
	log "creating instances and submissions for ${PLAYER_NAMES[$p_idx]}"
	for c_idx in "${!CHALL_IDS[@]}"; do
		expect_status 200 POST "${PLAYER_JARS[$p_idx]}" "/instances" "$(jq -nc \
			--argjson chall_id "${CHALL_IDS[$c_idx]}" \
			'{chall_id:$chall_id}')"
		instance_host="$(printf '%s' "${HTTP_BODY}" | jq -r '.host')"
		[[ -n "${instance_host}" && "${instance_host}" != "null" ]] || fail "missing instance host for challenge ${CHALL_IDS[$c_idx]}"
		instance_count=$((instance_count + 1))
		if [[ -z "${first_instance_host}" ]]; then
			first_instance_host="${instance_host}"
		fi

		submit_flag_checked "${PLAYER_JARS[$p_idx]}" "${CHALL_IDS[$c_idx]}" "WRONG{${PREFIX}_${p_idx}_${c_idx}}" "Wrong"
		wrong_count=$((wrong_count + 1))

		if ((( (p_idx + c_idx) % 2 ) == 0 )); then
			submit_flag_checked "${PLAYER_JARS[$p_idx]}" "${CHALL_IDS[$c_idx]}" "${CHALL_FLAGS[$c_idx]}" "Correct"
			correct_count=$((correct_count + 1))

			if ((( (p_idx + c_idx) % 4 ) == 0 )); then
				submit_flag_checked "${PLAYER_JARS[$p_idx]}" "${CHALL_IDS[$c_idx]}" "${CHALL_FLAGS[$c_idx]}" "Repeated"
				repeated_count=$((repeated_count + 1))
			fi
		fi
	done
done

[[ -n "${first_instance_host}" ]] || fail "no instance host recorded"

log "spot-checking ${first_instance_host}"
curl -sS -I "http://${first_instance_host}/" >/dev/null

log "verifying counts in postgres"
docker exec trxd-postgres-1 psql -U user -d postgres -c "
SELECT count(*) AS seeded_challenges
FROM challenges
WHERE category = '${CATEGORY}';

SELECT count(*) AS seeded_instances
FROM instances i
JOIN challenges c ON c.id = i.chall_id
WHERE c.category = '${CATEGORY}';

SELECT status, count(*) AS count
FROM submissions s
JOIN challenges c ON c.id = s.chall_id
WHERE c.category = '${CATEGORY}'
GROUP BY status
ORDER BY status;
" | tee "${WORKDIR}/summary.txt"

printf 'PREFIX=%s\n' "${PREFIX}" | tee -a "${WORKDIR}/summary.txt"
printf 'CATEGORY=%s\n' "${CATEGORY}" | tee -a "${WORKDIR}/summary.txt"
printf 'CHALLENGES=%s\n' "${CHALL_COUNT}" | tee -a "${WORKDIR}/summary.txt"
printf 'PLAYERS=%s\n' "${PLAYER_COUNT}" | tee -a "${WORKDIR}/summary.txt"
printf 'INSTANCES=%s\n' "${instance_count}" | tee -a "${WORKDIR}/summary.txt"
printf 'WRONG=%s\n' "${wrong_count}" | tee -a "${WORKDIR}/summary.txt"
printf 'CORRECT=%s\n' "${correct_count}" | tee -a "${WORKDIR}/summary.txt"
printf 'REPEATED=%s\n' "${repeated_count}" | tee -a "${WORKDIR}/summary.txt"
printf 'FIRST_HOST=%s\n' "${first_instance_host}" | tee -a "${WORKDIR}/summary.txt"

log "done"
