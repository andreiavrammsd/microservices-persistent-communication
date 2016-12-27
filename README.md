# WORK IN PROGRESS: Microservices Persistent Communication

## This is just an exercise, not considered totally reliable

### Persisting calls to microservices ###

Sometimes you really want to be sure the calls made to your services are really successful. If a service is down, important requests could be lost forever.
Instead of calling your services directly, you can proxy the calls through this tool, which will forward the calls to the actual services.
If a call is rejected because of service down, it will be retried until successful.

#### Setup

* Install Git
* Install Docker
* Install Docker compose
* Clone this repository
* Start: docker-compose up -d
