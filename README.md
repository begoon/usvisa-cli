The US Visa RESTful Web Service in Go
=====================================

* `make server` - start the multi-threaded web server
* `make print` - print out the current table (re-loads it if necessary)
* `make thrasher` - load the original table and verifys the responses from the server (`make server`) where each request is processed concurrently
