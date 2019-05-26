# JAProxy
Just Another Proxy it's a http proxy built in Go.
## Building
The project features a Makefile with the following commands:
 - **build** builds the project binary of the running platform
 - **build all** builds all the binaries for the platforms specified inside the Makefile
 - **push** builds the binary and images and then push the docker images to the repository configured in the makefile
 
_The build procedure is done inside a docker image, so naturally you'll need to have Docker installed and running on your machine_
## Features
 - Rate limiting with Redis support.
 - Can be configured through JSON or Consul with values JSON alike :full_moon_with_face:
 - Exposes metrics through Prometheus
## Known issues and follow up fixes :ghost:
 - Code should be cleaned up, unify the limiters behind one method to avoid duplicate code on rate limit handle(writing response, setting limits exceeded headers, etc)
