# apiWebserver

## Setup

Install mysql if necessary, on OSX I used brew

`brew install mysql`

Note if you are using brew if you turn your mac on and off you need to run 

`mysql.server restart`

and then the db can be accessed via

`mysql  -u root`

 


GET request sample

```
curl -H "Content-Type: application/json" -g  http://localhost:4000/v1/
```

POST request sample


```
curl -H "Content-Type: application/json" -d '{"name":"jimmy the greek", "age":25}'    http://localhost:4000/v1/
```

curl -H "Content-Type: application/json" -d '{"accountId":9,"details":"buying lots of products AGAIN","amount":201,"date":"2015-01-19T00:00:00Z","updated":0,"created":0}' http://localhost:8080/transactions