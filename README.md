# Fiber CLI

A command-line interface for Go Fiber framework, inspired by NestJS CLI. This tool helps you create and manage Go Fiber projects with support for modules and submodules.

## Features

- Project scaffolding with proper directory structure
- Component generation (controllers, services, models, repositories)
- Module and submodule support for hierarchical organization
- NestJS-like architecture adapted for Go Fiber

## Installation

### Option 1: Build from source

```bash
# Clone the repository
git clone https://github.com/yourusername/fiber-cli.git

# Navigate to the project directory
cd fiber-cli

# Build the CLI
go build -o fiber-cli

# Optional: Move to PATH
sudo mv fiber-cli /usr/local/bin/
```

### Option 2: Install with Go

```bash
go install github.com/yourusername/fiber-cli@latest
```

## Commands

### Create a new project

```bash
fiber-cli new my-app
```

This creates a new Go Fiber project with the following structure:

```
my-app/
├── config/
│   └── config.go
├── controllers/
├── middleware/
├── models/
├── modules/
│   └── module.go
├── repositories/
├── services/
├── utils/
├── go.mod
└── main.go
```

### Create a module

```bash
# Navigate to your project
cd my-app

# Create a root module
fiber-cli module users
```

This creates a new module with the following structure:

```
modules/
└── users/
    ├── controllers/
    │   └── users_controller.go
    ├── repositories/
    │   └── users_repository.go
    ├── services/
    │   └── users_service.go
    └── users_module.go
```

### Create a submodule

```bash
fiber-cli module profiles --parent /users
```

This creates a submodule under the users module:

```
modules/
└── users/
    ├── controllers/
    ├── profiles/
    │   ├── controllers/
    │   │   └── profiles_controller.go
    │   ├── repositories/
    │   │   └── profiles_repository.go
    │   ├── services/
    │   │   └── profiles_service.go
    │   └── profiles_module.go
    ├── repositories/
    ├── services/
    └── users_module.go
```

### Generate components

```bash
# Generate a controller
fiber-cli generate controller auth

# Generate a service in a module
fiber-cli generate service product --module modules/store

# Generate a repository in a submodule
fiber-cli generate repository comment --module modules/posts/comments
```

Available component types:
- controller
- service
- repository
- model
- middleware

## Module Structure

Each module contains its own:
- Controllers: Handle HTTP requests
- Services: Implement business logic
- Repositories: Handle data access

### Registering Modules in main.go

```go
package main

import (
	"log"
	"my-app/modules"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Register root module
	userModule := modules.NewUserModule()
	userModule.Register(app)

	// Create and register submodule
	profileModule := modules.NewProfileModule("/users")
	userModule.AddSubmodule(profileModule)

	// Start server
	log.Fatal(app.Listen(":3000"))
}
```

## Best Practices

### Module Organization

- **Root modules**: Represent core features (users, products, orders)
- **Submodules**: Represent related features (user-profiles, user-settings)
- **URL structure**: Will match your module hierarchy (/users/profiles)

### Route Definition

Routes are automatically registered when you call the `Register` method on your module:

```go
// From a controller file
func (c *UserController) Register(router fiber.Router) {
    users := router.Group("/users")
    users.Get("/", c.GetAll)
    users.Get("/:id", c.GetOne)
    // More routes...
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
