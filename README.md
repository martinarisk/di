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
	. "github.com/martinarisk/di/dependency_injection"

	"net/http/httptest"
)

type IServer interface {
	RunForever()
}

type Server struct{}

func (Server) RunForever() {

	// simulate http server running forever
	select {}
}

func NewServer(di *DependencyInjection) IServer {

	c := NewController1(di)

	// Register controllers
	di.Add(c)

	// simulate task
	go c.ServeHTTP(httptest.NewRecorder(), nil)

	// example
	return Server{}
}

func main() {
	// Initialize the DI container
	di := NewDependencyInjection()

	// Register configuration
	config := NewConfig()
	di.Add(config)

	// Start application server
	server := NewServer(di)
	server.RunForever()
}
```
## Controller example

```go
package main

import (
	"net/http"

	. "github.com/martinarisk/di/dependency_injection"
)

type ExampleController struct {
	usecase ExampleUseCase
}

func parse(r *http.Request) string {
	return "parsed"
}

func (c *ExampleController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	println("Serving")

	// parse json
	data := parse(r)

	// Handle HTTP request using the injected usecase
	response := c.usecase.Execute(data)

	println("Response:", response)

	w.WriteHeader(http.StatusOK)
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
package main

import (
	. "github.com/martinarisk/di/dependency_injection"
)

type ExampleUseCase struct {
	dependency1 IExampleService
	dependency2 IExampleService
}

func (uc *ExampleUseCase) Execute(data string) string {
	// Use dependencies to perform a task
	return uc.dependency1.DoSomething() + uc.dependency2.DoAnotherThing()
}

func NewUseCase(di *DependencyInjection) *ExampleUseCase {
	// Resolve dependencies
	dep1 := IExampleService(Ptr(MustNeed(di, NewExampleService)))
	dep2 := IExampleService(Ptr(MustNeed(di, NewExampleService)))
	return &ExampleUseCase{
		dependency1: dep1,
		dependency2: dep2,
	}
}
```

## Service example

```go
package main

import (
	. "github.com/martinarisk/di/dependency_injection"
)

type ExampleService struct{
	config IConfig
}

type IExampleService interface {
	DoSomething() string
	DoAnotherThing() string
}

func (s *ExampleService) DoSomething() string {
	// Process an incoming request
	return "Did something with " + s.config.GetVariable() + "\n"
}

func (s *ExampleService) DoAnotherThing() string {
	// Process an incoming request
	return "Did another thing with " + s.config.GetVariable() + "\n"
}

func NewExampleService(di *DependencyInjection) *ExampleService {

	config := MustAny[IConfig](di)

	return &ExampleService{config: config}
}
```

## Config example

```go
package main

type Config struct {
	Variable string
}

type IConfig interface {
	GetVariable() string
}

func (c *Config) GetVariable() string {
	return "custom variable"
}

func NewConfig() IConfig {
	return &Config{}
}
```
