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

## Create user
```sh
cd control_tool
go run main.go create-user guest
```

# Отказ от ответственности и условия использования

**Цель.** Данный проект создан в исследовательских и образовательных целях. Он предназначен для законного использования только теми, кто действует добросовестно и в рамках применимых законов и правил.

**Отсутствие гарантий.** Проект предоставляется «КАК ЕСТЬ» и «С ДОПУЩЕНИЕМ ОШИБОК». Никакие явные или подразумеваемые гарантии, включая пригодность для конкретной цели, не предоставляются.

**Ограничение ответственности.** Ни автор(ы), ни контрибьюторы не несут ответственности за прямые, косвенные, случайные, особые, штрафные убытки или потерю данных, дохода или прибыли, возникшие в результате использования, изменения или распространения данного программного обеспечения, даже если они были уведомлены о возможности таких убытков.

**Запрет на вредоносное использование.** Категорически запрещено использовать этот код для:
- совершения противоправных действий,
- нарушения безопасности, приватности или целостности чужих систем,
- подмены чисел на random.org,
- создания или распространения вредоносного ПО.

**Поддержка и ответственность сообщества.** Проект может не поддерживаться и не тестироваться на всех сценариях. Используйте на собственный риск.
