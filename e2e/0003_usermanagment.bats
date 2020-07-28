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

@test "user management - CRUD by root" {
    # Create a new user
    run curl --silent \
        -u ${DEFAULT_HTTP_AUTHORIZATION} \
        -i -X POST \
        -d "$(myroomies_user_json "Stephanie" "Drevet" "false" "stef" "ribo")" \
        ${MYROOMIES_SERVER_URL}/users
    assert_success
    assert_line --index 0 --partial 'HTTP/1.1 201 Created'
    assert_line --regexp '^Location: /users/[a-zA-Z]+'

    # Retrieve this user
    run curl --silent \
        -u ${DEFAULT_HTTP_AUTHORIZATION} \
        -X GET \
        ${MYROOMIES_SERVER_URL}/users/stef
    assert_success
    refute_output ''
    response=$output
    run echo $(echo "$response" | jq '.Lastname')
    assert_success
    assert_output '"Drevet"'

    # Changing the password and the lastname
    run curl --silent \
        -u ${DEFAULT_HTTP_AUTHORIZATION} \
        -X PUT \
        -d '{"Lastname":"Stacul", "Password":"juliette"}' \
        ${MYROOMIES_SERVER_URL}/users/stef
    assert_success
    assert_output ''

    # Retrieve the updated user
    run curl --silent \
        -u ${DEFAULT_HTTP_AUTHORIZATION} \
        -X GET \
        ${MYROOMIES_SERVER_URL}/users/stef
    assert_success
    refute_output ''
    response=$output
    run echo $(echo "$response" | jq '.Lastname')
    assert_success
    assert_output '"Stacul"'

    # Retrieve the updated user with his own account
    run curl --silent \
        -u stef:juliette \
        -X GET \
        ${MYROOMIES_SERVER_URL}/users/stef
    assert_success
    refute_output ''
    response=$output
    run echo $(echo "$response" | jq '.Lastname')
    assert_success
    assert_output '"Stacul"'

    # Delete user
    run curl --silent \
        -u ${DEFAULT_HTTP_AUTHORIZATION} \
        -X DELETE \
        ${MYROOMIES_SERVER_URL}/users/stef
    assert_success
    assert_output ''

    # Verify only the root is still here
    run curl --silent -u ${DEFAULT_HTTP_AUTHORIZATION} ${MYROOMIES_SERVER_URL}/users
    assert_success
    response="$output"
    run echo $(echo "$response" | jq '. | length')
    assert_success
    assert_output '1'
    run echo $(echo "$response" | jq '.[0].Login')
    assert_success
    assert_output '"root"'
}

@test "user management - CRUD by non-admin user" {
    # Create a new user
    run curl --silent \
        -u ${DEFAULT_HTTP_AUTHORIZATION} \
        -X POST \
        -d '{"Firstname":"Stephanie", "Lastname":"Drevet", "IsAdmin":false, "Login": "stef",  "Password":"ribo"}' \
        ${MYROOMIES_SERVER_URL}/users
    assert_success
    assert_output ''

    # Retrieve this user
    run curl --silent \
        -u stef:ribo \
        -X GET \
        ${MYROOMIES_SERVER_URL}/users/stef
    assert_success
    refute_output ''
    response=$output
    run echo $(echo "$response" | jq '.Lastname')
    assert_success
    assert_output '"Drevet"'

    # Changing the password and the lastname
    run curl --silent \
        -u stef:ribo \
        -X PUT \
        -d '{"Lastname":"Stacul", "Password":"juliette"}' \
        ${MYROOMIES_SERVER_URL}/users/stef
    assert_success
    assert_output ''

    # Retrieve the updated user with his own account
    run curl --silent \
        -u stef:juliette \
        -X GET \
        ${MYROOMIES_SERVER_URL}/users/stef
    assert_success
    refute_output ''
    response=$output
    run echo $(echo "$response" | jq '.Lastname')
    assert_success
    assert_output '"Stacul"'
}
