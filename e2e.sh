#!/bin/sh

# MARK: - Setup

build_if_needed() {
    # Check if there already is a pre-built executable
    if [ ! -f ./bin/commit ]; then
        # Build the executable, exit if build fails
        make build || exit 1
    fi
}

red() {
    echo "\033[0;31m$1\033[0m"
}

green() {
    echo "\033[0;32m$1\033[0m"
}

yellow() {
    echo "\033[0;33m$1\033[0m"
}

setup_global_git_config() {
    git config --global init.defaultBranch main
    git config --global user.name "GitHub Actions CI Runner"
    git config --global user.email "--"
}

setup_test_repository() { 
    # Create a new directory
    mkdir testdir && \
    cd testdir && \

    # Initialize a new repository
    git init && \
    touch file && git add file && git commit -m "Initial commit"
}

start_test() {
    echo ""
    echo $(yellow "TEST: $1")
}

fail_test() {
    echo $(red "FAIL: $1")
        cd ..
        rm -rf testdir
        exit 1
}

pass_test() {
    echo $(green "PASS: $1")
        cd ..
        rm -rf testdir
}

# MARK: - Test Cases

test_commit_from_current_directory_without_config() {
    TESTNAME="test_commit_from_current_directory_without_config"
    start_test $TESTNAME

    setup_test_repository &&\
    git checkout -b 1-initial-branch && \

    # Create a new file
    echo "Hello, World!" > hello.txt && \

    git add hello.txt && \

    # Commit the file
    ../bin/commit "Hello, git!"

    # Check if the commit was successful
    if [ $? -ne 0 ]; then
        fail_test $TESTNAME
    fi

    # Check if the commit message is correct
    if [ "$(git log -1 --pretty=%B)" != '#1: Hello, git!' ]; then
        fail_test $TESTNAME
    fi

    pass_test $TESTNAME
}

test_use_config_from_current_directory() {
    TESTNAME="test_use_config_from_current_directory"
    start_test $TESTNAME

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
        fail_test $TESTNAME
    fi

    # Check if the commit message is correct
    if [ "$(git log -1 --pretty=%B)" != '[DEV-38] Add a new file' ]; then
        fail_test $TESTNAME
    fi

    pass_test $TESTNAME
}

test_commit_from_subdirectory() {
    TESTNAME="test_commit_from_subdirectory"
    start_test $TESTNAME

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
        fail_test $TESTNAME
    fi

    # Check if the commit message is correct
    if [ "$(git log -1 --pretty=%B)" != '(cfg13) Do something very useful' ]; then
        cd ..
        fail_test $TESTNAME
    fi

    cd ..
    pass_test $TESTNAME
}

test_set_correct_author() {
    TESTNAME="test_set_correct_author"
    start_test $TESTNAME

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
        fail_test $TESTNAME
    fi

    # Check if the commit message is correct
    if [ "$(git log -1 --pretty=%B)" != '#1: Hello, git!' ]; then
        fail_test $TESTNAME
    fi

    # Check if the commit author name is correct
    ACTUAL_AUTHOR_NAME=$(git log -1 --pretty=%an)
    echo "Author name: ${ACTUAL_AUTHOR_NAME}"
    if [ "${ACTUAL_AUTHOR_NAME}" != "${EXPECTED_AUTHOR_NAME}" ]; then
        echo "Incorrect author name: expected ${EXPECTED_AUTHOR_NAME}, got ${ACTUAL_AUTHOR_NAME}"
        fail_test $TESTNAME
    fi

    # Check if the commit author email is correct
    ACTUAL_EMAIL=$(git log -1 --pretty=%ae)
    echo "Author email: ${ACTUAL_EMAIL}"
    if [ "${ACTUAL_EMAIL}" != "${EXPECTED_EMAIL}" ]; then
        echo "Incorrect author email: expected ${EXPECTED_EMAIL}, got ${ACTUAL_EMAIL}"
        fail_test $TESTNAME
    fi

    pass_test $TESTNAME
}

test_use_config_with_empty_regex() {
    TESTNAME="test_use_config_with_empty_regex"
    start_test $TESTNAME

    setup_test_repository &&\
    git checkout -b feature/WIP-88-add-privacy-manifest && \

    # Write a config file
    echo '
    { 
        "issueRegex": ""
    }
    ' > .commit.json && \

    # Create a new file
    echo "Hello, World!" > hello && \

    git add hello && \

    # Commit the file
    echo "Expecting exit with error..." && \
    ../bin/commit "Add missing privacy manifest"

    # Check if the commit was successful
    if [ $? -eq 0 ]; then
        echo "Expected exit with error, but the commit was successful"
        fail_test $TESTNAME
    fi 

    echo "Failed with error as expected!"
    pass_test $TESTNAME

}

# MARK: - Run Tests

build_if_needed
setup_global_git_config

test_commit_from_current_directory_without_config
test_use_config_from_current_directory
test_commit_from_subdirectory
test_set_correct_author
test_use_config_with_empty_regex