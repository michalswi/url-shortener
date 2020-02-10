### url-shortener

Default **ENV** variables: 

`SERVICE_ADDR=8080`   - url-shortener port  
`PPROF_ADDR=5050`     - pprof port [optional]  
`STORE_ADDR=6379`     - redis port [in progress]  
`DNS_NAME=localhost`  - if running on Azure instead of `localhost` has to be `FQDN`, example in the Azure part  

**pprof** endpoints:  
```
/debug/pprof/
/debug/pprof/trace
/debug/pprof/goroutine
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

$ curl localhost:8080
url-shortener
Hostname: 924d217ffca0; Version: 0.0.1


# redis [in progress]

$ curl localhost:8080/health | jq
{
  "state": "OK",
  "urlerrormessages": null,
  "redis": "NOK",
  "rediserrormessages": [
    "HealthError: Get http://localhost:6379: dial tcp 127.0.0.1:6379: connect: connection refused"
  ]
}


# check

$ curl -X POST \
-H "Content-Type: application/json" \
-d '{"longUrl":"https://google.com"}' \
localhost:8080/links | jq
{
  "id": "43600b9c",
  "longUrl": "https://google.com",
  "shortUrl": "http://localhost:8080/43600b9c",
  "createdAt": "2020-01-09 20:34:2419"
}

$ curl -X POST \
-H "Content-Type: application/json" \
-d '{"longUrl":"https://amazon.com"}' \
localhost:8080/links | jq
{
  "id": "614a1913",
  "longUrl": "https://amazon.com",
  "shortUrl": "http://localhost:8080/614a1913",
  "createdAt": "2020-01-09 20:34:3519"
}

$ curl http://localhost:8080/43600b9c
<a href="https://google.com">Moved Permanently</a>.

$ curl http://localhost:8080/614a1913
<a href="https://amazon.com">Moved Permanently</a>.


# clear

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

# HOME

$ curl urlshort-demo-863.westeurope.azurecontainer.io


# HEALTH

$ curl urlshort-demo-863.westeurope.azurecontainer.io/health | jq
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
-d '{"longUrl":"https://portal.azure.com/"}' \
http://urlshort-demo-863.westeurope.azurecontainer.io/links | jq
{
  "id": "08104971",
  "longUrl": "https://portal.azure.com/",
  "shortUrl": "http://urlshort-demo-863.westeurope.azurecontainer.io:80/08104971",
  "createdAt": "2020-01-16 19:50:33116"
}
```

#### API Management

OpenAPI specification available [here](./docs/swagger.json) . Once uploaded to APIs:

```
$ API_URL_SUFFIX=urlapi
$ SUBSCRIPTION_KEY=<your_key>
$ API_GATEWAY=<your_api_gateway_name>

# HOME

$ curl -H "Ocp-Apim-Subscription-Key: $SUBSCRIPTION_KEY" \
http://$API_GATEWAY.azure-api.net/$API_URL_SUFFIX/


# HEALTH

$ curl -H "Ocp-Apim-Subscription-Key: $SUBSCRIPTION_KEY" \
http://$API_GATEWAY.azure-api.net/$API_URL_SUFFIX/health


# GENERATE

$ curl -X POST \
-H "Content-Type: application/json" \
-d '{"longUrl":"https://portal.azure.com/"}' \
-H "Ocp-Apim-Subscription-Key: $SUBSCRIPTION_KEY" \
http://$API_GATEWAY.azure-api.net/$API_URL_SUFFIX/links
{
  "id": "c0f025ca",
  "longUrl": "https://portal.azure.com/",
  "shortUrl": "http://urlshort-demo-863.westeurope.azurecontainer.io:80/c0f025ca",
  "createdAt": "2020-01-21 15:05:45121"
}
```
