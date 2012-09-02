The US Visa RESTful Web Service in Go
=====================================

* `make server` - start a multi-threaded server
* `make print` - print out the current table (re-loads it if necessary)
* `make thrasher` - load the original table and verifies the responses from the server (`make server`) where each request is processed concurrently
