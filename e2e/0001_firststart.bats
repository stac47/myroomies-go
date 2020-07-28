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

@test "first server start - normal start for test suite" {
    run curl --silent -u ${DEFAULT_HTTP_AUTHORIZATION} http://localhost:8080/users
    assert_success
    # The 'Password' field of a User is always omited
    assert_output --partial '"Password":""'
    response="$output"

    run echo $(echo "$response" | jq '. | length')
    assert_success
    assert_output '1'
    run echo $(echo "$response" | jq '.[0].Login')
    assert_success
    assert_output '"root"'
}
