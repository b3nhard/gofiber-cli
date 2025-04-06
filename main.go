package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// CommandData represents the structure of the nest-cli.json equivalent
type CommandData struct {
	Collection   string `yaml:"collection"`
	SourceRoot   string `yaml:"sourceRoot"`
	ProjectName  string `yaml:"projectName"`
	ProjectType  string `yaml:"projectType"`
	ProjectRoot  string `yaml:"projectRoot"`
	Monorepo     bool   `yaml:"monorepo"`
	Dependencies struct {
		Core    []string `yaml:"core"`
		Dev     []string `yaml:"dev"`
		Test    []string `yaml:"test"`
		Scripts []string `yaml:"scripts"`
	} `yaml:"dependencies"`
}

// Improved template map with more component types
var templates = map[string]string{
	"controller": `package controllers

import (
	"github.com/gofiber/fiber/v2"
)

// {{.Name}}Controller handles {{.Resource}} related requests
type {{.Name}}Controller struct {
	service *services.{{.Name}}Service
}

// New{{.Name}}Controller creates a new {{.Name}}Controller
func New{{.Name}}Controller(service *services.{{.Name}}Service) *{{.Name}}Controller {
	return &{{.Name}}Controller{
		service: service,
	}
}

// Register registers routes for the {{.Name}}Controller
func (c *{{.Name}}Controller) Register(router fiber.Router) {
	{{.Resource}} := router.Group("/{{.ResourceURL}}")

	{{.Resource}}.Get("/", c.GetAll)
	{{.Resource}}.Get("/:id", c.GetOne)
	{{.Resource}}.Post("/", c.Create)
	{{.Resource}}.Put("/:id", c.Update)
	{{.Resource}}.Delete("/:id", c.Delete)
}

// GetAll returns all {{.Resource}}
func (c *{{.Name}}Controller) GetAll(ctx *fiber.Ctx) error {
	items, err := c.service.GetAll()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return ctx.JSON(fiber.Map{
		"data": items,
	})
}

// GetOne returns one {{.Resource}} by ID
func (c *{{.Name}}Controller) GetOne(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	item, err := c.service.GetByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "{{.Resource}} not found",
		})
	}
	
	return ctx.JSON(fiber.Map{
		"data": item,
	})
}

// Create creates a new {{.Resource}}
func (c *{{.Name}}Controller) Create(ctx *fiber.Ctx) error {
	dto := new(dtos.Create{{.Name}}DTO)
	
	if err := ctx.BodyParser(dto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	
	// Validate the DTO
	if err := dto.Validate(); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	item, err := c.service.Create(dto)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": item,
	})
}

// Update updates a {{.Resource}} by ID
func (c *{{.Name}}Controller) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	dto := new(dtos.Update{{.Name}}DTO)
	
	if err := ctx.BodyParser(dto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	
	// Validate the DTO
	if err := dto.Validate(); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	item, err := c.service.Update(id, dto)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return ctx.JSON(fiber.Map{
		"data": item,
	})
}

// Delete deletes a {{.Resource}} by ID
func (c *{{.Name}}Controller) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	err := c.service.Delete(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return ctx.Status(fiber.StatusNoContent).Send(nil)
}
`,
	"service": `package services

import (
	"{{.ProjectImport}}/{{.ModulePath}}/dtos"
	"{{.ProjectImport}}/{{.ModulePath}}/repositories"
	"{{.ProjectImport}}/models"
)

// {{.Name}}Service provides functionality for {{.Resource}}
type {{.Name}}Service struct {
	repository *repositories.{{.Name}}Repository
}

// New{{.Name}}Service creates a new {{.Name}}Service
func New{{.Name}}Service(repository *repositories.{{.Name}}Repository) *{{.Name}}Service {
	return &{{.Name}}Service{
		repository: repository,
	}
}

// GetAll returns all {{.Resource}}
func (s *{{.Name}}Service) GetAll() ([]models.{{.Name}}, error) {
	return s.repository.FindAll()
}

// GetByID returns a {{.Resource}} by ID
func (s *{{.Name}}Service) GetByID(id string) (*models.{{.Name}}, error) {
	return s.repository.FindByID(id)
}

// Create creates a new {{.Resource}}
func (s *{{.Name}}Service) Create(dto *dtos.Create{{.Name}}DTO) (*models.{{.Name}}, error) {
	// Map DTO to model
	model := &models.{{.Name}}{
		ID:        generateID(), // Implement ID generation
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		// Map other fields
	}
	
	if err := s.repository.Create(model); err != nil {
		return nil, err
	}
	
	return model, nil
}

// Update updates a {{.Resource}} by ID
func (s *{{.Name}}Service) Update(id string, dto *dtos.Update{{.Name}}DTO) (*models.{{.Name}}, error) {
	// Find existing model
	model, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	
	// Update model fields from DTO
	model.UpdatedAt = time.Now()
	// Update other fields
	
	if err := s.repository.Update(model); err != nil {
		return nil, err
	}
	
	return model, nil
}

// Delete deletes a {{.Resource}} by ID
func (s *{{.Name}}Service) Delete(id string) error {
	return s.repository.Delete(id)
}

// Helper function to generate IDs
func generateID() string {
	// Implement your ID generation logic here
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
`,
	"model": `package models

import (
	"time"
)

// {{.Name}} represents a {{.Resource}} entity
type {{.Name}} struct {
	ID        string    ` + "`json:\"id\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
	// Add your fields here
}
`,
	"repository": `package repositories

import (
	"errors"
	"{{.ProjectImport}}/models"
)

// {{.Name}}Repository handles database operations for {{.Resource}}
type {{.Name}}Repository struct {
	// Add database connection or ORM here
	store map[string]*models.{{.Name}} // In-memory store for demo
}

// New{{.Name}}Repository creates a new {{.Name}}Repository
func New{{.Name}}Repository() *{{.Name}}Repository {
	return &{{.Name}}Repository{
		store: make(map[string]*models.{{.Name}}),
	}
}

// FindAll returns all {{.Resource}}
func (r *{{.Name}}Repository) FindAll() ([]models.{{.Name}}, error) {
	items := make([]models.{{.Name}}, 0, len(r.store))
	for _, item := range r.store {
		items = append(items, *item)
	}
	return items, nil
}

