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

@test "expense management - CRUD basic cases" {
    # Create another user
    run curl --silent \
        -u ${DEFAULT_HTTP_AUTHORIZATION} \
        -i -X POST \
        -d "$(myroomies_user_json 'Stephanie' 'Stacul' 'false' 'stef' 'ribo')" \
        ${MYROOMIES_SERVER_URL}/users
    assert_success
    assert_line --index 0 --partial 'HTTP/1.1 201 Created'

    # Create a few expenses
    run curl --silent \
        -u ${DEFAULT_HTTP_AUTHORIZATION} \
        -i -X POST \
        -d "$(myroomies_expense_json '12.5' 'Shop 1' '2020-06-19T00:00:00.000Z' 'Expense 1')" \
        ${MYROOMIES_SERVER_URL}/expenses
    assert_success
    assert_line --index 0 --partial 'HTTP/1.1 201 Created'
    assert_line --regexp '^Location: /expenses/[0-9]+'

    run curl --silent \
        -u stef:ribo \
        -i -X POST \
        -d "$(myroomies_expense_json '50' 'Shop 2' '2020-07-01T00:00:00.000Z' 'Expense 2')" \
        ${MYROOMIES_SERVER_URL}/expenses
    assert_success
    assert_line --index 0 --partial 'HTTP/1.1 201 Created'
    assert_line --regexp '^Location: /expenses/[0-9]+'
    user_expense_uri=$(echo ${lines[1]} | grep 'Location' | cut -d ' ' -f 2 | tr -d '\r')

    run curl --silent \
        -u ${DEFAULT_HTTP_AUTHORIZATION} \
        -i -X POST \
        -d "$(myroomies_expense_json '470' 'Shop 3' '2020-06-30T00:00:00.000Z' 'Expense 3')" \
        ${MYROOMIES_SERVER_URL}/expenses
    assert_success
    assert_line --index 0 --partial 'HTTP/1.1 201 Created'
    assert_line --index 1 --regexp '^Location: /expenses/[0-9]+'

    # Getting a single expense
    expense_uri=$(echo ${lines[1]} | grep 'Location' | cut -d ' ' -f 2 | tr -d '\r')
    run curl --silent \
        -u ${DEFAULT_HTTP_AUTHORIZATION} \
        -i -X GET \
        ${MYROOMIES_SERVER_URL}${expense_uri}
    assert_success
    assert_line --index 0 --partial 'HTTP/1.1 200 OK'

    # Removing this last one
    run curl --silent \
        -u ${DEFAULT_HTTP_AUTHORIZATION} \
        -i -X DELETE \
        ${MYROOMIES_SERVER_URL}${expense_uri}
    assert_success
    assert_line --index 0 --partial 'HTTP/1.1 200 OK'

    run curl --silent \
        -u ${DEFAULT_HTTP_AUTHORIZATION} \
        -X GET \
        ${MYROOMIES_SERVER_URL}/expenses
    assert_success
    refute_output ''
    json_response=$output
    number_of_expenses=$(echo ${json_response} | jq '. | length')
    run echo ${number_of_expenses}
    assert_output '2'

    # Non-admin user deleting her own expense
    run curl --silent \
        -u stef:ribo \
        -i -X DELETE \
        ${MYROOMIES_SERVER_URL}${user_expense_uri}
    assert_success
    assert_line --index 0 --partial 'HTTP/1.1 200 OK'

    run curl --silent \
        -u ${DEFAULT_HTTP_AUTHORIZATION} \
        -X GET \
        ${MYROOMIES_SERVER_URL}/expenses
    assert_success
    refute_output ''
    json_response=$output
    number_of_expenses=$(echo ${json_response} | jq '. | length')
    run echo ${number_of_expenses}
    assert_output '1'
}

@test "expense management - CRUD error cases" {
}
