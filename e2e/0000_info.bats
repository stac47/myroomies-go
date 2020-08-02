#!/usr/bin/env bats

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'common'

setup myroomies_common_setup

function setup_file() {
    MYROOMIES_SERVER_PID="$(myroomies_start_server)"
}

function teardown_file() {
    myroomies_stop_server ${MYROOMIES_SERVER_PID}
}

function teardown() {
    myroomies_reset_server
}

@test "server global info" {
    expected_version="\"$(myroomies_get_version)\""
    run curl --silent ${MYROOMIES_SERVER_URL}/info
    assert_success
    response="${output}"
    run echo $(echo "${response}" | jq '.Version')
    assert_success
    assert_output ${expected_version}
    run echo $(echo "${response}" | jq '.Name')
    assert_success
    assert_output '"MyRoomies"'
    run echo $(echo "${response}" | jq '.Licence')
    assert_success
    assert_output '"GPL Version 3"'
    run echo $(echo "${response}" | jq '.Creator')
    assert_success
    assert_output '"stac47 - https://github.com/stac47"'
}
