#!/usr/bin/env bats

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'common'

setup myroomies_common_setup

@test "https server start - error missing key file" {
    run ${SERVER_BIN} --cert-file '/tmp/fake-cert'
    assert_failure
    assert_line --index 0 'Error: A path to a certificate was given but no path to a private key.'
}

@test "https server start - error missing certificate file" {
    run ${SERVER_BIN} --key-file '/tmp/fake-key'
    assert_failure
    assert_line --index 0 'Error: A path to a private key was given but no path to a certificate.'
}

@test "https server start - OK" {
    cert_path="${BATS_TMPDIR}/myroomies.cert"
    key_path="${BATS_TMPDIR}/myroomies.key"
    myroomies_generate_keys ${cert_path} ${key_path}
    server_pid=$(myroomies_start_server "" "secret" "8443" "${cert_path}" "${key_path}")
    run curl --silent \
        -u root:secret \
        --insecure \
        https://localhost:8443/users
    assert_success
    [[ "$(echo $output | jq '.[] | length')" == '1' ]]
    [[ "$(echo $output | jq '.[0].Login')" == '"root"' ]]
    myroomies_stop_server ${server_pid}
}
