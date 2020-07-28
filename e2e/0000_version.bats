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

@test "server version" {
    run curl --silent ${MYROOMIES_SERVER_URL}/version
    assert_success
    assert_output "0.0.1"
}
