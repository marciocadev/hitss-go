## HITSS

# Detalhes do projeto

## scripts
A pasta scripts possuem scripts para a definição do banco de dados que será criádo no `Postgres` na criação do projeto

Basicamente será criado um schema `hitss` e a tabela `cliente` neste schema

# Inicialização do projeto
Para inicializar o projeto para teste é necessário executar os comandos abaixo

  > Eu utilizo o linux para construir caso seja utilizado o Windows para testar será necessário setar as variáveis de ambiente USER, PWD, QUEUE_INSERT

* A variável de ambiente `USER` e `PWD` são utilizadas tanto para a definição do usuário e da senha do `RabbitMQ` quanto do `Postgres`
* A variável de ambiente `QUEUE_INSERT` define o nome da fila do `RabbitMQ` que entrega dados para ser inserido no `Postgres`
```bash
$ USER=guest PWD=guest QUEUE_INSERT=insert DBNAME=hitssdb docker compose build
$ USER=guest PWD=guest QUEUE_INSERT=insert DBNAME=hitssdb docker compose up -d
```