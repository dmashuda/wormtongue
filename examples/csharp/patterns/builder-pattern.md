# Builder Pattern

Construct complex objects step by step with a fluent API.

## When to Use

- Objects with many optional parameters
- When construction requires validation or ordering constraints
- Immutable objects that can't be modified after creation

## Example

```csharp
public class HttpRequest
{
    public string Url { get; }
    public string Method { get; }
    public Dictionary<string, string> Headers { get; }
    public string? Body { get; }
    public TimeSpan Timeout { get; }

    private HttpRequest(Builder builder)
    {
        Url = builder.Url;
        Method = builder.Method;
        Headers = new Dictionary<string, string>(builder.Headers);
        Body = builder.Body;
        Timeout = builder.Timeout;
    }

    public class Builder
    {
        internal string Url { get; }
        internal string Method { get; private set; } = "GET";
        internal Dictionary<string, string> Headers { get; } = new();
        internal string? Body { get; private set; }
        internal TimeSpan Timeout { get; private set; } = TimeSpan.FromSeconds(30);

        public Builder(string url)
        {
            Url = url ?? throw new ArgumentNullException(nameof(url));
        }

        public Builder WithMethod(string method)
        {
            Method = method;
            return this;
        }

        public Builder WithHeader(string key, string value)
        {
            Headers[key] = value;
            return this;
        }

        public Builder WithBody(string body)
        {
            Body = body;
            return this;
        }

        public Builder WithTimeout(TimeSpan timeout)
        {
            Timeout = timeout;
            return this;
        }

        public HttpRequest Build()
        {
            if (Body != null && Method == "GET")
                throw new InvalidOperationException("GET requests cannot have a body");

            return new HttpRequest(this);
        }
    }
}

// Usage:
// var request = new HttpRequest.Builder("https://api.example.com/users")
//     .WithMethod("POST")
//     .WithHeader("Content-Type", "application/json")
//     .WithBody("{\"name\": \"Alice\"}")
//     .WithTimeout(TimeSpan.FromSeconds(10))
//     .Build();
```

## Key Points

- The constructor is private — only the builder can create instances
- Each `With*` method returns `this` for fluent chaining
- `Build()` is the place for validation and constraint checking
- The resulting object is immutable — all properties are read-only
- Required parameters go in the builder constructor; optional ones are methods
