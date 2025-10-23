# Go Learning Guide - Based on Your Pokedex Project

## Table of Contents
1. [Key Concepts Review](#key-concepts-review)
2. [Common Go Idioms](#common-go-idioms)
3. [Practice Questions](#practice-questions)
4. [Advanced Exercises](#advanced-exercises)

---

## Key Concepts Review

### 1. Error Handling in Go

**The Golden Rule: Errors Come Last**

In Go, the convention is that error returns come as the **last** return value:

```go
// ✅ Correct - idiomatic Go
func fetchData(url string) (Data, error) {
    // implementation
}

// ❌ Incorrect - non-idiomatic
func fetchData(url string) (error, Data) {
    // implementation
}
```

**Why?** This makes it easier to scan code and identify error handling patterns. It's so consistent that Go developers expect it everywhere.

**Your code example (line 112 in main.go):**
```go
// Current (non-idiomatic)
func fetchLocations(l string, c *config) (error, LocationResponse) {
    // ...
    return nil, response
}

// Should be
func fetchLocations(l string, c *config) (LocationResponse, error) {
    // ...
    return response, nil
}
```

---

### 2. Concurrency and Mutexes

You did an excellent job with the cache! Let's break down why:

**sync.RWMutex vs sync.Mutex:**

```go
type Cache struct {
    mu    sync.RWMutex  // Read-Write mutex
    cache map[string]cacheEntry
}
```

- `sync.Mutex`: Only one goroutine can hold the lock (read OR write)
- `sync.RWMutex`: Multiple readers OR one writer

**When to use which:**

```go
// For reading (multiple goroutines can read simultaneously)
p.mu.RLock()
defer p.mu.RUnlock()
val := p.cache[key]

// For writing (exclusive access)
p.mu.Lock()
defer p.mu.Unlock()
p.cache[key] = value
```

**Key Point:** Always use `defer` to unlock! This ensures the lock is released even if a panic occurs.

---

### 3. Goroutines and Background Tasks

**Your reapLoop pattern (pokecache.go:37):**

```go
func (p *Cache) reapLoop(interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()  // Clean up resources

    for range ticker.C {  // Loop over channel values
        // cleanup logic
    }
}
```

**Important concepts:**
- `time.NewTicker` creates a channel that sends values at regular intervals
- `defer ticker.Stop()` prevents memory leaks
- `for range ticker.C` is idiomatic for consuming channel values
- This runs forever in a goroutine (started in NewCache)

---

### 4. Struct Tags for JSON

**Your LocationResponse (main.go:35):**

```go
type LocationResponse struct {
    Count     int        `json:"count"`
    Next      string     `json:"next"`
    Previous  string     `json:"previous"`
    Locations []Location `json:"results"`  // Note: JSON field is "results"
}
```

**Key learning:**
- Struct tags tell `json.Unmarshal` how to map JSON fields to struct fields
- Field names don't have to match JSON keys (see `Locations` ↔ `"results"`)
- Uppercase field names = exported (accessible from other packages)
- Lowercase field names = private (package-only)

---

### 5. The defer Keyword

**Three common uses:**

```go
// 1. Closing resources
defer res.Body.Close()

// 2. Unlocking mutexes
defer p.mu.Unlock()

// 3. Cleanup functions
defer ticker.Stop()
```

**Execution order:** LIFO (Last In, First Out)

```go
defer fmt.Println("1")
defer fmt.Println("2")
defer fmt.Println("3")
// Prints: 3, 2, 1
```

---

### 6. Pointer vs Value Receivers

**Your code uses pointer receivers correctly:**

```go
func (p *Cache) Add(key string, value []byte) {
    // p is a pointer to Cache
}

func (c *config) someMethod() error {
    // c is a pointer to config
}
```

**When to use pointers:**
- ✅ Method modifies the receiver
- ✅ Receiver is a large struct (avoid copying)
- ✅ Consistency (if some methods use pointers, all should)

**When values are okay:**
- Small, immutable structs
- Types like `time.Time`

---

## Common Go Idioms

### 1. Early Returns

**Instead of:**
```go
if condition {
    // do something
} else {
    // handle error
    return err
}
```

**Prefer:**
```go
if !condition {
    // handle error
    return err
}
// do something (happy path at the lowest indentation level)
```

Your code in `commandPrevMap` (line 158) could be simplified this way.

---

### 2. Variable Declaration

**Multiple forms:**

```go
// 1. Declaration with type
var response LocationResponse

// 2. Short declaration with inference
response := LocationResponse{}

// 3. Declaration where used (often better)
res, err := http.Get(url)
if err != nil {
    return err
}
```

**Tip:** Declare variables close to where they're used, not at the top of functions.

---

### 3. Zero Values

Go initializes variables to "zero values" automatically:

```go
var s string        // ""
var i int           // 0
var b bool          // false
var p *int          // nil
var slice []int     // nil
var m map[string]int // nil
```

**Your code (line 89-90):**
```go
locations := []Location{}  // Empty slice
response := LocationResponse{}  // Zero-valued struct
```

These get overwritten later, so you could declare them where first assigned.

---

### 4. Map Lookup Pattern

**Your command dispatch (line 201):**

```go
if cmd, ok := commands[command]; ok {
    err := cmd.callback(&apiConfig)
    // ...
} else {
    fmt.Println("Unknown command:", command)
}
```

This is the idiomatic way to check if a map key exists. The `ok` boolean tells you if the key was found.

---

## Practice Questions

### Level 1: Understanding Your Code

1. **Why does `Cache` use `sync.RWMutex` instead of `sync.Mutex`?**
   - What's the benefit in this specific use case?
   - When would a regular `Mutex` be sufficient?

2. **What happens if you remove `defer ticker.Stop()` from `reapLoop`?**
   - Is this a memory leak?
   - How would you detect this problem?

3. **In `fetchLocations`, what's the difference between these error checks?**
   ```go
   if err != nil {
       return err, response
   }
   ```
   vs
   ```go
   if res.StatusCode != http.StatusOK {
       return fmt.Errorf("failed to fetch data: %s", res.Status), response
   }
   ```

4. **Why is `apiConfig` passed as a pointer (`&apiConfig`) to callbacks?**
   - What would break if you passed it by value?

5. **What does the `defer res.Body.Close()` line prevent?**

### Level 2: Code Analysis

6. **Race Condition Challenge:**
   - If you removed all mutex locks from `pokecache.go`, describe exactly what could go wrong
   - What Go tool can detect race conditions? (Hint: look up `go run -race`)

7. **Goroutine Lifecycle:**
   ```go
   cache := pokecache.NewCache(5 * time.Minute)
   ```
   - When does the goroutine started in `NewCache` end?
   - Is this a problem? Why or why not?

8. **Error Handling:**
   - Find all places in `main.go` where errors are handled
   - Are there any places where errors are ignored?
   - Should they be?

9. **Memory Usage:**
   - What happens to the cache if you fetch 10,000 different URLs?
   - Does the reaper prevent unbounded growth?
   - What's a potential problem with the current reaping strategy?

10. **Type Safety:**
    - Why does the cache store `[]byte` instead of `LocationResponse`?
    - What are the tradeoffs?

### Level 3: Design Questions

11. **Command Pattern:**
    - How would you add a new command called `inspect` that takes a location name as an argument?
    - Where would you parse the argument from the user input?

12. **Error Types:**
    - How could you distinguish between "network error" and "invalid JSON" errors?
    - Look up custom error types in Go

13. **Testing:**
    - How would you test `commandMap` without making actual HTTP requests?
    - Look up "mocking" in Go

14. **Graceful Shutdown:**
    - How would you stop the reaper goroutine when the program exits?
    - Look up "context cancellation" in Go

15. **Cache Invalidation:**
    - How would you add a method to manually clear the entire cache?
    - Would you need to lock the mutex? Why?

### Level 4: Advanced Challenges

16. **Interface Design:**
    ```go
    type HttpClient interface {
        Get(url string) (*http.Response, error)
    }
    ```
    - How could you use this interface to make `fetchLocations` more testable?
    - Rewrite `fetchLocations` to accept this interface

17. **Context Usage:**
    - Research Go's `context` package
    - How would you add timeout support to the HTTP request in `fetchLocations`?

18. **Generics (Go 1.18+):**
    - Could you make the cache generic to store any type, not just `[]byte`?
    - What would the signature of `NewCache` look like?

19. **Buffered Channels:**
    - How could you use a buffered channel to limit concurrent HTTP requests?
    - Why might this be useful?

20. **Method Sets:**
    - If you changed `Add` to use a value receiver `(c Cache)` instead of `(c *Cache)`, what would break?
    - Try to explain using Go's method set rules

---

## Advanced Exercises

### Exercise 1: Add Caching Metrics

Add the following to your cache:
- Track cache hits vs misses
- Add a `Stats()` method that returns hit rate
- Make it thread-safe!

**Hints:**
- Add fields to `Cache` struct for counters
- Use atomic operations or mutex protection
- Calculate hit rate as: `hits / (hits + misses)`

---

### Exercise 2: Command with Arguments

Implement a new command `explore <location-name>` that shows Pokemon in that location.

**Requirements:**
- Parse command arguments in the REPL loop
- Make a new API call to `/location-area/{name}`
- Cache the results
- Handle "location not found" errors gracefully

**API endpoint:** `https://pokeapi.co/api/v2/location-area/{name}`

---

### Exercise 3: Graceful Shutdown

Make the reaper goroutine stoppable:

**Requirements:**
- Create a `Stop()` method on `Cache`
- Use a channel or context to signal the goroutine to exit
- Call `Stop()` when the program exits
- Bonus: Use `signal.Notify` to handle Ctrl+C gracefully

---

### Exercise 4: Configuration File

Add support for a config file:

**Requirements:**
- Create a `.pokedexrc` file with JSON config (cache duration, API URL, etc.)
- Load config at startup
- Use appropriate defaults if file doesn't exist
- Handle JSON parsing errors

---

### Exercise 5: Better Error Messages

Create custom error types:

```go
type NetworkError struct {
    URL string
    Err error
}

func (e *NetworkError) Error() string {
    return fmt.Sprintf("network error fetching %s: %v", e.URL, e.Err)
}
```

**Requirements:**
- Create at least 3 custom error types
- Use them in `fetchLocations`
- Add helpful context to errors
- Research `errors.Is` and `errors.As`

---

## Recommended Reading

### Official Resources
1. [Effective Go](https://go.dev/doc/effective_go) - The definitive style guide
2. [Go by Example](https://gobyexample.com/) - Practical examples
3. [A Tour of Go](https://go.dev/tour/) - Interactive tutorial

### Specific Topics
- **Concurrency:** "Go Concurrency Patterns" on the Go blog
- **Error Handling:** "Error handling and Go" on the Go blog
- **Testing:** "Testing" section in Effective Go

### Books
- "The Go Programming Language" by Donovan & Kernighan
- "Concurrency in Go" by Katherine Cox-Buday

---

## Quick Reference: Common Patterns

### Error Handling
```go
if err != nil {
    return fmt.Errorf("context: %w", err)  // Wrap error
}
```

### Defer Pattern
```go
resource, err := acquire()
if err != nil {
    return err
}
defer resource.Close()
// use resource
```

### Map Initialization
```go
m := make(map[string]int)  // Preferred
m := map[string]int{}      // Also fine
```

### Channel Iteration
```go
for value := range channel {
    // process value
}
```

### Struct Initialization
```go
// Named fields (preferred for clarity)
c := Config{
    URL: "https://api.example.com",
    Timeout: 5 * time.Second,
}

// Positional (fragile, avoid)
c := Config{"https://api.example.com", 5 * time.Second}
```

---

## Answers to Common Questions

**Q: When should I use pointers?**
A: When the struct is large, when you need to modify it, or when methods already use pointers.

**Q: What's the difference between `make` and `new`?**
A: `make` initializes slices, maps, and channels. `new` allocates zeroed memory and returns a pointer.

**Q: How do I handle optional parameters?**
A: Use the "functional options" pattern or accept a config struct.

**Q: Should I worry about goroutine leaks?**
A: Yes! Always ensure goroutines can exit. Use context cancellation or stop channels.

---

Good luck with your Go learning journey! Feel free to experiment with these exercises and questions. The best way to learn Go is by writing more Go code.
