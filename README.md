# JAProxy
Just Another Proxy it's a http proxy built in Go.
## Requirements
You should have dep installed because that's what we use to manage the dependencies, go mod maybe some day.

###### And Go installed, right?
## Building
The project features a Makefile with the following commands:
 - **build** builds the project binary of the running platform
 - **build all** builds all the binaries for the platforms specified inside the Makefile
 - **push** builds the binary and images and then push the docker images to the repository configured in the makefile
 
_The build procedure is done inside a docker image, so naturally you'll need to have Docker installed and running on your machine_
## Features
 - Target balance atm it's only implemented with a Round Robin strategy
 - Rate limiting with Redis support implemented on top of AWS Elasticache
 - Can be configured through JSON or Consul with values JSON alike :full_moon_with_face:
 - Exposes metrics through Prometheus
 - Possibility of implementing middleware functions on request and client responses
## Known issues and follow up fixes :ghost:
 - Code should be cleaned up, unify the limiters behind one method to avoid duplicate code on rate limit handle(writing response, setting limits exceeded headers, etc)
 - Currently we don't feature a circuit breaker so if a target is dead... good luck.
 - Implement a key value Consul configuration instead of parsing a JSON... which is the same as using a JSON file for config, just remotely fetched :sweat_smile:
 - Current redis implementation doesn't support Elasticache cluster node autodiscovery through configuration endpoint.
