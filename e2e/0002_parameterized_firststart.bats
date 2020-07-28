#!/usr/bin/env bats

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'common'

setup myroomies_common_setup

@test "first server start - error no root password" {
    run ${SERVER_BIN}
    assert_failure
}

@test "first server start - error root login given but no root password" {
    run MYROOMIES_ROOT_LOGIN=stac ${SERVER_BIN}
    assert_failure
}

@test "first server start - only password given" {
    server_pid=$(myroomies_start_server "" "anotherpass" "8081")
    run curl --silent -u root:anotherpass http://localhost:8081/users
    assert_success
    [[ "$(echo $output | jq '.[] | length')" == '1' ]]
    [[ "$(echo $output | jq '.[0].Login')" == '"root"' ]]
    myroomies_stop_server ${server_pid}
}

@test "first server start - login and password given" {
    server_pid=$(myroomies_start_server "stac" "mysecret" "8081")
    run curl --silent -u stac:mysecret http://localhost:8081/users
    assert_success
    [[ "$(echo $output | jq '.[] | length')" == '1' ]]
    [[ "$(echo $output | jq '.[0].Login')" == '"stac"' ]]
    myroomies_stop_server ${server_pid}
}
