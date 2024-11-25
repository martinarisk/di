// Package dependency_injection serves as a generics-based container for dependency injection.
package dependency_injection

import (
	"errors"
	"reflect"
	"sync"
)

// ErrDependencyNotFound is returned by Any(...) when no corresponding dependency is found.
var ErrDependencyNotFound = errors.New("dependency not found")

type dependencyInjection struct {
	dependencies map[string]map[interface{}]struct{}
	mutex sync.RWMutex
}

// DependencyInjection acts as a container for managing dependencies.
type DependencyInjection struct {
	info *dependencyInjection
}

// NewDependencyInjection initializes and returns a new instance of DependencyInjection.
func NewDependencyInjection() (di *DependencyInjection) {
	di = &DependencyInjection{info: &dependencyInjection{}}

	data := make(map[string]map[interface{}]struct{})
	
	di.info.dependencies = data

	return
}

// Add registers a dependency within the container.
func (di *DependencyInjection) Add(dep interface{}) {
	di.info.mutex.Lock()

	var t0 = "*" + reflect.TypeOf(dep).String()
	const t1 = ""

	if di.info.dependencies[t0] == nil {
		di.info.dependencies[t0] = make(map[interface{}]struct{})
	}
	di.info.dependencies[t0][dep] = struct{}{}

	if di.info.dependencies[t1] == nil {
		di.info.dependencies[t1] = make(map[interface{}]struct{})
	}
	di.info.dependencies[t1][dep] = struct{}{}

	di.info.mutex.Unlock()
}

// Remove unregisters a dependency from the container.
func (di *DependencyInjection) Remove(dep interface{}) {
	di.info.mutex.Lock()

	var t0 = "*" + reflect.TypeOf(dep).String()
	const t1 = ""

	delete(di.info.dependencies[t0], dep)
	delete(di.info.dependencies[t1], dep)

	di.info.mutex.Unlock()
}

// MustNeed injects a dependency of type T using the given constructor function and
// panics if the injection is unsuccessful.
func MustNeed[T any](di *DependencyInjection, newer func(di *DependencyInjection) *T) (result T) {
	err := Any[T](di, &result)
	if err != nil {
		result = *newer(di)
		di.Add(result)
	}
	return
}

// MustAny retrieves and returns a dependency of type T, panicking if the retrieval fails.
func MustAny[T any](di *DependencyInjection) (result T) {
	err := Any(di, &result)
	if err != nil {
		panic(err.Error())
	}
	return
}

// Any assigns a dependency of type T to the provided res pointer.
func Any[T any](di *DependencyInjection, res *T) error {
	if di == nil {
		return ErrDependencyNotFound
	}
	di.info.mutex.RLock()

	var t0 = reflect.TypeOf(res).String()
	const t1 = ""

	var deps0 = di.info.dependencies[t0]
	for dep := range deps0 {
		result, ok := (dep).(T)
		if ok {
			*res = result
			di.info.mutex.RUnlock()
			return nil
		}
	}
	var deps1 = di.info.dependencies[t1]
	for dep := range deps1 {
		result, ok := (dep).(T)
		if ok {
			*res = result
			di.info.mutex.RUnlock()
			return nil
		}
	}
	di.info.mutex.RUnlock()
	if t0 != "**dependency_injection.DependencyInjection" && (interface{}(di) != interface{}(*res)) && Any[*DependencyInjection](di, &di) == nil {
		return Any[T](di, res)
	}
	return ErrDependencyNotFound
}

// Ptr returns the pointer to any variable. Useful to make reference to values returned by MustAny() or MustNeed()
func Ptr[T any](val T) *T {
	return &val
}
