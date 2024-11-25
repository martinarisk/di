#  Dependency Injection for Golang

## Constructor: `NewDependencyInjection`
```go
func NewDependencyInjection() (di *DependencyInjection)
```
The `NewDependencyInjection` function initializes a new instance of the `DependencyInjection` container. Typically, it is called once at the application's entry point (e.g., in main) to create the DI container that manages all dependencies.

Example:
```go
di := NewDependencyInjection()
```
## Adding and Removing Dependencies

Add:

```go
di.Add(obj interface{})
```

Registers an object with the DI container. Once added, the object is available for resolution.

Example:
```go
config := NewConfig()
di.Add(config)
```

Remove:
```go
di.Remove(obj interface{})
```

Removes an object from the DI container.

Example:
```go
di.Remove(config)
```

## Resolving Dependencies
### Non-interface Object Creation

Use `MustNeed` to resolve a dependency or create a new instance using a constructor function (if dependency not found).
```go
dep1 := MustNeed(di, NewExampleService)
```
`MustNeed`:

Resolves the dependency if it already exists in the DI container.
If not found, it calls the provided constructor (`NewExampleService` in this case), registers the result, and returns it.

Example:
```go
usecase := MustNeed(di, NewUseCase)
```

### Interface Object Creation

When dealing with interfaces, wrap the resolved object in a type conversion.
```go
dep1 := IExampleService(Ptr(MustNeed(di, NewExampleService)))
```
Purpose: Converts the resolved object to the required interface type.

**Ptr**: Helper function that dereferences the value for pointer-based interfaces.

### Utility Functions

#### Ptr:
```go
func Ptr[T any](val T) *T
```
Dereferences the value and returns a pointer to it.

Example:
```go
ptr := Ptr(someValue)
```

#### Any:
```go
func Any[T any](di *DependencyInjection, res *T) error
```
Attempts to resolve a dependency and populate res. Returns an error if the dependency is not found.

Example:
```go
var config IConfig
err := Any(di, &config)
```

#### MustAny:
```go
func MustAny[T any](di *DependencyInjection) (result T)
```
Resolves a dependency and returns it. Panics if the dependency is not found.

Example:
```go
config := MustAny[IConfig](di)
```

#### MustNeed:
```go
func MustNeed[T any](di *DependencyInjection, newer func(di *DependencyInjection) *T) (result T)
```
Resolves or creates a dependency using the provided constructor function. Panics if the dependency cannot be created.

Example:
```go
service := MustNeed(di, NewExampleService)
```

## Using Lifetimes in Dependency Injection

The DI container supports various lifetimes to manage the lifecycle of dependencies.

### Singleton:
Default lifetime. A single instance is shared globally.

### Scoped:
A new instance is created for each logical "scope" (e.g., per request). Use NewScopedDependencyInjection.
```go
func NewScopedDependencyInjection(di *DependencyInjection) *DependencyInjection
```
Example:
```go
scopedDi := NewScopedDependencyInjection(di)
```

### Transient:
A new instance is created every time the dependency is requested. Use NewTransientDependencyInjection.
```go
func NewTransientDependencyInjection(di *DependencyInjection) *DependencyInjection
```
Example:
```go
transientDi := NewTransientDependencyInjection(di)
```

### Pooled:
Maintains a pool of objects that are reused and garbage-collected as needed. Use NewPooledDependencyInjection.

```go
func NewPooledDependencyInjection(di *DependencyInjection) *DependencyInjection
```

Example:
```go
pooledDi := NewPooledDependencyInjection(di)
```

## Reserved Internal Methods

`IsTransient()` and `SetTransient()`:
Reserved for internal use to manage transient dependencies.

## Features Recap

- Generics-based design for type safety.
- Thread-safe operations for concurrent use.
- Supports Singleton, Scoped, Transient, and Pooled lifetimes.
- Enables lazy initialization for dependencies.

By adopting these patterns, `dependency_injection` simplifies dependency management in Go applications, promoting modularity, scalability, and testability.
