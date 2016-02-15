# golangExampleSimpleApi

This is a very simple example of using an API with Go and shows how easy it is to setup right out of the box. Currently contains no validation or security.

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


