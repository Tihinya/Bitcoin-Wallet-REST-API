
<div align=center>

<h1>Bitcoin-Wallet-REST-API</h1>

</div>

###


## 📄 Description

Simple Bitcoin Wallet API that facilitates transactions and provides essential functionalities such as viewing all transactions, checking the current balance in BTC and EUR, and initiating money transfers/spending.



## 🧑‍💻 Running

Change .example.env to .env and add:

    DB_HOST="localhost"
    DB_PORT="5432"
    DB_USER="postgres"
    DB_PASSWORD="secret"
    DB_NAME="postgres""

After run:

    docker compose up -d

    go run main.go

For check use PostMan:

For transfer

POST http://localhost:8080/transfer

Json:

    {
    "amount": 5
    }

For Get amount

GET http://localhost:8080/balance/get

For top-up:

POST http://localhost:8080/top-up

Json:

    {
    "amount": 10
    }



## Implementation


<div style="display:flex; align-items:center">
    <h4 style="padding-right:7px">Backend</h4>
    <img src="https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white" style="height:30px; padding-right:7px">
</div>

<div style="display:flex; align-items:center">
    <h4 style="padding-right:7px">Database</h4>
    <img src ="https://img.shields.io/badge/postgres-%23316192.svg?&style=for-the-badge&logo=postgresql&logoColor=white" style="height:30px; padding-right:7px">
</div>

<div style="display:flex; align-items:center">
    <h4 style="padding-right:7px">Container</h4>
    <img src="https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white" style="height:30px; padding-right:7px">


</div>




## 🤝 Authors

- **Stepan Tihinya** 
