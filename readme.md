# 🧺 Store Golang Microservices Sample

[![CI](https://github.com/mehdihadeli/store-golang-microservices/actions/workflows/ci.yml/badge.svg?branch=main&style=flat-square)](https://github.com/mehdihadeli/store-golang-microservices/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mehdihadeli/store-golang-microservice)](https://goreportcard.com/report/github.com/mehdihadeli/store-golang-microservice)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.19-61CFDD.svg?style=flat-square)
[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod&style=flat-square&color=ff69b4)](https://gitpod.io/#https://github.com/mehdihadeli/store-golang-microservices)

> A practical store sample, built with Golang and different software architecture and technologies like **Microservices Architecture**, **Vertical Slice Architecture** , **CQRS Pattern**, **Domain Driven Design (DDD)**, **Event Sourcing**, **Event Driven Architecture**. For communication between independent services, We use asynchronous messaging with using RabbitMQ, and sometimes we use synchronous communication for real-time communications with using REST and gRPC calls.

💡 This application is not business oriented and my focus is mostly on technical part, I just want to implement a sample with using different technologies, software architecture design, principles and all the thing we need for creating a microservices app.

🌀 This Application is `in-progress` and I will add new features and thecnologies over time. 

## Features
- ✅ Using `Vertical Slice Architecture` as a high level architecture
- ✅ Using `Event Driven Architecture` on top of RabbitMQ Message Broker with a custom [Event Bus](pkg/messaging/bus/)
- ✅ Using `Data Centeric Architecture` based on CRUD in [Catalogs Read Service](services/catalogs/read_service/)
- ✅ Using `Event Sourcing` in `Audit Based` services like [Orders Service](services/orders/) 
- ✅ Using `CQRS Pattern` and `Mediator Pattern`on top of [Go-MediatR](https://github.com/mehdihadeli/Go-MediatR) library
- ✅ Using RESTFul api with [Echo](https://github.com/labstack/echo) framework and using swagger with [swaggo/swag](https://github.com/swaggo/swag) library
- ✅ Using gRpc for internal service communication
- ✅ Using [go-playground/validator](https://github.com/go-playground/validator) for validating input data in the REST and gRpc
- ✅ Using `Postgres` and `EventStoreDB` for write databases as relational DB and `MongoDB` and `Elasric Search` for read databases
- ✅ Using `OpenTelemetry` for collection `Metrics` and `Distributed Tracing`
- ✅ Using docker and `docker-compose` for deployment

## Roadmap
- 🚧 Add Dependecny Injections with using [uber-go/fx](https://github.com/uber-go/fx)
- 🚧 Using `Domain Driven Design` in some of services like [Catalogs Write Service](services/catalogs/write_service/) 
- 🚧 Using `Helm` and `Kubernetes` for deployment
- 🚧 Add identity management and oauth

## Technologies - Libraries

- ✔️ **[`labstack/echo`](https://github.com/labstack/echo)** - High performance, minimalist Go web framework
- ✔️ **[`uber-go/zap`](https://github.com/uber-go/zap)** - Blazing fast, structured, leveled logging in Go.
- ✔️ **[`emperror/errors`](https://github.com/emperror/errors)** - Drop-in replacement for the standard library errors package and github.com/pkg/errors
- ✔️ **[`open-telemetry/opentelemetry-go`](https://github.com/open-telemetry/opentelemetry-go/)** - OpenTelemetry Go API and SDK
- ✔️ **[`open-telemetry/opentelemetry-go-contrib`](https://github.com/open-telemetry/opentelemetry-go-contrib)** - Collection of extensions for OpenTelemetry-Go.
- ✔️ **[`rabbitmq/amqp091-go`](https://github.com/rabbitmq/amqp091-go)** - An AMQP 0-9-1 Go client maintained by the RabbitMQ team. Originally by @streadway: `streadway/amqp`
- ✔️ **[`stretchr/testify`](https://github.com/stretchr/testify)** - A toolkit with common assertions and mocks that plays nicely with the standard library
- ✔️ **[`mehdihadeli/go-mediatr`](https://github.com/mehdihadeli/go-mediatr)** - Mediator pattern implementation in Golang and helpful in creating CQRS based applications.
- ✔️ **[`grpc-ecosystem/go-grpc-middleware`](https://github.com/grpc-ecosystem/go-grpc-middleware)** - Golang gRPC Middlewares: interceptor chaining, auth, logging, retries and more
- ✔️ **[`grpc/grpc-go`](https://github.com/grpc/grpc-go)** - The Go language implementation of gRPC. HTTP/2 based RPC
- ✔️ **[`elastic/go-elasticsearch`](https://github.com/elastic/go-elasticsearch)** - The official Go client for Elasticsearch
- ✔️ **[`avast/retry-go`](https://github.com/avast/retry-go)** - Simple golang library for retry mechanism
- ✔️ **[`ahmetb/go-linq`](https://github.com/ahmetb/go-linq)** - .NET LINQ capabilities in Go
- ✔️ **[`EventStore/EventStore-Client-Go`](https://github.com/EventStore/EventStore-Client-Go)** - Go Client for Event Store version 20 and above.
- ✔️ **[`olivere/elastic/v7`](https://github.com/olivere/elastic/v7)** - Deprecated: Use the official Elasticsearch client for Go at
- ✔️ **[`swaggo/swag`](https://github.com/swaggo/swag)** - Automatically generate RESTful API documentation with Swagger 2.0 for Go.
- ✔️ **[`prometheus/client_golang`](github.com/prometheus/client_golang)** - Prometheus instrumentation library for Go applications
- ✔️ **[`mongodb/mongo-go-driver`](https://github.com/mongodb/mongo-go-driver)** - The Go driver for MongoDB
- ✔️ **[`go-redis/redis`](https://github.com/go-redis/redis)** - Type-safe Redis client for Golang
- ✔️ **[`go-gorm/gorm`](https://github.com/go-gorm/gorm)** - The fantastic ORM library for Golang, aims to be developer friendly
- ✔️ **[`go-playground/validator`](https://github.com/go-playground/validator)** - Go Struct and Field validation, including Cross Field, Cross Struct, Map, Slice and Array diving

## System Architecture

![](./assets/system-architecture-diagram.png)

## Application Structure

In this project I used [vertical slice architecture](https://jimmybogard.com/vertical-slice-architecture/) or [Restructuring to a Vertical Slice Architecture](https://codeopinion.com/restructuring-to-a-vertical-slice-architecture/) also I used [feature folder structure](http://www.kamilgrzybek.com/design/feature-folders/) in this project.

- We treat each request as a distinct use case or slice, encapsulating and grouping all concerns from front-end to back.
- When We adding or changing a feature in an application in n-tire architecture, we are typically touching many different "layers" in an application. we are changing the user interface, adding fields to models, modifying validation, and so on. Instead of coupling across a layer, we couple vertically along a slice and each change affects only one slice.
- We `Minimize coupling` `between slices`, and `maximize coupling` `in a slice`.
- With this approach, each of our vertical slices can decide for itself how to best fulfill the request. New features only add code, we're not changing shared code and worrying about side effects. For implementing vertical slice architecture using cqrs pattern is a good match.

![](./assets/vertical-slice-architecture.jpg)


![](./assets/vsa2.png)

Also here I used [CQRS](https://www.eventecommerce.com/cqrs-pattern) for decompose my features to very small parts that makes our application:

- maximize performance, scalability and simplicity.
- adding new feature to this mechanism is very easy without any breaking change in other part of our codes. New features only add code, we're not changing shared code and worrying about side effects.
- easy to maintain and any changes only affect on one command or query (or a slice) and avoid any breaking changes on other parts
- it gives us better separation of concerns and cross cutting concern (with help of MediatR behavior pipelines) in our code instead of a big service class for doing a lot of things.

With using [CQRS](https://event-driven.io/en/cqrs_facts_and_myths_explained/), our code will be more aligned with [SOLID principles](https://en.wikipedia.org/wiki/SOLID), especially with:

- [Single Responsibility](https://en.wikipedia.org/wiki/Single-responsibility_principle) rule - because logic responsible for a given operation is enclosed in its own type.
- [Open-Closed](https://en.wikipedia.org/wiki/Open%E2%80%93closed_principle) rule - because to add new operation you don’t need to edit any of the existing types, instead you need to add a new file with a new type representing that operation.

Here instead of some [Technical Splitting](http://www.kamilgrzybek.com/design/feature-folders/) for example a folder or layer for our `services`, `controllers` and `data models` which increase dependencies between our technical splitting and also jump between layers or folders, We cut each business functionality into some vertical slices, and inner each of these slices we have [Technical Folders Structure](http://www.kamilgrzybek.com/design/feature-folders/) specific to that feature (command, handlers, infrastructure, repository, controllers, data models, ...).

Usually, when we work on a given functionality we need some technical things for example:

- API endpoint (Controller)
- Request Input (Dto)
- Request Output (Dto)
- Some class to handle Request, For example Command and Command Handler or Query and Query Handler
- Data Model

Now we could all of these things beside each other and it decrease jumping and dependencies between some layers or folders.

Keeping such a split works great with CQRS. It segregates our operations and slices the application code vertically instead of horizontally. In Our CQRS pattern each command/query handler is a separate slice. This is where you can reduce coupling between layers. Each handler can be a separated code unit, even copy/pasted. Thanks to that, we can tune down the specific method to not follow general conventions (e.g. use custom SQL query or even different storage). In a traditional layered architecture, when we change the core generic mechanism in one layer, it can impact all methods.

### Catalogs Service Structure

![](./assets/catalogs-service.png)

### High Level Structure

![](./assets/high-level.png)


## References

- [https://github.com/AleksK1NG/Go-CQRS-Kafka-gRPC-Microservices](https://github.com/AleksK1NG/Go-CQRS-Kafka-gRPC-Microservices)
- [https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example](https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example)
- [https://github.com/kgrzybek/modular-monolith-with-ddd](https://github.com/kgrzybek/modular-monolith-with-ddd)
- [https://github.com/jbogard/ContosoUniversityDotNetCore-Pages](https://github.com/jbogard/ContosoUniversityDotNetCore-Pages)
