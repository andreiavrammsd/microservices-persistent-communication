# Microservices Persistent Communication

## This is just an exercise, not considered totally reliable

### Persisting calls to microservices

Sometimes you really want to be sure the calls made to your services are really successful. If a service is down, important requests could be lost forever.
Instead of calling your services directly, you can proxy the calls through this tool, which will forward the calls to the actual services.
If a call is rejected because service is down, it will be retried until successful (configurable).

#### Setup

* Install Git
* Install Docker
* Install Docker compose
* Clone this repository
* cp env.dist.yml env.yml
* docker-compose up -d (if RABBITMQ_HIPE_COMPILE option set to 1 - see Configuration -, the service can take up to a few minutes to start)

#### Configuration (env.yml)
```
QUEUE_HOST: "queue" Queue service name as defined in docker-compose.yml
QUEUE_PORT: "5672" Port for queue instance
QUEUE_NAME: "services" Name of queue
QUEUE_NUMBER_OF_CONSUMERS: 3 Number of queue consumers
RETRY_FAILED_AFTER_MILLISECONDS: 5000 Time to retry failed calls after
FILE_LOG_ENABLED: 1 If 1, logs will pe published to file
FAST_PUBLISH: 0 Performance tweaking. If 1, data validation will not be performed when making a request. Request is ignored only if body is empty, any other input is accepted. If no url, the consumers will ignore the calls and will remove them from queue. 
RABBITMQ_DEFAULT_USER: "usernameforrabbitmq" Username for rabbitmq instance
RABBITMQ_DEFAULT_PASS: "passwordforrabbitmq" Password for rabbitmq instance
RABBITMQ_DEFAULT_VHOST: "/" Vhost for rabbitmq instance
RABBITMQ_HIPE_COMPILE: 1 See "hipe_compile" https://www.rabbitmq.com/configure.html
```

#### Request

* Url: 0.0.0.0:8008
* Body: Json string with the following properties:
    * url (string) The service url (accepted protocols: http, https )
    * method (string, optional, default GET) GET, HEAD, POST, PUT, PATCH, DELETE, OPTIONS, CONNECT, TRACE
    * headers (object, optional) Key - value object with headers to send
    * body (string, optional) String to send as body (accepted content type: json, xml)
    * retry (bool, optional, default true) Whether to retry failed calls or not
* Examples
    ```
    {
        "url" : "https://api.myservicedomain.tld/notification",
        "method": "POST",
        "headers": {"X-Auth": "MyAuthKey", "Content-Type": "application/json; charset=UTF-8"},
        "body": "{\"receiver\": \"john.doe@domain.tld\", \"text\": \"Welcome to the machine!\"}"
    }
    ```

    ```
    {
        "url" : "https://api.myservicedomain.tld/ping",
        "retry": false
    }
    ```
* Test
    ```
    curl -X POST -k -d '{
                "url" : "https://api.myservicedomain.tld/notification",
                "method": "POST",
                "headers": {"X-Auth": "MyAuthKey", "Content-Type": "application/json; charset=UTF-8"},
                "body": "{\"receiver\": \"john.doe@domain.tld\", \"text\": \"Welcome to the machine!\"}"
            }' http://0.0.0.0:8008/
    ```

#### Response

* Status codes
    * 200 Only for requests with FAST_PUBLISH set to 1.
    * 202 Request was accepted
    * 400 Invalid body properties
    * 404
    * 422 Invalid body syntax
* Body: Json string with the following properties:
    * error (bool) Whether the publish was successful or not
    * message (string) Success or error message

#### Workflow

* You make a request to this service.
* The request is validated (if FAST_PUBLISH is set to 0).
* Your request is sent to a queue.
* One or many consumers (QUEUE_NUMBER_OF_CONSUMERS) process the queue and perform the calls.
* If a call fails (and retry option was not set to false), it will be requeued after a specified time (RETRY_FAILED_AFTER_MILLISECONDS).

#### Call validation rules

Your services must return an http status code in 2xx class (Success) for a call to be considered successful.

#### Scaling

* You can set up this tool on multiple machines and have a load balancer in front.

#### Performance testing

* https://httpd.apache.org/docs/2.4/programs/ab.html
* Save a body (see examples) in a file (post.json)
* ab -n 1000 -c 100 -v 1 -p post.json -v 0 http://0.0.0.0:8008/

#### Development

* Rebuild and run service: docker-compose up -d --build
* Rebuild and run app: docker-compose exec app go get && docker-compose restart app
* Logs: docker-compose logs -f app
