### url-shortener

Default **ENV** variables: 

`SERVICE_ADDR=8080`   - url-shortener port  
`PPROF_ADDR=5050`     - pprof port [optional]  
`STORE_ADDR=6379`     - redis port [in progress]  
`DNS_NAME=localhost`  - if running on Azure instead of `localhost` has to be `FQDN`, example in the Azure part  

#### \# ENDPOINTS

**pprof**:  
```
/debug/pprof/
/debug/pprof/trace
/debug/pprof/goroutine
```  

**url-shortener**:
```
# PREFIX
/us

# GET
/home
/health
/healthz

# POST
/links
```

#### \# HELP

```
$ make help
```


#### \# TEST

```
$ make go-run
```

#### \# DOCKER

```
$ make docker-build

$ make docker-run

# HOME

$ curl localhost:8080/us/home
url-shortener
Hostname: 924d217ffca0; Version: 0.0.1

# HEALTH

$ curl localhost:8080/us/healthz
OK

# redis [in progress]
$ curl localhost:8080/us/health | jq
{
  "state": "OK",
  "urlerrormessages": null,
  "redis": "NOK",
  "rediserrormessages": [
    "HealthError: Get http://localhost:6379: dial tcp 127.0.0.1:6379: connect: connection refused"
  ]
}

# GENERATE

$ curl -X POST \
-H "Content-Type: application/json" \
-d '{"longUrl":"https://google.com"}' \
localhost:8080/us/links | jq
{
  "id": "43600b9c",
  "longUrl": "https://google.com",
  "shortUrl": "http://localhost:8080/43600b9c",
  "createdAt": "2020-01-09 20:34:2419"
}

$ curl http://localhost:8080/43600b9c
<a href="https://google.com">Moved Permanently</a>.

# REMOVE

$ make docker-stop
```

#### \# AZURE

[![Deploy to Azure](http://azuredeploy.net/deploybutton.png)](https://azuredeploy.net/)  


#### Container Instance

```
$ DNS_NAME_LABEL=urlshort-demo-$RANDOM
$ LOCATION=westeurope
$ RGNAME=<enterYourRG>

$ az container create \
  --resource-group $RGNAME \
  --name urlshortener \
  --image michalsw/url-shortener \
  --restart-policy Always \
  --ports 80 5050 6379 \
  --dns-name-label $DNS_NAME_LABEL \
  --location $LOCATION \
  --environment-variables \
    SERVICE_ADDR=80 \
    PPROF_ADDR=5050 \
    STORE_ADDR=6379 \
    DNS_NAME=$DNS_NAME_LABEL.$LOCATION.azurecontainer.io

# GENERATE

$ curl -X POST \
-H "Content-Type: application/json" \
-d '{"longUrl":"https://portal.azure.com/"}' \
http://urlshort-demo-12060.westeurope.azurecontainer.io/us/links | jq

{
  "id": "11624ba8",
  "longUrl": "https://portal.azure.com/",
  "shortUrl": "http://urlshort-demo-12060.westeurope.azurecontainer.io:80/us/11624ba8",
  "createdAt": "2020-03-16 11:22:22316"
}
```

#### API Management

OpenAPI specification available [here](./docs/swagger.json) . Once your API Management service is up and running:
- edit specification adding valid **host** (URL from step above)
- add a new OpenAPI - upload specification (don't forget about **API URL suffix**)

```
$ API_URL_SUFFIX=urlapi
$ SUBSCRIPTION_KEY=<your_key>
$ API_GATEWAY=<your_api_gateway_name>

# GENERATE

$ curl -X POST \
-H "Content-Type: application/json" \
-d '{"longUrl":"https://portal.azure.com/"}' \
-H "Ocp-Apim-Subscription-Key: $SUBSCRIPTION_KEY" \
http://$API_GATEWAY.azure-api.net/$API_URL_SUFFIX/us/links

{
  "id": "c0f025ca",
  "longUrl": "https://portal.azure.com/",
  "shortUrl": "http://urlshort-demo-863.westeurope.azurecontainer.io:80/c0f025ca",
  "createdAt": "2020-01-21 15:05:45121"
}
```
