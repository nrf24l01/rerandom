# Rerandom

## How to run main ngixn
```sh
cd /home/technik12345/vscodeprojects/rerandom && docker run --rm --network host \
-v "$PWD/nginx/nginx.conf:/etc/nginx/nginx.conf:ro" \
-v "$PWD/ssl/server.crt:/etc/ssl/certs/server.crt:ro" \
-v "$PWD/ssl/server.key:/etc/ssl/private/server.key:ro" \
-v "$PWD/nginx:/var/www/html:ro" \
nginx:alpine
```