#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SERVER_BIN="${SCRIPT_DIR}/../myroomies-server"
CLIENT_BIN="${SCRIPT_DIR}/../myroomies-client"

DEFAULT_ROOT_LOGIN="root"
DEFAULT_ROOT_PASSWORD="defaultpass"
DEFAULT_HTTP_AUTHORIZATION="${DEFAULT_ROOT_LOGIN}:${DEFAULT_ROOT_PASSWORD}"

MYROOMIES_SERVER_PORT=8080
MYROOMIES_SERVER_URL="http://localhost:${MYROOMIES_SERVER_PORT}"

function myroomies_common_setup() {
    if [[ ! -x ${SERVER_BIN} ]]; then
        pushd "$(dirname ${SERVER_BIN})"
        make build-server
        popd
    fi
}

function myroomies_start_server() {
    local root_login=${1:-${DEFAULT_ROOT_LOGIN}}
    local root_password=${2:-${DEFAULT_ROOT_PASSWORD}}
    local port=${3:-"${MYROOMIES_SERVER_PORT}"}
    export MYROOMIES_ROOT_LOGIN="${root_login}"
    export MYROOMIES_ROOT_PASSWORD="${root_password}"
    local mongodb_uri=${MYROOMIES_E2E_TESTS_MONGODB_ADDRESS}
    local server_options="--bind-to :${port}"
    if [[ -n "${mongodb_uri}" ]]; then
        server_options="${server_options} --storage ${mongodb_uri}"
    fi
    ${SERVER_BIN} ${server_options} >>e2e_test_server.log 2>&1 &
    local server_pid=$!
    until $(curl 2>/dev/null --output /dev/null --silent --fail http://localhost:${port}/version); do
        sleep 0.1
    done
    echo ${server_pid}
}

function myroomies_stop_server() {
    local server_pid=${1:?"Missing the server PID"}
    kill ${server_pid}
}

function myroomies_reset_server() {
    curl --silent \
        -u ${DEFAULT_HTTP_AUTHORIZATION} \
        -X POST \
        ${MYROOMIES_SERVER_URL}/reset
}


function myroomies_user_json() {
    local firstname=${1:?"Missing user's firstname"}
    local lastname=${2:?"Missing user's lastname"}
    local is_admin=${3:?"Missing admin right"}
    local login=${4:?"Missing user's login"}
    local password=${5:?"Missing user's password"}
    echo "{\"Firstname\":\"${firstname}\", \"Lastname\":\"${lastname}\", \"IsAdmin\":${is_admin}, \"Login\": \"${login}\", \"Password\":\"${password}\"}"
}

function myroomies_expense_json() {
    local amount=${1:?"Missing expense amount"}
    local recipient=${2:?"Missing expense recipient"}
    local date=${3:?"Missing expense date"}
    local description=${4:?"Missing expense description"}
    echo "{\"Amount\": ${amount}, \"Recipient\":\"${recipient}\", \"Date\":\"${date}\", \"Description\": \"${description}\"}"
}
