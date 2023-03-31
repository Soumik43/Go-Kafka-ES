
# Go - Kafka - Elasticsearch

A simple basic endpoint server which adds message to ElasticSearch from Kafka and fetches the messages.


## Run Locally

Clone the project

```bash
  git clone https://github.com/Soumik43/Go-Kafka-ES
```

Install dependencies

```bash
  go mod download
```

Go to the project directory

```bash
  cd /main
```

Build the code

```bash
  go build main.go
```

Run the server

```
  go run main.go
```


## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`KAFKA_ADDRESS (default : "localhost:9092")`

`ES_ADDRESS (default :"http://localhost:9200")` 

`ES_USERNAME (your elasticsearch username)` 

`ES_PASSWORD (your elasticsearch password)`

```
ES_INDEX 

(you'll have to create an index by the following steps)

curl -XPUT "http://192.168.29.218:9200/<your_index_name>" -H "kbn-xsrf: reporting" 
```

```
KAFKA_TOPIC

(you'll have to create an topic in kafka by the following steps)


Redirect to the folder where kafka is present 

(Start the ZooKeeper service in another terminal)
$ bin/zookeeper-server-start.sh config/zookeeper.properties

(Start the Kafka broker service in another terminal)
$ bin/kafka-server-start.sh config/server.properties

(Create a topic <your_topic_name> in another terminal)
bin/kafka-topics.sh --create --topic <your_topic_name> --bootstrap-server localhost:9092

``` 

## Roadmap

- Add parallel batching support for incoming message streams from kafka.

- Publish this to docker and containerize it.

