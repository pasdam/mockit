name: Continuous integration

on: [push]

jobs:

  checks:
    name: Unit test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          submodules: true
      - uses: actions/setup-go@v3
        with:
          go-version: '1.14'
      - run: make go-coverage
      - uses: actions/upload-artifact@v3
        with:
          name: coverage.html
          path: /tmp/coverage.html
