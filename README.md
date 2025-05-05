# Clean Architecture - Desafio

## Passos para executar o desafio

### Inicializar os serviços

* com o docker client em funcionamento, execute:  
_docker-compose up -d_  
* certifique-se de criar um database "orders" no Mysql e a tabela com os seguintes campos: id do tipo string, price, tax e final_price do tipo float.  
* execute a aplicação com o comando, dentro da pasta cmd/ordersystem:  
_go run main.go wire_gen.go_  
* a aplicação subirá três serviços:
  * web server em http://localhost:8080  
  * gRPC server na porta 50051  
  * GraphQL server na porta 8080

### Testar o serviço Web Server
* utilize o arquivo /api/create_order.http e acione as requisições createOrder e order para criar um "order" e listas os "ordens", respectivamente.
### Testar o serviço gRPC Server  
* execute o comando abaixo para utilizar a ferramenta Evans:  
    _evans -r repl --host localhost --port 50051_  
* dentro da interface interativa, execute os comandos:  
    _call CreateOrder_ (para criar um "order")
    _call ListOrders_ (para listar os "orders")

### Testar o serviço GraphQL Server
* a partir do webbrowser, acione o Playground GraphQL em _http://localhost:8080_
* para criar um "order", execute uma mutation, conforme exemplo a seguir:
```mutation createOrder {
  createOrder (input: {id: "aaaaa", Price: 200.00, Tax: 35.00}) { 
  id
  Price
  Tax
  FinalPrice
	}
}
```
* para listar os "orders", execute uma query, conforme exemplo a seguir:
```query listOrders {
  listOrders { 
  id
  Price
  Tax
  FinalPrice
	}
}
```

## RabbitMQ
* O RabbitMQ é utilizado para implementar o serviço de mensageria.
* Para testar o serviço de mensageria, acesse _localhost://15672_ (usuário: guest, senha: guest)
* Crie uma fila chamada "orders"
* Dentro da fila, crie um "bind" denominado "amq.direct"
* a partir de agora, o serviço de mensageria está pronto para receber as mensagens, tanto para o serviço de criação de "order", quanto para o serviço de listagem de "orders"








