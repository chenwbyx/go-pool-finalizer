# go-pool-finalizer
implement cgo-buf pools based on the sync.Pool standard library and use runtime.SetFinalizer to solve the problem of releasing cgo buffers periodically.
