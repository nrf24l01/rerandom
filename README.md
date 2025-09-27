# Rerandom

## HOW TO RUN
### On server
- Get google cloud account with google sheet api
- Configure .env like in .env.example
- Gen certs
```sh
cd nginx
openssl req -x509 -newkey rsa:2048 -nodes -keyout random-org.key -out random-org.crt -days 365 -subj "/C=RU/ST=Moscow/L=Moscow/O=ReRandom/OU=IT/CN=random.org" -addext "subjectAltName=DNS:random.org,DNS:www.random.org"
cd ..
```
- Run main docker compose
```sh
docker compose -f dev.docker-compose.yml up --build
```
- Create user
```sh
docker-compose exec -T postgres psql -U postgres -d postgres -c "INSERT INTO users (username, password) VALUES ('admin', '\$argon2id\$v=19\$m=16,t=2,p=1\$aFVpcHZsZlB1RngxbHdmSA\$fSTS6yUOxQCgoIJRJbzP1Q');"
```
creds - admin:12345678.
- Open admin panel
`<server-ip>/admin`

### On client
*I WILL HATE YOU IF YOU DO IT*
**ONLY ON YOUR MACHINE FOR EDUCATION PURPOSES**
- Download cert
```sh
wget <server-ip>/cert
```
- Add cert to thrusted
*differs on you os*
- Clear chrome hsts cache
[chrome://net-internals/#hsts](chrome://net-internals/#hsts)
Enter to delete security policies
`www.random.org` and `random.org`
- Reboot system or chrome
- hate yourself


## DEV
### How to run main ngixn
```sh
docker run --rm --network host \
-v "$PWD/nginx/nginx.conf:/etc/nginx/nginx.conf:ro" \
-v "$PWD/ssl/server.crt:/etc/ssl/certs/server.crt:ro" \
-v "$PWD/ssl/server.key:/etc/ssl/private/server.key:ro" \
-v "$PWD/nginx:/var/www/html:ro" \
nginx:alpine
```

### Keygen
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

### Create user
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

# Disclaimer and Terms of use

** The goal.** This project was created for research and educational purposes. It is intended for legitimate use only by those who act in good faith and within the framework of applicable laws and regulations.

**Lack of guarantees.** The project is provided "AS IS" and "WITH THE ASSUMPTION OF ERRORS". No express or implied warranties, including fitness for a particular purpose, are provided.

**Limitation of liability.** Neither the author(s) nor the contributors are liable for direct, indirect, incidental, special, punitive, or loss of data, income, or profits resulting from the use, modification, or distribution of this software, even if they have been advised of the possibility of such damages.

**Prohibition of malicious use.** It is strictly forbidden to use this code for:
- committing illegal actions,
- violating the security, privacy, or integrity of other people's systems,
- substitution of numbers for random.org,
- creation or distribution of malware.

**Community support and responsibility.** The project may not be supported or tested in all scenarios. Use it at your own risk.