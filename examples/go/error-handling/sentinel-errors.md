# Sentinel Errors

Define package-level error values for callers to check with `errors.Is`.

## When to Use

- When callers need to distinguish specific failure modes
- When error identity matters more than error message
- When wrapping errors through multiple layers of a call stack

## Example

```go
package storage

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound   = errors.New("storage: not found")
	ErrConflict   = errors.New("storage: conflict")
	ErrValidation = errors.New("storage: validation failed")
)

func GetUser(id string) (User, error) {
	row, err := db.QueryRow("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, fmt.Errorf("user %s: %w", id, ErrNotFound)
		}
		return User{}, fmt.Errorf("querying user %s: %w", id, err)
	}
	return scanUser(row)
}

// Caller code:
func handleGetUser(id string) {
	user, err := storage.GetUser(id)
	switch {
	case errors.Is(err, storage.ErrNotFound):
		respondNotFound()
	case err != nil:
		respondInternalError(err)
	default:
		respondOK(user)
	}
}
```

## Key Points

- Use `errors.New` at package level for sentinel errors
- Wrap with `%w` to preserve the error chain for `errors.Is`
- Prefix sentinel error messages with the package name for clarity
- Prefer `errors.Is` over string comparison — it works through wrapped errors
- Don't overuse sentinels; only create them when callers need to branch on the error
