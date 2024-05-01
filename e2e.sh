#!/bin/sh

# MARK: - Setup

function build_if_needed() {
    # Check if there already is a pre-built executable
    if [ ! -f ./bin/commit ]; then
        # Build the executable, exit if build fails
        make build || exit 1
    fi
}

function red() {
    echo "\033[0;31m$1\033[0m"
}

function green() {
    echo "\033[0;32m$1\033[0m"
}

function yellow() {
    echo "\033[0;33m$1\033[0m"
}

function setup_test_repository() { 
    # Create a new directory
    mkdir testdir && \
    cd testdir && \

    # Initialize a new repository
    git init && \
    touch file && git add file && git commit -m "Initial commit"
}

function start_test() {
    echo ""
    echo $(yellow "TEST: $1")
}

function fail_test() {
    echo $(red "FAIL: $1")
        cd ..
        rm -rf testdir
        exit 1
}

function pass_test() {
    echo $(green "PASS: $1")
        cd ..
        rm -rf testdir
}

# MARK: - Test Cases

function test_commit_from_current_directory_without_config() {
    start_test $FUNCNAME

    setup_test_repository &&\
    git checkout -b 1-initial-branch && \

    # Create a new file
    echo "Hello, World!" > hello.txt && \

    git add hello.txt && \

    # Commit the file
    ../bin/commit "Hello, git!"

    # Check if the commit was successful
    if [ $? -ne 0 ]; then
        fail_test $FUNCNAME
    fi

    # Check if the commit message is correct
    if [ "$(git log -1 --pretty=%B)" != '#1: Hello, git!' ]; then
        fail_test $FUNCNAME
    fi

    pass_test $FUNCNAME
}

function test_use_config_from_current_directory() {
    start_test $FUNCNAME

    setup_test_repository &&\
    git checkout -b feature/DEV-38-setup-new-module && \

    # Write a config file
    echo '
    { 
        "issueRegex": "DEV-[0-9]+",
        "outputIssuePrefix": "[",
        "outputIssueSuffix": "]",
        "outputStringPrefix": "",
        "outputStringSuffix": " "
    }
    ' > .commit.json && \

    # Create a new file
    echo "Hello, World!" > hello && \

    git add hello && \

    # Commit the file
    ../bin/commit "Add a new file"

    # Check if the commit was successful
    if [ $? -ne 0 ]; then
        fail_test $FUNCNAME
    fi

    # Check if the commit message is correct
    if [ "$(git log -1 --pretty=%B)" != '[DEV-38] Add a new file' ]; then
        fail_test $FUNCNAME
    fi

    pass_test $FUNCNAME
}

function test_commit_from_subdirectory() {
    start_test $FUNCNAME

    setup_test_repository &&\
    git checkout -b prepare-for-cfg13 && \

    # Write a config file
    echo '
    { 
        "issueRegex": "cfg[0-9]+",
        "outputIssuePrefix": "",
        "outputIssueSuffix": "",
        "outputStringPrefix": "(",
        "outputStringSuffix": ") "
    }
    ' > .commit.json && \

    # Create a new file
    echo "Hello, World!" > hello.world && \

    git add hello.world && \

    # Create a new subdirectory
    mkdir testsubdir && cd testsubdir && \

    # Commit the file
    ../../bin/commit "Do something very useful"

    # Check if the commit was successful
    if [ $? -ne 0 ]; then
        cd ..
        fail_test $FUNCNAME
    fi

    # Check if the commit message is correct
    if [ "$(git log -1 --pretty=%B)" != '(cfg13) Do something very useful' ]; then
        cd ..
        fail_test $FUNCNAME
    fi

    cd ..
    pass_test $FUNCNAME
}

function test_set_correct_author() {
    start_test $FUNCNAME

    EXPECTED_AUTHOR_NAME="John Doe"
    EXPECTED_EMAIL="johntheprogrammer@commit.commit"

    setup_test_repository &&\
    git config --local user.name "${EXPECTED_AUTHOR_NAME}" && \
    git config --local user.email "${EXPECTED_EMAIL}" && \
    git checkout -b 1-initial-branch && \

    # Create a new file
    echo "Hello, World!" > hello.txt && \

    git add hello.txt && \

    # Commit the file
    ../bin/commit "Hello, git!"

    # Check if the commit was successful
    if [ $? -ne 0 ]; then
        fail_test $FUNCNAME
    fi

    # Check if the commit message is correct
    if [ "$(git log -1 --pretty=%B)" != '#1: Hello, git!' ]; then
        fail_test $FUNCNAME
    fi

    # Check if the commit author name is correct
    ACTUAL_AUTHOR_NAME=$(git log -1 --pretty=%an)
    echo "Author name: ${ACTUAL_AUTHOR_NAME}"
    if [ "${ACTUAL_AUTHOR_NAME}" != "${EXPECTED_AUTHOR_NAME}" ]; then
        echo "Incorrect author name: expected ${EXPECTED_AUTHOR_NAME}, got ${ACTUAL_AUTHOR_NAME}"
        fail_test $FUNCNAME
    fi

    # Check if the commit author email is correct
    ACTUAL_EMAIL=$(git log -1 --pretty=%ae)
    echo "Author email: ${ACTUAL_EMAIL}"
    if [ "${ACTUAL_EMAIL}" != "${EXPECTED_EMAIL}" ]; then
        echo "Incorrect author email: expected ${EXPECTED_EMAIL}, got ${ACTUAL_EMAIL}"
        fail_test $FUNCNAME
    fi

    pass_test $FUNCNAME
}

# MARK: - Run Tests

build_if_needed

test_commit_from_current_directory_without_config
test_use_config_from_current_directory
test_commit_from_subdirectory
test_set_correct_author