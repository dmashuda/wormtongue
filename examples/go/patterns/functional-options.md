# Functional Options Pattern

Configure structs using composable option functions instead of large config structs or many constructor parameters.

## When to Use

- Public API with many optional configuration knobs
- When you want sensible defaults with selective overrides
- Library code where backwards compatibility matters

## Example

```go
package server

import (
	"time"
)

type Server struct {
	addr         string
	readTimeout  time.Duration
	writeTimeout time.Duration
	maxConns     int
}

type Option func(*Server)

func WithReadTimeout(d time.Duration) Option {
	return func(s *Server) {
		s.readTimeout = d
	}
}

func WithWriteTimeout(d time.Duration) Option {
	return func(s *Server) {
		s.writeTimeout = d
	}
}

func WithMaxConns(n int) Option {
	return func(s *Server) {
		s.maxConns = n
	}
}

func New(addr string, opts ...Option) *Server {
	s := &Server{
		addr:         addr,
		readTimeout:  5 * time.Second,
		writeTimeout: 10 * time.Second,
		maxConns:     100,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Usage:
// srv := server.New(":8080",
//     server.WithReadTimeout(10*time.Second),
//     server.WithMaxConns(200),
// )
```

## Key Points

- The constructor sets sensible defaults; options override selectively
- Each option is a self-contained function — easy to add new ones without breaking the API
- Callers get a clean, readable configuration style
- This is the standard Go pattern used by `grpc.NewServer`, `zap.New`, and many others
