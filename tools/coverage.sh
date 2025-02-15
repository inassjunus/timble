#!/bin/bash
#
# Code coverage generation

COVERAGE_DIR="${COVERAGE_DIR:-coverage}"

# Create the coverage files directory
mkdir -p "$COVERAGE_DIR";

# Create a coverage file for each package
go test --cover -covermode=count  -coverprofile="${COVERAGE_DIR}"/coverage.cov `go list ./... | grep -v /mocks/` ;

# Display the global code coverage
go tool cover -func="${COVERAGE_DIR}"/coverage.cov ;

# If needed, generate HTML report
if [ "$1" == "html" ]; then
    go tool cover -html="${COVERAGE_DIR}"/coverage.cov -o coverage.html ;
fi

# If needed, generate cobertura report
if [ "$1" == "cobertura" ]; then
    gocover-cobertura < "${COVERAGE_DIR}"/coverage.cov > coverage.xml ;
fi

# Remove the coverage files directory
rm -rf "$COVERAGE_DIR";
