# go dependency_injection

`dependency_injection` is a lightweight, generics-based dependency injection container for Go, designed to simplify dependency management in modular applications.
It provides thread-safe registration, retrieval, and management of dependencies while promoting loose coupling and better testability.

## Features

- **Generics-based design**: Type-safe operations for dependency injection.
- **Thread-safe**: Built-in synchronization for concurrent environments.
- **Flexible**: Easily register, resolve, and remove dependencies.
- **Supports lazy initialization**: Automatically initialize and register dependencies when needed.

---

## Installation

Install the package via `go get`:

```bash
go get github.com/martinarisk/di/dependency_injection
```

## Main module example

```go
package main

import (
	"github.com/martinarisk/di/dependency_injection"
)

func main() {
	// Initialize the DI container
	di := dependency_injection.NewDependencyInjection()

	// Register configuration
	config := NewConfig()
	di.Add(config)

	// Register services and dependencies
	di.Add(NewService1(di))
	di.Add(NewService2(di))

	// Register controllers
	di.Add(NewController1(di))
	di.Add(NewController2(di))

	// Start application server
	server := NewServer(di)
	di.Add(server)
	server.RunForever()
}
```
## Controller example

```go
package controllers

import (
	"net/http"

	. "github.com/martinarisk/di/dependency_injection"
)

type ExampleController struct {
	usecase ExampleUsecase
}

func (c *ExampleController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// parse json
	data := parse(r)

	// Handle HTTP request using the injected usecase
	response := c.usecase.Execute(data)
	w.Write([]byte(response))
}

func NewController1(di *DependencyInjection) *ExampleController {
	// Resolve the usecase dependency
	usecase := MustNeed(di, NewUseCase)
	return &ExampleController{usecase: usecase}
}
```

## Use case example


```go
package usecases

import (
	. "github.com/martinarisk/di/dependency_injection"
)

type ExampleUseCase struct {
	dependency1 ExampleService
	dependency2 ExampleService
}

func (uc *ExampleUseCase) Execute() string {
	// Use dependencies to perform a task
	return uc.dependency1.DoSomething() + uc.dependency2.DoAnotherThing()
}

func NewUseCase(di *DependencyInjection) *ExampleUseCase {
	// Resolve dependencies
	dep1 := MustNeed(di, NewExampleService)
	dep2 := MustNeed(di, NewExampleService)
	return &ExampleUseCase{
		dependency1: dep1,
		dependency2: dep2,
	}
}
```

## Service example

```go
package services

import (
	. "github.com/martinarisk/di/dependency_injection"
)

type ExampleService struct{}

func (s *ExampleService) DoSomething() string {
	// Process an incoming request
	return "Did something"
}

func (s *ExampleService) DoAnotherThing() string {
	// Process an incoming request
	return "Did another thing"
}


func NewExampleService(di *DependencyInjection) *ExampleService {
	return &ExampleService{}
}
```
