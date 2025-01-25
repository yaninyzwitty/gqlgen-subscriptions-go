
# GQLGEN Subscriptions Go Project

This project primarily focuses on implementing message queues using Kafka, with an emphasis on consuming messages for gqlgen subscriptions.

*Messages are stored in Scylla DB.*

## Features

- Perform GraphQL mutations with ease, such as creating users, rooms, and messages.
- Receive real-time updates for messages via gqlgen subscriptions.
- Query data to verify successful database insertion.
- Remove eager fetching to reduce the cost of fetching objects.

## Credits

- [Astra-streaming-Kafka-integration](https://www.datastax.com/products/astra-streaming) for serverless Kafka integration.
- [Scylla DB Cloud, Cluster](https://cloud.scylladb.com) for serverless Scylla DB.

## Tools

- [Sonyflake](https://github.com/sony/sonyflake) for generating unique snowflake-like IDs, similar to Twitter IDs.
- [GQLGEN](https://gqlgen.com/getting-started/) for building servers with graphql
