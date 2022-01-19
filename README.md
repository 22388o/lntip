<h1 align="center">Welcome to lntip 👋</h1>
<p>
  <img alt="Version" src="https://img.shields.io/badge/version-1.0.0-blue.svg?cacheSeconds=2592000" />
  <a href="#" target="_blank">
    <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" />
  </a>
</p>

> lntip is a bot that allows you to tip Discord users with Bitcoin through the Lightning Network 

## Demo
Join this [Discord server](https://discord.gg/EjgmBe4HhU) to test it out!  

## Usage
There are two ways to use this bot, either self-hosted or using the hosted version.

### Hosted

Simply invite the bot to your server using this URL.
> https://discord.com/api/oauth2/authorize?client_id=932252820768448552&permissions=137439276096&scope=bot

### Self-hosted using Docker
You need:
- Docker & docker-compose
- A Bitcoin node (Bitcoin Core recommended)
- An LND node (With Tor recommended)

Here is a docker-compose.yml sample, edit it to your needs.  

```docker
version: "3.1"
services:
  db:
    container_name: lntipdb
    image: mariadb
    volumes:
      - ./data:/var/lib/mysql
    environment:
      - MARIADB_ROOT_PASSWORD=your_password
    ports: # optional
      - "127.0.0.1:3306:3306"
    networks:
      - db
  lntip:
    container_name: lntip
    build:
      context: lntip
      dockerfile: Dockerfile
    volumes:
      - ./config.yml:/lntip/config.yml
      - /path/to/tls.cert:/lntip/tls.cert
      - /path/to/macaroons/mainnet:/lntip/macaroons
    networks:
      - db
      - lnd

networks:
  db:
  lnd:
    external: true
```

Then you must configure a config.yml. Here is a sample:  
```yaml
bot:
  token: your.discord.token

database:
  name: lntip
  host: db
  port: 3306
  user: root
  password: your_password

lnd:
  host: lnd
  tls_path: /lntip/tls.cert
  macaroon_path: /lntip/macaroons
  network: mainnet
```

Start the database and the bot:

```bash
$ docker-compose up -d --build
```

Initialize the database:
```
echo 'CREATE DATABASE lntip' | docker exec -it lntipdb mysql -uroot -pyour_pasword

# Run SQL migrations
docker run -v $PWD/migrations:/migrations --net=host migrate/migrate -path=/migrations/ -database "mysql://root:your_password@(localhost:3306)/lntip" up
```

🎉

## Author

👤 Aurèle Oulès

* Github: [@aureleoules](https://github.com/aureleoules)
* LinkedIn: [@https:\/\/www.linkedin.com\/in\/aureleoules\/](https://linkedin.com/in/https:\/\/www.linkedin.com\/in\/aureleoules\/)

## License
[MIT](https://github.com/aureleoules/lntip/blob/master/LICENSE) - [Aurèle Oulès](https://www.aureleoules.com)