// FindByID finds a {{.Resource}} by ID
func (r *{{.Name}}Repository) FindByID(id string) (*models.{{.Name}}, error) {
	item, ok := r.store[id]
	if !ok {
		return nil, errors.New("{{.Resource}} not found")
	}
	return item, nil
}

// Create creates a new {{.Resource}}
func (r *{{.Name}}Repository) Create({{.Resource}} *models.{{.Name}}) error {
	r.store[{{.Resource}}.ID] = {{.Resource}}
	return nil
}

// Update updates a {{.Resource}} by ID
func (r *{{.Name}}Repository) Update({{.Resource}} *models.{{.Name}}) error {
	_, ok := r.store[{{.Resource}}.ID]
	if !ok {
		return errors.New("{{.Resource}} not found")
	}
	r.store[{{.Resource}}.ID] = {{.Resource}}
	return nil
}

// Delete deletes a {{.Resource}} by ID
func (r *{{.Name}}Repository) Delete(id string) error {
	_, ok := r.store[id]
	if !ok {
		return errors.New("{{.Resource}} not found")
	}
	delete(r.store, id)
	return nil
}
`,
	"dto": `package dtos

import (
	"errors"
)

// Create{{.Name}}DTO is the data transfer object for creating a {{.Resource}}
type Create{{.Name}}DTO struct {
	// Add fields here
	Name string ` + "`json:\"name\"`" + `
	// Add more fields
}

// Update{{.Name}}DTO is the data transfer object for updating a {{.Resource}}
type Update{{.Name}}DTO struct {
	// Add fields here
	Name string ` + "`json:\"name\"`" + `
	// Add more fields
}

// Validate validates the Create{{.Name}}DTO
func (dto *Create{{.Name}}DTO) Validate() error {
	if dto.Name == "" {
		return errors.New("name is required")
	}
	// Add more validations
	return nil
}

// Validate validates the Update{{.Name}}DTO
func (dto *Update{{.Name}}DTO) Validate() error {
	// Add validations
	return nil
}
`,
	"module": `package {{.ModulePackage}}

import (
	"{{.ProjectImport}}/{{.ModulePath}}/controllers"
	"{{.ProjectImport}}/{{.ModulePath}}/repositories"
	"{{.ProjectImport}}/{{.ModulePath}}/services"

	"github.com/gofiber/fiber/v2"
)

// {{.Name}}Module represents the {{.Resource}} module
type {{.Name}}Module struct {
	controller    *controllers.{{.Name}}Controller
	service       *services.{{.Name}}Service
	repository    *repositories.{{.Name}}Repository
	submodules    []Module
	prefix        string
}

// New{{.Name}}Module creates a new {{.Name}}Module
func New{{.Name}}Module() *{{.Name}}Module {
	repository := repositories.New{{.Name}}Repository()
	service := services.New{{.Name}}Service(repository)
	controller := controllers.New{{.Name}}Controller(service)

	return &{{.Name}}Module{
		controller:    controller,
		service:       service,
		repository:    repository,
		submodules:    []Module{},
		prefix:        "/{{.ResourceURL}}",
	}
}

// AddSubmodule adds a submodule to this module
func (m *{{.Name}}Module) AddSubmodule(submodule Module) {
	m.submodules = append(m.submodules, submodule)
}

// SetPrefix sets the URL prefix for this module
func (m *{{.Name}}Module) SetPrefix(prefix string) {
	m.prefix = prefix
}

// Register registers this module's routes
func (m *{{.Name}}Module) Register(app *fiber.App) {
	// Create a router group for this module
	router := app.Group(m.prefix)

	// Register this module's controller
	m.controller.Register(router)

	// Register all submodules
	for _, submodule := range m.submodules {
		submodule.Register(app)
	}
}
`,
	"submodule": `package {{.ModulePackage}}

import (
	"{{.ProjectImport}}/{{.ModulePath}}/controllers"
	"{{.ProjectImport}}/{{.ModulePath}}/repositories"
	"{{.ProjectImport}}/{{.ModulePath}}/services"

	"github.com/gofiber/fiber/v2"
)

// {{.Name}}Module represents the {{.Resource}} submodule
type {{.Name}}Module struct {
	controller    *controllers.{{.Name}}Controller
	service       *services.{{.Name}}Service
	repository    *repositories.{{.Name}}Repository
	submodules    []Module
	prefix        string
	parent        string
}

// New{{.Name}}Module creates a new {{.Name}}Module
func New{{.Name}}Module(parentPrefix string) *{{.Name}}Module {
	repository := repositories.New{{.Name}}Repository()
	service := services.New{{.Name}}Service(repository)
	controller := controllers.New{{.Name}}Controller(service)

	return &{{.Name}}Module{
		controller:    controller,
		service:       service,
		repository:    repository,
		submodules:    []Module{},
		prefix:        "/{{.ResourceURL}}",
		parent:        parentPrefix,
	}
}

// AddSubmodule adds a submodule to this module
func (m *{{.Name}}Module) AddSubmodule(submodule Module) {
	m.submodules = append(m.submodules, submodule)
}

// SetPrefix sets the URL prefix for this module
func (m *{{.Name}}Module) SetPrefix(prefix string) {
	m.prefix = prefix
}

// Register registers this module's routes
func (m *{{.Name}}Module) Register(app *fiber.App) {
	// Create a router group for this module
	fullPath := m.parent + m.prefix
	router := app.Group(fullPath)

	// Register this module's controller
	m.controller.Register(router)

	// Register all submodules
	for _, submodule := range m.submodules {
		submodule.Register(app)
	}
}
`,
	"middleware": `package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// {{.Name}}Middleware provides {{.Resource}} related middleware
func {{.Name}}Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Implement middleware logic here
		return c.Next()
	}
}
`,
	"guard": `package guards

import (
	"github.com/gofiber/fiber/v2"
)

// {{.Name}}Guard handles authorization for {{.Resource}} resources
func {{.Name}}Guard() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Implement authorization logic here
		// For example, check if the user has the right permissions

		// If not authorized
		// return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		//     "error": "Unauthorized access",
		// })

		// If authorized, continue
		return c.Next()
	}
}
`,
	"pipe": `package pipes

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

// {{.Name}}ValidationPipe validates requests for {{.Resource}}
func {{.Name}}ValidationPipe(schema interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create a new instance of the schema
		val := reflect.New(reflect.TypeOf(schema)).Interface()
		
		// Parse request body into the schema
		if err := c.BodyParser(val); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}
		
		// If the schema has a Validate method, call it
		if validator, ok := val.(interface{ Validate() error }); ok {
			if err := validator.Validate(); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
		}
		
		// Store the validated data in locals for the next handler
		c.Locals("data", val)
		
		return c.Next()
	}
}
`,
	"test": `package {{.PackageName}}_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	
	"{{.ProjectImport}}/{{.ModulePath}}/controllers"
	"{{.ProjectImport}}/{{.ModulePath}}/repositories"
	"{{.ProjectImport}}/{{.ModulePath}}/services"
)

func Setup{{.Name}}Test() (*fiber.App, *controllers.{{.Name}}Controller) {
	app := fiber.New()
	repository := repositories.New{{.Name}}Repository()
	service := services.New{{.Name}}Service(repository)
	controller := controllers.New{{.Name}}Controller(service)
	
	return app, controller
}

func Test{{.Name}}GetAll(t *testing.T) {
	app, controller := Setup{{.Name}}Test()
	
	// Setup route
	app.Get("/{{.ResourceURL}}", controller.GetAll)
	
	// Create request
	req := httptest.NewRequest(http.MethodGet, "/{{.ResourceURL}}", nil)
	req.Header.Set("Content-Type", "application/json")
	
	// Execute request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	// Check response body
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "data")
}

func Test{{.Name}}GetOne(t *testing.T) {
	app, controller := Setup{{.Name}}Test()
	
	// Setup route
	app.Get("/{{.ResourceURL}}/:id", controller.GetOne)
	
	// Create request
	req := httptest.NewRequest(http.MethodGet, "/{{.ResourceURL}}/1", nil)
	req.Header.Set("Content-Type", "application/json")
	
	// Execute request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	// Expect 404 since no items exist yet
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func Test{{.Name}}Create(t *testing.T) {
	app, controller := Setup{{.Name}}Test()
	
	// Setup route
	app.Post("/{{.ResourceURL}}", controller.Create)
	
	// Create request with JSON body
	body := strings.NewReader(` + "`{\"name\":\"Test {{.Name}}\"}`" + `)
	req := httptest.NewRequest(http.MethodPost, "/{{.ResourceURL}}", body)
	req.Header.Set("Content-Type", "application/json")
	
	// Execute request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}
`,
	"filter": `package filters

import (
	"github.com/gofiber/fiber/v2"
)

// {{.Name}}ErrorFilter handles errors for {{.Resource}}
func {{.Name}}ErrorFilter() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Get the status code from the error
		code := fiber.StatusInternalServerError
		
		// You can check for specific error types and set appropriate status codes
		// if err == someSpecificError {
		//     code = fiber.StatusBadRequest
		// }
		
		// Return JSON error response
		return c.Status(code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
}
`,
	"decorator": `package decorators

import (
	"github.com/gofiber/fiber/v2"
	"time"
)

// {{.Name}}Decorator decorates handlers for {{.Resource}}
func {{.Name}}Decorator(handler fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Pre-processing
		startTime := time.Now()
		
		// Call the original handler
		err := handler(c)
		
		// Post-processing
		duration := time.Since(startTime)
		// Log duration or do other operations
		c.Set("X-Response-Time", duration.String())
		
		return err
	}
}
`,
	"module-interface": `package modules

import (
	"github.com/gofiber/fiber/v2"
)

// Module represents a feature module in the application
type Module interface {
	Register(app *fiber.App)
	AddSubmodule(submodule Module)
	SetPrefix(prefix string)
}
`,
	"main": `package main

import (
	"fmt"
	"log"
	"os"
	"{{.ProjectImport}}/config"
	"{{.ProjectImport}}/modules"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg := config.Load()
	
	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      cfg.AppName,
		ErrorHandler: createErrorHandler(),
	})

	// Middleware
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(recover.New())
	app.Use(cors.New())

	// Health check route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   fmt.Sprintf("%v", time.Now()),
		})
	})

	// Register modules
	// Example:
	// userModule := modules.NewUserModule()
	// userProfileModule := modules.NewUserProfileModule("/users")
	// userModule.AddSubmodule(userProfileModule)
	// userModule.Register(app)

	// Start server
	port := fmt.Sprintf(":%d", cfg.AppPort)
	log.Printf("Starting server on port %s", port)
	log.Fatal(app.Listen(port))
}

func createErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Default status code
		code := fiber.StatusInternalServerError
		
		// Check if it's a Fiber error
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}
		
		// Return JSON error response
		return c.Status(code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
}
`,
	"config": `package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	AppName         string
	AppEnv          string
	AppPort         int
	DatabaseURL     string
	DatabaseName    string
	JWTSecret       string
	LogLevel        string
}

// Load loads configuration from environment variables
func Load() *Config {
	// Default values
	config := &Config{
		AppName:      getEnv("APP_NAME", "FiberApp"),
		AppEnv:       getEnv("APP_ENV", "development"),
		AppPort:      getEnvAsInt("APP_PORT", 3000),
		DatabaseURL:  getEnv("DATABASE_URL", "mongodb://localhost:27017"),
		DatabaseName: getEnv("DATABASE_NAME", "fiberapp"),
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
	}

	return config
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
`,
	"dockerfile": `FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Use a small alpine image
FROM alpine:latest

# Install CA certificates
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/app .

# Expose the application port
EXPOSE 3000

# Command to run the application
CMD ["./app"]
`,
	"makefile": `# Go parameters
BINARY_NAME=app
MAIN_PATH=main.go
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

# Build parameters
BUILD_DIR=./bin
BUILD_FLAGS=-ldflags="-s -w"

# Docker parameters
DOCKER_IMG={{.ProjectName}}
DOCKER_TAG=latest

.PHONY: all build clean run test lint vet fmt docker-build docker-run

all: build

build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	@go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete."

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "Clean complete."

run:
	@go run $(MAIN_PATH)

test:
	@echo "Running tests..."
	@go test ./... -v

lint:
	@echo "Running linter..."
	@golint ./...

vet:
	@echo "Running vet..."
	@go vet ./...

fmt:
	@echo "Formatting code..."
	@gofmt -s -w $(GO_FILES)

docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMG):$(DOCKER_TAG) .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 3000:3000 $(DOCKER_IMG):$(DOCKER_TAG)

dev:
	@echo "Running in development mode with live reload..."
	@air
`,
}

