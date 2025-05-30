# Information for confluent-kafka-go developers

## Development process

1. Use go1.21 (and related tooling) for development on confluent-kafka-go.
2. Make sure to run `gofmt` and `go vet` on your code.
3. While there is no hard-limit, try to keep your line length under 80
   characters.
3. [Test](#testing) your changes and create a PR.


NOTE: Whenever librdkafka error codes are updated make sure to run generate
before building:

```
  $ make -f mk/Makefile generr
  $ go build ./...
```



## Testing

Some of the tests included in this directory, the benchmark and integration tests in particular,
require an existing Kafka cluster and a testconf.json configuration file to
provide tests with bootstrap brokers, topic name, etc.

The format of testconf.json is a JSON object:
```
{
  "Brokers": "<bootstrap-brokers>",
  "Topic": "<test-topic-name>"
}
```

See testconf-example.json for an example and full set of available options.


To run unit-tests:
```
$ go test
```

To run benchmark tests:
```
$ go test -bench .
```

For the code coverage:
```
$ go test -coverprofile=coverage.out -bench=.
$ go tool cover -func=coverage.out
```


## Build tags

Different build types are supported through Go build tags (`-tags ..`),
these tags should be specified on the **application** build/get/install command.

 * By default the bundled platform-specific static build of librdkafka will
   be used. This works out of the box on Mac OSX and glibc-based Linux distros,
   such as Ubuntu and CentOS.
 * `-tags musl` - must be specified when building on/for musl-based Linux
   distros, such as Alpine. Will use the bundled static musl build of
   librdkafka.
 * `-tags dynamic` - link librdkafka dynamically. A shared librdkafka library
   must be installed manually through other means (apt-get, yum, build from
   source, etc).



## Release process

For each release candidate and final release, perform the following steps:

### Review the CHANGELOG

### Update bundle to latest librdkafka

See instructions in [kafka/librdkafka_vendor/README.md](librdkafka_vendor/README.md).


### Update librdkafka version requirement

Update the minimum required librdkafka version in `kafka/00version.go`
and `README.md` and the version in `examples/go.mod` and `mk/doc-gen.py`.

### Update error codes

Error codes can be automatically generated from the current librdkafka version.


Update generated error codes:

    $ make -f mk/Makefile generr
    # Verify by building


## Generating HTML documentation

To generate one-page HTML documentation run the mk/doc-gen.py script from the
top-level directory. This script requires the beautifulsoup4 Python package.

```
$ source .../your/virtualenv/bin/activate
$ pip install beautifulsoup4
...
$ make -f mk/Makefile docs
```


### Rebuild everything

    $ go clean -i ./...
    $ go build ./...


### Run full test suite

Set up a test cluster using whatever mechanism you typically use
(docker, trivup, ccloud, ..).

Make sure to update `kafka/testconf.json` as needed (broker list, $BROKERS)

Run test suite:

    $ go test ./...


### Verify examples

Manually verify that the examples/ applications work.

Also make sure the examples in README.md work.

### Commit any changes

Make sure to push to github before creating the tag to have CI tests pass.


### Create and push tag

    $ git tag v1.3.0
    $ git push --dry-run origin v1.3.0
    # Remove --dry-run and re-execute if it looks ok.


### Create release notes page on github

### Update version in Confluent docs

Put the new version in settings.sh of these two repos

https://github.com/confluentinc/docs

https://github.com/confluentinc/docs-platform

### Don't forget tweeting it!
