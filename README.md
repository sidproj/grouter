# grouter

A lightweight, dependency-free HTTP router for Go with support for:

- ✅ Dynamic route parameters (e.g. `/users/:id`)
- ✅ RESTful HTTP method routing (`GET`, `POST`, `PUT`, `DELETE`)
- ✅ Trie-based routing tree for efficient nested path resolution
- ✅ Easy path parameter extraction from requests
- ✅ Custom 404 page handling
- ✅ Centralized request dispatching

---

## 🚀 Features

⚡ Fast and minimalistic trie-based route matching

🔧 Clean route registration via Get, Post, Put, Delete functions

🔄 Built-in support for dynamic parameters like /users/:id

🧠 Parameters accessible from any handler using context

❌ Customizable 404 handler

🌐 CORS-friendly OPTIONS method support

---

## 📦 Installation

```bash
go get github.com/sidproj/grouter
```

## 🛠️ Usage

### Define Routes

```go
router.Get("/hello", func(w http.ResponseWriter, r \*http.Request) {
    fmt.Fprint(w, "Hello, world!")
})
```

```go
router.Post("/submit", func(w http.ResponseWriter, r \*http.Request) {
    fmt.Fprint(w, "Submitted successfully!")
})
```

### Dynamic Parameters

```go
router.Get("/users/:userId", func(w http.ResponseWriter, r \*http.Request) {
    params := router.GetPathParams(r)
    fmt.Fprintf(w, "User ID: %s", params["userId"])
})
```

### Launch the Server

```go
router.LoadRoutes()
http.ListenAndServe(":8080", nil)
```

---

## ✨ Example

```go
package main

import (
    "fmt"
    "net/http"
    "github.com/sidproj/grouter"
)

func main() {
    router.Get("/greet", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "Welcome to go-router!")
    })

    router.Get("/user/:id", func(w http.ResponseWriter, r *http.Request) {
        params := router.GetPathParams(r)
        fmt.Fprintf(w, "User ID: %s", params["id"])
    })

    router.Set404Path("views/404.html") // Optional custom 404 page

    router.LoadRoutes()
    fmt.Println("Server running on http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
```

## ⚙️ API Reference

### Route Registration

- List all public functions with brief descriptions:

```go
// Registers a GET route handler
router.Get(path string, handler http.HandlerFunc)

// Registers a POST route handler
router.Post(path string, handler http.HandlerFunc)

// Registers a PUT route handler
router.Put(path string, handler http.HandlerFunc)

// Registers a DELETE route handler
router.Delete(path string, handler http.HandlerFunc)

// Returns dynamic route parameters from the request context
router.GetPathParams(r *http.Request) map[string]string

// Initializes routing and registers the root HTTP handler
router.LoadRoutes()

// Sets the file path for the custom 404 response HTML
router.Set404Path(path string)
```

---

## 🔧 Internal Design

Router uses a recursive tree (RouterNode) structure to match incoming paths.

Each node can represent a static or dynamic path segment.

Dynamic segments start with : (e.g. :id) and are extracted into the request context.

All routing happens through a single wrapper() function to reduce handler duplication.

---

## ❗ Notes

Trailing slashes are normalized (/foo/ → /foo)

Dynamic parameters are stored in the request context and accessible via GetPathParams()

Routes not found are served via views/404.html (customize or change this)

OPTIONS preflight requests are handled based on the Access-Control-Request-Method header
