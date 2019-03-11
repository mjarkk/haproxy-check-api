# `haproxy-check-api` - A api to check haproxy configs
A small docker container just for checking if a haproxy config doesn't contains errors

## *1* Setup
**Make sure you have installed Docker**
```
$ wget https://github.com/mjarkk/haproxy-check-api/releases/download/0.1/release.zip
$ unzip release.zip
$ docker build --no-cache --tag haproxy-check-api:latest .
```

## *2* Run
```
$ docker run --restart always --name haproxyCheckApi -d -p 8223:8223 haproxy-check-api
```

## *3* Usage
```
$ curl -X POST http://localhost:8223/checkHaProxy -F "file=@./haproxyConfig.cfg" -H "Content-Type: multipart/form-data"
```
If there are any errors it will return status code **400** with the error as response, if the config is oke it will return **OK** with status code **200**
