# go dependency_injection

`dependency_injection` is a lightweight, generics-based dependency injection container for Go, designed to simplify dependency management in modular applications.
It provides thread-safe registration, retrieval, and management of dependencies while promoting loose coupling and better testability.

## Features

- **Generics-based design**: Type-safe operations for dependency injection.
- **Thread-safe**: Built-in synchronization for concurrent environments.
- **Flexible**: Easily register, resolve, and remove dependencies.
- **Supports lazy initialization**: Automatically initialize and register dependencies when needed.
- **Lifetimes**: Singleton (global instance), Transient (each time fresh allocation), Scoped (shared per request or scope), Pooled (maintain a small pool of objects that is auto garbage collected)
- **Examples**: See below
---

## Installation

Install the package via `go get`:

```bash
go get github.com/martinarisk/di/dependency_injection
```

# Singleton

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

# Scoped

This will create a new object, every time an object of unique type is requested. Total number of objects created is the number of Execute() calls.

```go
package main

import (
	. "github.com/martinarisk/di/dependency_injection"
)

type ExampleUseCase struct {
	*DependencyInjection
}

func (uc *ExampleUseCase) Execute(data string) string {


	di := NewScopedDependencyInjection(uc.DependencyInjection)

	// Scoped dependencies (once per request)
	dep1 := IExampleService(Ptr(MustNeed(di, NewExampleService)))
	dep2 := IExampleService(Ptr(MustNeed(di, NewExampleService)))
	


	// Use dependencies to perform a task
	return dep1.DoSomething() + dep2.DoAnotherThing()
}

func NewUseCase(di *DependencyInjection) *ExampleUseCase {
	// pass DependencyInjection to handlers pattern
	return &ExampleUseCase{DependencyInjection: di}
}
```


# Transient

This will create a new object, every time an object is requested. Total number of objects created is twice the number of Execute() calls (MustNeed used for two instances of NewExampleService).

```go
package main

import (
	. "github.com/martinarisk/di/dependency_injection"
)

type ExampleUseCase struct {
	*DependencyInjection
}

func (uc *ExampleUseCase) Execute(data string) string {


	di := NewTransientDependencyInjection(uc.DependencyInjection)

	// Transient dependencies (twice per request)
	dep1 := IExampleService(Ptr(MustNeed(di, NewExampleService)))
	dep2 := IExampleService(Ptr(MustNeed(di, NewExampleService)))
	


	// Use dependencies to perform a task
	return dep1.DoSomething() + dep2.DoAnotherThing()
}

func NewUseCase(di *DependencyInjection) *ExampleUseCase {
	// pass DependencyInjection to handlers pattern
	return &ExampleUseCase{DependencyInjection: di}
}
```


# Pooled

This will maintain a small pool of objects, which will be garbage collected after use automatically.

```go
package main

import (
	. "github.com/martinarisk/di/dependency_injection"
)

type ExampleUseCase struct {
	*DependencyInjection
}

func (uc *ExampleUseCase) Execute(data string) string {


	di := NewPooledDependencyInjection(uc.DependencyInjection)

	// Pooled dependencies (few times globally)
	dep1 := IExampleService(Ptr(MustNeed(di, NewExampleService)))
	dep2 := IExampleService(Ptr(MustNeed(di, NewExampleService)))
	


	// Use dependencies to perform a task
	return dep1.DoSomething() + dep2.DoAnotherThing()
}

func NewUseCase(di *DependencyInjection) *ExampleUseCase {
	// pass DependencyInjection to handlers pattern
	return &ExampleUseCase{DependencyInjection: di}
}
```

