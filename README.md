# learn-pub-sub-starter (Peril)

## This is the starter code used in Boot.dev's [Learn Pub/Sub](https://learn.boot.dev/learn-pub-sub) course

#### To run rabbitMQ

```bash
./rabbit.sh start

```

UI for rabbitMQ will be served on <http://localhost:15672> guest:guest

#### To stop

```bash
./rabbit.sh stop
```

#### For rabbitMQ logs

```bash
./rabbit.sh logs
```

#### To run server

```bash
go run ./cmd/server
```

#### To run client

```bash
go run ./cmd/client
```
