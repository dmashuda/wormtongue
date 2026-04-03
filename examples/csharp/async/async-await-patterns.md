# Async/Await Patterns

Structured patterns for asynchronous programming in C# using async/await.

## When to Use

- I/O-bound operations (HTTP calls, database queries, file access)
- When you need non-blocking execution without manual thread management
- Parallelizing independent async operations with `Task.WhenAll`

## Example

```csharp
using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Threading.Tasks;

public class UserService
{
    private readonly HttpClient _httpClient;
    private readonly IUserRepository _repository;

    public UserService(HttpClient httpClient, IUserRepository repository)
    {
        _httpClient = httpClient;
        _repository = repository;
    }

    // Basic async/await — let the task propagate naturally
    public async Task<User> GetUserAsync(int id)
    {
        return await _repository.FindByIdAsync(id)
            ?? throw new NotFoundException($"User {id} not found");
    }

    // Parallel async — fire independent tasks, await together
    public async Task<UserProfile> GetFullProfileAsync(int userId)
    {
        var userTask = _repository.FindByIdAsync(userId);
        var ordersTask = _repository.GetOrdersAsync(userId);
        var prefsTask = _repository.GetPreferencesAsync(userId);

        await Task.WhenAll(userTask, ordersTask, prefsTask);

        return new UserProfile
        {
            User = userTask.Result,
            Orders = ordersTask.Result,
            Preferences = prefsTask.Result
        };
    }

    // Cancellation support via CancellationToken
    public async Task<string> FetchExternalDataAsync(
        string url,
        CancellationToken cancellationToken = default)
    {
        var response = await _httpClient.GetAsync(url, cancellationToken);
        response.EnsureSuccessStatusCode();
        return await response.Content.ReadAsStringAsync(cancellationToken);
    }
}
```

## Key Points

- Always use `Async` suffix for async method names by convention
- Prefer `await` over `.Result` or `.Wait()` to avoid deadlocks
- Use `Task.WhenAll` to parallelize independent operations
- Pass `CancellationToken` through the call chain for cooperative cancellation
- Avoid `async void` — use `async Task` so exceptions propagate correctly
