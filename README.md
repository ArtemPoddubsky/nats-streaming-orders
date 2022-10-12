# nats-streaming-orders 📦
Service that listens on nats-streaming, stores models into Postgres and inmemory cache to provide HTML UI.

<h2>Overview</h2>

Сервис подписывается на канал nats-streaming, записывает полученные данные в Postgres и in-memory cache. Поднимается http сервер, который отвечает за предоставления HTML UI, обращение к которому осуществляется через параметр через ID в запросе. Html страница генерируется на основе данных из кеша. В случае падения сервиса, cache восстанавливается из Postgres.

<h2>Model</h2>

    {
      "order_uid": "b563feb7b2b84b6test",
      "track_number": "WBILMTESTTRACK",
      "entry": "WBIL",
      "delivery": {
        "name": "Test Testov",
        "phone": "+9720000000",
        "zip": "2639809",
        "city": "Kiryat Mozkin",
        "address": "Ploshad Mira 15",
        "region": "Kraiot",
        "email": "test@gmail.com"
      },
      "payment": {
        "transaction": "b563feb7b2b84b6test",
        "request_id": "",
        "currency": "USD",
        "provider": "wbpay",
        "amount": 1817,
        "payment_dt": 1637907727,
        "bank": "alpha",
        "delivery_cost": 1500,
        "goods_total": 317,
        "custom_fee": 0
      },
      "items": [
        {
          "chrt_id": 9934930,
          "track_number": "WBILMTESTTRACK",
          "price": 453,
          "rid": "ab4219087a764ae0btest",
          "name": "Mascaras",
          "sale": 30,
          "size": "0",
          "total_price": 317,
          "nm_id": 2389212,
          "brand": "Vivienne Sabo",
          "status": 202
        }
      ],
      "locale": "en",
      "internal_signature": "",
      "customer_id": "test1",
      "delivery_service": "meest",
      "shardkey": "9",
      "sm_id": 99,
      "date_created": "2021-11-26T06:22:19Z",
      "oof_shard": "2"
    }

<h2>How to use</h2>

    make
    make publish    # публикация модели из примера.
    make test       # запуск тестов (требуется make publish).
