package dependency_injection

import "runtime"

// NewTransientDependencyInjection creates a DependencyInjection for injection using
// the Transient lifetime. Each MustNew(...) object made from the result is newly allocated.
func NewTransientDependencyInjection(di *DependencyInjection) (*DependencyInjection) {
	child := NewDependencyInjection()
	// must be before SetTransient
	child.Add(di)
	// freeze it
	child.SetTransient(true)
	return child
}

// NewScopedDependencyInjection creates a DependencyInjection for injection using
// the Scoped lifetime. Each MustNew(...) object made from the result is scoped,
// multiple instances for equal type objects are not newly allocated (one singleton per type).
func NewScopedDependencyInjection(di *DependencyInjection) (*DependencyInjection) {
	child := NewDependencyInjection()
	child.Add(di)
	return child
}

// NewPooledDependencyInjection creates a DependencyInjection for injection using
// the Pooled lifetime. Each MustNew(...) object made from the result is from a pool
// of small number of objects, dynamically adjusting to load.
func NewPooledDependencyInjection(di *DependencyInjection) (*DependencyInjection) {
	return Ptr(MustNeed(di, func (parent *DependencyInjection) (*DependencyInjection) {
		child := NewDependencyInjection()
		clone := Ptr(*parent)
		clone.Add(child)
		runtime.SetFinalizer(clone, func(s *DependencyInjection) {
			clone.Remove(*child)
			clone = nil
			parent = nil
			child = nil
		})
		return clone
	}))
}
