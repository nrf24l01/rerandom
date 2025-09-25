# Rerandom

## How to run main ngixn
```sh
docker run --rm --network host \
-v "$PWD/nginx/nginx.conf:/etc/nginx/nginx.conf:ro" \
-v "$PWD/ssl/server.crt:/etc/ssl/certs/server.crt:ro" \
-v "$PWD/ssl/server.key:/etc/ssl/private/server.key:ro" \
-v "$PWD/nginx:/var/www/html:ro" \
nginx:alpine
```

## Keygen
```sh
openssl req -x509 -newkey rsa:2048 -nodes \                    
-keyout random-org.key -out random-org.crt -days 365 \
-subj "/C=RU/ST=Moscow/L=Moscow/O=ReRandom/OU=IT/CN=random.org" \
-addext "subjectAltName=DNS:random.org,DNS:www.random.org"
```
or
```sh
openssl req -x509 -newkey rsa:2048 -nodes -keyout random-org.key -out random-org.crt -days 365 -subj "/C=RU/ST=Moscow/L=Moscow/O=ReRandom/OU=IT/CN=random.org" -addext "subjectAltName=DNS:random.org,DNS:www.random.org"
```