// Project templates
var projectTemplates = map[string]string{
	"rest-api":     "REST API project with CRUD operations",
	"graphql":      "GraphQL API project",
	"microservice": "Microservice architecture project",
	"websocket":    "WebSocket server project",
}

// Database options
var databaseOptions = map[string]string{
	"mongodb":  "MongoDB",
	"postgres": "PostgreSQL",
	"mysql":    "MySQL",
	"sqlite":   "SQLite",
	"redis":    "Redis",
}

type templateData struct {
	Name          string
	Resource      string
	ResourceURL   string
	ProjectImport string
	ModulePath    string
	ModulePackage string
	ParentModule  string
	ProjectName   string
	PackageName   string
}

func main() {
	// Create root command
	rootCmd := &cobra.Command{
		Use:     "fiber-cli",
		Version: "1.0.0",
		Short:   "A CLI for Go Fiber framework, inspired by NestJS CLI",
		Long:    `A command line interface for Go Fiber framework that helps you to create and manage your Fiber projects with support for modules and submodules.`,
	}

	// Define colors for better output
	successColor := color.New(color.FgGreen).SprintFunc()
	infoColor := color.New(color.FgCyan).SprintFunc()
	errorColor := color.New(color.FgRed).SprintFunc()

	// new command - creates a new Fiber project with improved templates
	var newCmd = &cobra.Command{
		Use:   "new [project-name]",
		Short: "Create a new Go Fiber project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			projectName := args[0]

			// Get template type
			templateType, _ := cmd.Flags().GetString("template")

			// Get database type
			dbType, _ := cmd.Flags().GetString("database")

			// Show spinner during creation
			s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
			s.Prefix = "Creating project: "
			s.Start()

			err := createProject(projectName, templateType, dbType)
			s.Stop()

			if err != nil {
				fmt.Println(errorColor("Error: " + err.Error()))
				return
			}

			fmt.Println(successColor("\n✅ Project created successfully! 🚀"))
			fmt.Printf(infoColor("\nTo get started:\n"))
			fmt.Printf("  cd %s\n", projectName)
			fmt.Printf("  go mod tidy\n")
			fmt.Printf("  go run main.go\n")
			fmt.Printf(infoColor("\nCreate your first module:\n"))
			fmt.Printf("  fiber-cli generate module users\n")
			fmt.Printf(infoColor("\nCreate a submodule:\n"))
			fmt.Printf("  fiber-cli generate module profiles --parent /users\n")
		},
	}

	// Add flags for new command
	newCmd.Flags().StringP("template", "t", "rest-api", "Project template (rest-api, graphql, microservice, websocket)")
	newCmd.Flags().StringP("database", "d", "mongodb", "Database type (mongodb, postgres, mysql, sqlite, redis)")

	// generate command - generates components with more options
	var generateCmd = &cobra.Command{
		Use:     "generate [component-type] [component-name]",
		Aliases: []string{"g"},
		Short:   "Generate a component",
		Args:    cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			componentType := args[0]
			componentName := args[1]

			// Get module path from flags
			modulePath, _ := cmd.Flags().GetString("module")
			parentModule, _ := cmd.Flags().GetString("parent")
			flat, _ := cmd.Flags().GetBool("flat")
			skipTests, _ := cmd.Flags().GetBool("skip-tests")

			// Show spinner during generation
			s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
			s.Prefix = "Generating component: "
			s.Start()

			err := generateComponent(componentType, componentName, modulePath, parentModule, flat, skipTests)
			s.Stop()

			if err != nil {
				fmt.Println(errorColor("Error: " + err.Error()))
				return
			}

			fmt.Println(successColor(fmt.Sprintf("✅ %s generated successfully!", componentType)))
		},
	}

	// module command - create modules and submodules with better UX
	var moduleCmd = &cobra.Command{
		Use:   "module [module-name]",
		Short: "Generate a new module",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			moduleName := args[0]
			parentModule, _ := cmd.Flags().GetString("parent")
			skipTests, _ := cmd.Flags().GetBool("skip-tests")

			// Show spinner during creation
			s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
			s.Prefix = "Creating module: "
			s.Start()

			err := createModule(moduleName, parentModule, skipTests)
			s.Stop()

			if err != nil {
				fmt.Println(errorColor("Error: " + err.Error()))
				return
			}

			fmt.Println(successColor("✅ Module created successfully!"))

			// Print usage instructions
			if parentModule != "" {
				fmt.Printf(infoColor("\nTo use this submodule, add it to your parent module:\n"))
				fmt.Printf("\n  parentModule := modules.New%sModule()\n", toCamelCase(strings.TrimPrefix(parentModule, "/")))
				fmt.Printf("  submodule := modules.New%sModule(\"%s\")\n", toCamelCase(moduleName), parentModule)
				fmt.Printf("  parentModule.AddSubmodule(submodule)\n")
			} else {
				fmt.Printf(infoColor("\nTo use this module, add it to your main.go:\n"))
				fmt.Printf("\n  %sModule := %s.New%sModule()\n", toSnakeCase(moduleName), moduleName, toCamelCase(moduleName))
				fmt.Printf("  %sModule.Register(app)\n", toSnakeCase(moduleName))
			}
		},
	}

	// Add flags for module command
	moduleCmd.Flags().StringP("parent", "p", "", "Parent module path (for submodules)")
	moduleCmd.Flags().BoolP("skip-tests", "s", false, "Skip test file generation")

	// Add build command
	var buildCmd = &cobra.Command{
		Use:   "build",
		Short: "Build the application",
		Run: func(cmd *cobra.Command, args []string) {
			prod, _ := cmd.Flags().GetBool("prod")

			fmt.Println(infoColor("Building application..."))

			var buildCmd *exec.Cmd
			if prod {
				buildCmd = exec.Command("go", "build", "-ldflags", "-s -w", "-o", "app")
			} else {
				buildCmd = exec.Command("go", "build", "-o", "app")
			}

			output, err := buildCmd.CombinedOutput()
			if err != nil {
				fmt.Println(errorColor("Build failed:"))
				fmt.Println(string(output))
				return
			}

			fmt.Println(successColor("✅ Build completed successfully!"))
		},
	}

	// Add flags for build command
	buildCmd.Flags().BoolP("prod", "p", false, "Build for production with optimizations")

	// Add start command
	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the application",
		Run: func(cmd *cobra.Command, args []string) {
			watch, _ := cmd.Flags().GetBool("watch")

			if watch {
				fmt.Println(infoColor("Starting application in watch mode..."))
				watchCmd := exec.Command("air")
				watchCmd.Stdout = os.Stdout
				watchCmd.Stderr = os.Stderr
				err := watchCmd.Run()
				if err != nil {
					fmt.Println(errorColor("Failed to start application in watch mode:"))
					fmt.Println("Make sure 'air' is installed. Run: go install github.com/cosmtrek/air@latest")
					return
				}
			} else {
				fmt.Println(infoColor("Starting application..."))
				runCmd := exec.Command("go", "run", "main.go")
				runCmd.Stdout = os.Stdout
				runCmd.Stderr = os.Stderr
				err := runCmd.Run()
				if err != nil {
					fmt.Println(errorColor("Failed to start application:"))
					fmt.Println(err)
					return
				}
			}
		},
	}

	// Add flags for start command
	startCmd.Flags().BoolP("watch", "w", false, "Watch for changes and reload (requires 'air')")

	// Add test command
	var testCmd = &cobra.Command{
		Use:   "test",
		Short: "Run tests",
		Run: func(cmd *cobra.Command, args []string) {
			verbose, _ := cmd.Flags().GetBool("verbose")
			coverage, _ := cmd.Flags().GetBool("coverage")

			fmt.Println(infoColor("Running tests..."))

			testArgs := []string{"test", "./..."}
			if verbose {
				testArgs = append(testArgs, "-v")
			}
			if coverage {
				testArgs = append(testArgs, "-cover")
			}

			testCmd := exec.Command("go", testArgs...)
			testCmd.Stdout = os.Stdout
			testCmd.Stderr = os.Stderr
			err := testCmd.Run()
			if err != nil {
				fmt.Println(errorColor("Tests failed"))
				return
			}

			fmt.Println(successColor("✅ All tests passed!"))
		},
	}

	// Add flags for test command
	testCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	testCmd.Flags().BoolP("coverage", "c", false, "Show coverage")

	// Add add command for adding libraries
	var addCmd = &cobra.Command{
		Use:   "add [packages...]",
		Short: "Add packages to the project",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dev, _ := cmd.Flags().GetBool("dev")

			fmt.Println(infoColor("Adding packages: " + strings.Join(args, ", ")))

			var getCmd *exec.Cmd
			if dev {
				// Add as development dependency
				arguments := append([]string{"get", "-t"}, args...)
				getCmd = exec.Command("go", arguments...)
			} else {
				arguments := append([]string{"get"}, args...)
				getCmd = exec.Command("go", arguments...)
			}

			output, err := getCmd.CombinedOutput()
			if err != nil {
				fmt.Println(errorColor("Failed to add packages:"))
				fmt.Println(string(output))
				return
			}

			// Update fiber-cli.yaml
			updateCliConfig(args, dev)

			fmt.Println(successColor("✅ Packages added successfully!"))
		},
	}

	// Add flags for add command
	addCmd.Flags().BoolP("dev", "d", false, "Add as development dependency")

	// Add info command
	var infoCmd = &cobra.Command{
		Use:   "info",
		Short: "Display project information",
		Run: func(cmd *cobra.Command, args []string) {
			// Get project info
			projectInfo, err := getProjectInfo()
			if err != nil {
				fmt.Println(errorColor("Error getting project information:"))
				fmt.Println(err)
				return
			}

			fmt.Println(infoColor("Project Information:"))
			fmt.Printf("  Project Name: %s\n", projectInfo.ProjectName)
			fmt.Printf("  Project Type: %s\n", projectInfo.ProjectType)
			fmt.Printf("  Source Root: %s\n", projectInfo.SourceRoot)

			fmt.Println(infoColor("\nDependencies:"))
			fmt.Println("  Core:")
			for _, dep := range projectInfo.Dependencies.Core {
				fmt.Printf("    - %s\n", dep)
			}

			fmt.Println("  Dev:")
			for _, dep := range projectInfo.Dependencies.Dev {
				fmt.Printf("    - %s\n", dep)
			}

			fmt.Println(infoColor("\nModules:"))
			modules, err := getModules()
			if err != nil {
				fmt.Println("  No modules found")
			} else {
				for _, module := range modules {
					fmt.Printf("  - %s\n", module)
				}
			}
		},
	}

	// Add all commands to root command
	rootCmd.AddCommand(newCmd, generateCmd, moduleCmd, buildCmd, startCmd, testCmd, addCmd, infoCmd)

	// Execute the CLI
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(errorColor("Error:"), err)
		os.Exit(1)
	}
}

