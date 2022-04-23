# Syslog Sink Integration Test

### Start syslog
```bash
docker-compose up -d syslog
```

### Running the integration test
```bash
# one time test
docker-compose up app

# repeated testing, much faster
docker-compose run --entrypoint=/bin/sh app
$ cd tests/syslog
$ go install
$ go run main.go

```