// createProject creates a new Fiber project with enhanced options
func createProject(projectName, templateType, dbType string) error {
	fmt.Printf("Creating a new Go Fiber project: %s\n", projectName)
	fmt.Printf("Template: %s, Database: %s\n", templateType, dbType)

	// Get GitHub username
	username, err := getGitUsername()
	if err != nil {
		return fmt.Errorf("failed to get GitHub username: %w", err)
	}

	// Format the full project import path
	fullProjectPath := fmt.Sprintf("github.com/%s/%s", username, projectName)

	// Create project directory
	if err := os.MkdirAll(projectName, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Create basic directory structure
	dirs := []string{
		"controllers",
		"services",
		"models",
		"repositories",
		"config",
		"modules",
		"middleware",
		"utils",
		"tests",
		"internal",
		"docs",
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(projectName, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create module interface
	moduleInterfaceFile, err := os.Create(filepath.Join(projectName, "modules", "module.go"))
	if err != nil {
		return fmt.Errorf("failed to create module interface file: %w", err)
	}
	defer moduleInterfaceFile.Close()

	// Parse and execute the template
	tmpl, err := template.New("module-interface").Parse(templates["module-interface"])
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	data := templateData{
		Name:          projectName,
		Resource:      strings.ToLower(projectName),
		ResourceURL:   strings.ToLower(projectName),
		ProjectImport: fullProjectPath,
		ProjectName:   projectName,
	}

	if err := tmpl.Execute(moduleInterfaceFile, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Create main.go
	mainFile, err := os.Create(filepath.Join(projectName, "main.go"))
	if err != nil {
		return fmt.Errorf("failed to create main.go: %w", err)
	}
	defer mainFile.Close()

	// Parse and execute the template
	tmpl, err = template.New("main").Parse(templates["main"])
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	if err := tmpl.Execute(mainFile, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Create config.go
	configFile, err := os.Create(filepath.Join(projectName, "config", "config.go"))
	if err != nil {
		return fmt.Errorf("failed to create config.go: %w", err)
	}
	defer configFile.Close()

	// Parse and execute the template
	tmpl, err = template.New("config").Parse(templates["config"])
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	if err := tmpl.Execute(configFile, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Create Dockerfile
	dockerFile, err := os.Create(filepath.Join(projectName, "Dockerfile"))
	if err != nil {
		return fmt.Errorf("failed to create Dockerfile: %w", err)
	}
	defer dockerFile.Close()

	// Parse and execute the template
	tmpl, err = template.New("dockerfile").Parse(templates["dockerfile"])
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	if err := tmpl.Execute(dockerFile, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Create Makefile
	makeFile, err := os.Create(filepath.Join(projectName, "Makefile"))
	if err != nil {
		return fmt.Errorf("failed to create Makefile: %w", err)
	}
	defer makeFile.Close()

	// Parse and execute the template
	tmpl, err = template.New("makefile").Parse(templates["makefile"])
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	if err := tmpl.Execute(makeFile, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Create go.mod
	goModContent := fmt.Sprintf(`module %s
	
	go 1.21
	
	require (
		github.com/gofiber/fiber/v2 v2.52.0
		github.com/AlecAivazis/survey/v2 v2.3.7
		github.com/spf13/cobra v1.8.0
		github.com/stretchr/testify v1.8.4
		github.com/briandowns/spinner v1.23.0
		github.com/fatih/color v1.16.0
		gopkg.in/yaml.v3 v3.0.1
	)
	`, fullProjectPath)

	goModFile, err := os.Create(filepath.Join(projectName, "go.mod"))
	if err != nil {
		return fmt.Errorf("failed to create go.mod: %w", err)
	}
	defer goModFile.Close()

	if _, err := goModFile.WriteString(goModContent); err != nil {
		return fmt.Errorf("failed to write to go.mod: %w", err)
	}

	// Create air.toml for hot reload
	airContent := `root = "."
	tmp_dir = "tmp"
	
	[build]
	  bin = "./tmp/main"
	  cmd = "go build -o ./tmp/main ."
	  delay = 1000
	  exclude_dir = ["assets", "tmp", "vendor"]
	  exclude_file = []
	  exclude_regex = [".*_test.go"]
	  exclude_unchanged = false
	  follow_symlink = false
	  full_bin = ""
	  include_dir = []
	  include_ext = ["go", "tpl", "tmpl", "html"]
	  kill_delay = "0s"
	  log = "build-errors.log"
	  send_interrupt = false
	  stop_on_error = true
	
	[color]
	  app = ""
	  build = "yellow"
	  main = "magenta"
	  runner = "green"
	  watcher = "cyan"
	
	[log]
	  time = false
	
	[misc]
	  clean_on_exit = false
	`

	airFile, err := os.Create(filepath.Join(projectName, ".air.toml"))
	if err != nil {
		return fmt.Errorf("failed to create .air.toml: %w", err)
	}
	defer airFile.Close()

	if _, err := airFile.WriteString(airContent); err != nil {
		return fmt.Errorf("failed to write to .air.toml: %w", err)
	}

	// Create .gitignore
	gitignoreContent := `# Binaries for programs and plugins
	*.exe
	*.exe~
	*.dll
	*.so
	*.dylib
	bin/
	tmp/
	
	# Test binary, built with 'go test -c'
	*.test
	
	# Output of the go coverage tool, specifically when used with LiteIDE
	*.out
	
	# Dependency directories (remove the comment below to include it)
	vendor/
	
	# IDE specific files
	.idea/
	.vscode/
	*.swp
	*.swo
	
	# Environment variables
	.env
	.env.local
	.env.development
	.env.test
	.env.production
	
	# Log files
	*.log
	
	# OS specific files
	.DS_Store
	`

	gitignoreFile, err := os.Create(filepath.Join(projectName, ".gitignore"))
	if err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}
	defer gitignoreFile.Close()

	if _, err := gitignoreFile.WriteString(gitignoreContent); err != nil {
		return fmt.Errorf("failed to write to .gitignore: %w", err)
	}

	// Create README.md
	readmeContent := fmt.Sprintf(`# %s
	
	A Fiber project generated using fiber-cli.
	
	## Getting Started
	
	### Prerequisites
	
	- Go 1.21 or later
	- MongoDB (or your selected database)
	
	### Installation
	
	1. Clone the repository
	   `+"```"+`sh
	   git clone https://github.com/%s/%s.git
	   `+"```"+`
	
	2. Install dependencies
	   `+"```"+`sh
	   cd %s
	   go mod tidy
	   `+"```"+`
	
	3. Run the application
	   `+"```"+`sh
	   # Development mode with hot reload
	   go run github.com/cosmtrek/air
	   
	   # Or standard mode
	   go run main.go
	   `+"```"+`
## Project Structure

`+"```"+`
.
├── config/             # Application configuration
├── controllers/        # Request handlers
├── middleware/         # Middleware components
├── models/             # Data models
├── modules/            # Feature modules
├── repositories/       # Data access layer
├── services/           # Business logic
├── utils/              # Utility functions
├── tests/              # Test files
├── .air.toml           # Hot reload configuration
├── .gitignore          # Git ignore file
├── Dockerfile          # Docker configuration
├── Makefile            # Make commands
├── go.mod              # Go modules
├── go.sum              # Go modules checksum
├── main.go             # Application entry point
└── README.md           # Project documentation
`+"```"+`

## Available Commands

- Build the application
  `+"```"+`sh
  make build
  `+"```"+`

- Run tests
  `+"```"+`sh
  make test
  `+"```"+`

- Run in development mode
  `+"```"+`sh
  make dev
  `+"```"+`

- Build and run Docker container
  `+"```"+`sh
  make docker-build
  make docker-run
  `+"```"+`

## Using fiber-cli

Create a new module:
`+"```"+`sh
fiber-cli generate module users
`+"```"+`

Generate a controller:
`+"```"+`sh
fiber-cli generate controller users
`+"```"+`

## License

This project is licensed under the MIT License - see the LICENSE file for details.
`, projectName, username, projectName, projectName)

	readmeFile, err := os.Create(filepath.Join(projectName, "README.md"))
	if err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}
	defer readmeFile.Close()

	if _, err := readmeFile.WriteString(readmeContent); err != nil {
		return fmt.Errorf("failed to write to README.md: %w", err)
	}

	// Create fiber-cli.yaml configuration file
	cliConfig := CommandData{
		Collection:  "fiber-standard",
		SourceRoot:  "./",
		ProjectName: projectName,
		ProjectType: templateType,
		ProjectRoot: "./",
		Monorepo:    false,
		Dependencies: struct {
			Core    []string `yaml:"core"`
			Dev     []string `yaml:"dev"`
			Test    []string `yaml:"test"`
			Scripts []string `yaml:"scripts"`
		}{
			Core: []string{
				"github.com/gofiber/fiber/v2",
				"github.com/AlecAivazis/survey/v2",
				"github.com/spf13/cobra",
			},
			Dev: []string{
				"github.com/cosmtrek/air",
			},
			Test: []string{
				"github.com/stretchr/testify",
			},
			Scripts: []string{
				"build: go build -o bin/app main.go",
				"start: go run main.go",
				"test: go test ./...",
				"dev: air",
			},
		},
	}

	cliConfigData, err := yaml.Marshal(cliConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal CLI config: %w", err)
	}

	cliConfigFile, err := os.Create(filepath.Join(projectName, "fiber-cli.yaml"))
	if err != nil {
		return fmt.Errorf("failed to create fiber-cli.yaml: %w", err)
	}
	defer cliConfigFile.Close()

	if _, err := cliConfigFile.Write(cliConfigData); err != nil {
		return fmt.Errorf("failed to write to fiber-cli.yaml: %w", err)
	}

	return nil
}

// createModule creates a module with improved organization
func createModule(moduleName, parentModule string, skipTests bool) error {
	// Convert module name to proper format
	name := toCamelCase(moduleName)
	resource := toSnakeCase(moduleName)
	resourceURL := toKebabCase(moduleName)

	// Get current directory (project root)
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Get project name from go.mod
	projectName, err := getProjectName(currentDir)
	if err != nil {
		return fmt.Errorf("failed to get project name: %w", err)
	}

	// Determine if this is a submodule
	isSubmodule := parentModule != ""

	// repository should be plural as repositories
	// Create module directory structure
	modulePath := filepath.Join("modules", resourceURL)
	if isSubmodule {
		// Remove leading slash if it exists
		if strings.HasPrefix(parentModule, "/") {
			parentModule = parentModule[1:]
		}
		modulePath = filepath.Join("modules", parentModule, resourceURL)
	}

	// Create directories
	dirs := []string{
		filepath.Join(modulePath, "controllers"),
		filepath.Join(modulePath, "services"),
		filepath.Join(modulePath, "repositories"),
		filepath.Join(modulePath, "dtos"),
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(currentDir, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create the module file
	moduleFilePath := filepath.Join(currentDir, modulePath, fmt.Sprintf("%s_module.go", resourceURL))

	moduleFile, err := os.Create(moduleFilePath)
	if err != nil {
		return fmt.Errorf("failed to create module file: %w", err)
	}
	defer moduleFile.Close()

	// Create template data
	data := templateData{
		Name:          name,
		Resource:      resource,
		ResourceURL:   resourceURL,
		ProjectImport: projectName,
		ModulePath:    modulePath,
		ModulePackage: filepath.Base(modulePath),
		ParentModule:  parentModule,
		ProjectName:   projectName,
		PackageName:   resourceURL,
	}

	// Use appropriate template based on whether this is a submodule
	templateKey := "module"
	if isSubmodule {
		templateKey = "submodule"
	}

	// Parse and execute the template
	tmpl, err := template.New(templateKey).Parse(templates[templateKey])
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	if err := tmpl.Execute(moduleFile, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Generate controller, service, repository, and DTO
	components := []string{"controller", "service", "repository", "dto"}
	for _, component := range components {

		if err := generateComponent(component, moduleName, filepath.Join(modulePath), parentModule, false, skipTests); err != nil {
			return fmt.Errorf("failed to generate %s: %w", component, err)
		}
	}

	// Generate tests if not skipped
	if !skipTests {
		if err := generateComponent("test", moduleName, filepath.Join(modulePath, "controllers"), parentModule, false, false); err != nil {
			return fmt.Errorf("failed to generate tests: %w", err)
		}
	}

	// Auto-register module in parent module if applicable
	if isSubmodule {
		// Find the parent module file
		parentModuleName := filepath.Base(parentModule)
		parentModulePath := filepath.Join(currentDir, "modules", parentModule, fmt.Sprintf("%s_module.go", parentModuleName))

		// Check if parent module exists
		if _, err := os.Stat(parentModulePath); err == nil {

			// TODO: Update parent module to include this submodule
			// This would require parsing Go code which is beyond the scope of this example
		}
	}

	return nil
}

// generateComponent generates a component with provided details
func generateComponent(componentType, componentName, modulePath, parentModule string, flat, skipTests bool) error {
	// Validate component type
	validTypes := map[string]bool{
		"controller": true,
		"service":    true,
		"repository": true,
		"model":      true,
		"dto":        true,
		"middleware": true,
		"guard":      true,
		"pipe":       true,
		"filter":     true,
		"decorator":  true,
		"test":       true,
	}

	if !validTypes[componentType] {
		return fmt.Errorf("invalid component type: %s", componentType)
	}

	// Convert component name to proper format
	name := toCamelCase(componentName)
	resource := toSnakeCase(componentName)
	resourceURL := toKebabCase(componentName)

	// Get current directory (project root)
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Get project name from go.mod
	projectName, err := getProjectName(currentDir)
	if err != nil {
		return fmt.Errorf("failed to get project name: %w", err)
	}

	// Determine file path and name
	componentDir := ""
	packageName := ""

	// For tests, use the module path provided
	if componentType == "test" {
		componentDir = modulePath
		packageName = filepath.Base(modulePath)
	} else if flat {
		// If flat structure, place directly in module path
		componentDir = modulePath
		packageName = filepath.Base(modulePath)
	} else if componentType == "repository" {
		// repository should be plural as repositories
		componentDir = filepath.Join(modulePath, "repositories")
		packageName = "repositories"
	} else {
		// Otherwise use standard structure with component type as directory
		componentDir = filepath.Join(modulePath, componentType+"s")
		packageName = componentType + "s"
	}
	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Join(currentDir, componentDir), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate file name
	fileName := ""
	if componentType == "test" {
		fileName = fmt.Sprintf("%s_test.go", resourceURL)
	} else {
		fileName = fmt.Sprintf("%s_%s.go", resourceURL, componentType)
	}

	filePath := filepath.Join(currentDir, componentDir, fileName)

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Create template data
	data := templateData{
		Name:          name,
		Resource:      resource,
		ResourceURL:   resourceURL,
		ProjectImport: projectName,
		ModulePath:    modulePath,
		ModulePackage: filepath.Base(modulePath),
		ParentModule:  parentModule,
		ProjectName:   projectName,
		PackageName:   packageName,
	}

	// Parse and execute the template
	tmpl, err := template.New(componentType).Parse(templates[componentType])
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Generate test file if appropriate and not skipped
	if componentType == "controller" && !skipTests {
		if err := generateComponent("test", componentName, filepath.Join(modulePath, "controllers"), parentModule, false, true); err != nil {
			return fmt.Errorf("failed to generate test: %w", err)
		}
	}

	return nil
}

// getProjectInfo gets project information from fiber-cli.yaml
func getProjectInfo() (*CommandData, error) {
	// Check if fiber-cli.yaml exists
	if _, err := os.Stat("fiber-cli.yaml"); os.IsNotExist(err) {
		return nil, fmt.Errorf("fiber-cli.yaml not found, are you in a Fiber CLI project?")
	}

	// Read and parse fiber-cli.yaml
	data, err := os.ReadFile("fiber-cli.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read fiber-cli.yaml: %w", err)
	}

	var config CommandData
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse fiber-cli.yaml: %w", err)
	}

	return &config, nil
}

// updateCliConfig updates fiber-cli.yaml with new dependencies
func updateCliConfig(packages []string, isDev bool) error {
	config, err := getProjectInfo()
	if err != nil {
		return err
	}

	// Add packages to appropriate section
	if isDev {
		config.Dependencies.Dev = append(config.Dependencies.Dev, packages...)
	} else {
		config.Dependencies.Core = append(config.Dependencies.Core, packages...)
	}

	// Write updated config
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile("fiber-cli.yaml", data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// getModules returns a list of available modules
func getModules() ([]string, error) {
	modules := []string{}

	// Check if modules directory exists
	modulesDir := "modules"
	if _, err := os.Stat(modulesDir); os.IsNotExist(err) {
		return modules, nil
	}

	// Walk through modules directory
	err := filepath.Walk(modulesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Look for *_module.go files
		if !info.IsDir() && strings.HasSuffix(info.Name(), "_module.go") {
			// Extract module name
			moduleName := strings.TrimSuffix(info.Name(), "_module.go")

			// Get relative path from modules directory
			relPath, err := filepath.Rel(modulesDir, filepath.Dir(path))
			if err != nil {
				return err
			}

			// If it's in the root modules directory
			if relPath == "." {
				modules = append(modules, moduleName)
			} else {
				// For nested modules, format as parent/module
				modules = append(modules, filepath.Join(relPath, moduleName))
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list modules: %w", err)
	}

	return modules, nil
}

// getProjectName extracts project name from go.mod
func getProjectName(projectDir string) (string, error) {
	goModPath := filepath.Join(projectDir, "go.mod")

	// Check if go.mod exists
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return "", fmt.Errorf("go.mod not found, are you in a Go project?")
	}

	// Read go.mod file
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}

	// Extract module name
	content := string(data)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}

	return "", fmt.Errorf("module name not found in go.mod")
}

// getGitUsername attempts to get the GitHub username
func getGitUsername() (string, error) {
	// First, try to get from git config
	cmd := exec.Command("git", "config", "--get", "user.username")
	output, err := cmd.Output()
	if err == nil && len(output) > 0 {
		return strings.TrimSpace(string(output)), nil
	}

	// If that fails, prompt the user
	username := ""
	prompt := &survey.Input{
		Message: "Enter your GitHub username:",
	}
	err = survey.AskOne(prompt, &username)
	if err != nil {
		return "", err
	}

	if username == "" {
		return "user", nil // Default fallback
	}

	return username, nil
}

// Helper functions for string formatting
func toCamelCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '-' || r == '_' || r == ' '
	})

	result := ""
	for _, word := range words {
		if len(word) == 0 {
			continue
		}
		result += strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
	}

	return result
}

func toSnakeCase(s string) string {
	// First convert to kebab case
	kebab := toKebabCase(s)
	// Then replace hyphens with underscores
	return strings.ReplaceAll(kebab, "-", "_")
}

func toKebabCase(s string) string {
	// Add space before uppercase letters
	s = regexp.MustCompile(`([a-z0-9])([A-Z])`).ReplaceAllString(s, "$1 $2")

	// Convert to lowercase and replace spaces with hyphens
	return strings.ToLower(strings.ReplaceAll(s, " ", "-"))
}
