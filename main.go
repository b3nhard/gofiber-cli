package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var templates = map[string]string{
	"controller": `package controllers

import (
	"github.com/gofiber/fiber/v2"
)

// {{.Name}}Controller handles {{.Resource}} related requests
type {{.Name}}Controller struct {
	// Add dependencies here
}

// New{{.Name}}Controller creates a new {{.Name}}Controller
func New{{.Name}}Controller() *{{.Name}}Controller {
	return &{{.Name}}Controller{}
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
	return ctx.JSON(fiber.Map{
		"message": "Get all {{.Resource}}",
	})
}

// GetOne returns one {{.Resource}} by ID
func (c *{{.Name}}Controller) GetOne(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	return ctx.JSON(fiber.Map{
		"message": "Get {{.Resource}} with ID: " + id,
	})
}

// Create creates a new {{.Resource}}
func (c *{{.Name}}Controller) Create(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"message": "Create {{.Resource}}",
	})
}

// Update updates a {{.Resource}} by ID
func (c *{{.Name}}Controller) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	return ctx.JSON(fiber.Map{
		"message": "Update {{.Resource}} with ID: " + id,
	})
}

// Delete deletes a {{.Resource}} by ID
func (c *{{.Name}}Controller) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	return ctx.JSON(fiber.Map{
		"message": "Delete {{.Resource}} with ID: " + id,
	})
}
`,
	"service": `package services

// {{.Name}}Service provides functionality for {{.Resource}}
type {{.Name}}Service struct {
	// Add dependencies here
}

// New{{.Name}}Service creates a new {{.Name}}Service
func New{{.Name}}Service() *{{.Name}}Service {
	return &{{.Name}}Service{}
}

// GetAll returns all {{.Resource}}
func (s *{{.Name}}Service) GetAll() ([]interface{}, error) {
	// Implement this
	return []interface{}{}, nil
}

// GetByID returns a {{.Resource}} by ID
func (s *{{.Name}}Service) GetByID(id string) (interface{}, error) {
	// Implement this
	return nil, nil
}

// Create creates a new {{.Resource}}
func (s *{{.Name}}Service) Create(data interface{}) (interface{}, error) {
	// Implement this
	return nil, nil
}

// Update updates a {{.Resource}} by ID
func (s *{{.Name}}Service) Update(id string, data interface{}) (interface{}, error) {
	// Implement this
	return nil, nil
}

// Delete deletes a {{.Resource}} by ID
func (s *{{.Name}}Service) Delete(id string) error {
	// Implement this
	return nil
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
	"{{.ProjectImport}}/models"
)

// {{.Name}}Repository handles database operations for {{.Resource}}
type {{.Name}}Repository struct {
	// Add database connection or ORM here
}

// New{{.Name}}Repository creates a new {{.Name}}Repository
func New{{.Name}}Repository() *{{.Name}}Repository {
	return &{{.Name}}Repository{}
}

// FindAll returns all {{.Resource}}
func (r *{{.Name}}Repository) FindAll() ([]models.{{.Name}}, error) {
	// Implement this
	return []models.{{.Name}}{}, nil
}

// FindByID finds a {{.Resource}} by ID
func (r *{{.Name}}Repository) FindByID(id string) (*models.{{.Name}}, error) {
	// Implement this
	return nil, nil
}

// Create creates a new {{.Resource}}
func (r *{{.Name}}Repository) Create({{.Resource}} *models.{{.Name}}) error {
	// Implement this
	return nil
}

// Update updates a {{.Resource}} by ID
func (r *{{.Name}}Repository) Update({{.Resource}} *models.{{.Name}}) error {
	// Implement this
	return nil
}

// Delete deletes a {{.Resource}} by ID
func (r *{{.Name}}Repository) Delete(id string) error {
	// Implement this
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
	service := services.New{{.Name}}Service()
	controller := controllers.New{{.Name}}Controller()

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
	service := services.New{{.Name}}Service()
	controller := controllers.New{{.Name}}Controller()

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
	"log"
	"{{.ProjectImport}}/config"
	"{{.ProjectImport}}/modules"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())

	// Register modules
	// Example:
	// userModule := modules.NewUserModule()
	// userProfileModule := modules.NewUserProfileModule("/users")
	// userModule.AddSubmodule(userProfileModule)
	// userModule.Register(app)

	// Start server
	log.Fatal(app.Listen(":3000"))
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
}

type templateData struct {
	Name          string
	Resource      string
	ResourceURL   string
	ProjectImport string
	ModulePath    string
	ModulePackage string
	ParentModule  string
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "fiber-cli",
		Short: "A CLI for Go Fiber framework, inspired by NestJS CLI",
		Long:  `A command line interface for Go Fiber framework that helps you to create and manage your Fiber projects with support for modules and submodules.`,
	}

	// new command - creates a new Fiber project
	var newCmd = &cobra.Command{
		Use:   "new [project-name]",
		Short: "Create a new Go Fiber project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			projectName := args[0]
			createProject(projectName)
		},
	}

	// generate command - generates components
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

			generateComponent(componentType, componentName, modulePath, parentModule)
		},
	}

	// Add flags for generate command
	generateCmd.Flags().StringP("module", "m", "", "Module path for the component (e.g., 'users' or 'users/profiles')")
	generateCmd.Flags().StringP("parent", "p", "", "Parent module for submodule generation (e.g., '/users')")

	// Create module command
	var moduleCmd = &cobra.Command{
		Use:   "module [module-name]",
		Short: "Generate a new module",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			moduleName := args[0]
			parentModule, _ := cmd.Flags().GetString("parent")
			createModule(moduleName, parentModule)
		},
	}

	// Add flags for module command
	moduleCmd.Flags().StringP("parent", "p", "", "Parent module path (for submodules)")

	rootCmd.AddCommand(newCmd, generateCmd, moduleCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createProject(projectName string) {
	fmt.Printf("Creating a new Go Fiber project: %s\n", projectName)

	// Get GitHub username
	username, err := getGitUsername()
	if err != nil {
		fmt.Printf("Error getting GitHub username: %v\n", err)
		return
	}

	// Format the full project import path
	fullProjectPath := fmt.Sprintf("github.com/%s/%s", username, projectName)
	fmt.Printf("Creating a new Go Fiber project: %s\n", fullProjectPath)

	// Create project directory
	if err := os.MkdirAll(projectName, 0755); err != nil {
		fmt.Printf("Error creating project directory: %v\n", err)
		return
	}

	// Create basic directory structure
	dirs := []string{
		"controllers",
		"services",
		"models",
		"repositories",
		"config",
		"modules",
		"middlewares",
		"utils",
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(projectName, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			return
		}
	}

	// Create module interface
	moduleInterfaceFile, err := os.Create(filepath.Join(projectName, "modules", "module.go"))
	if err != nil {
		fmt.Printf("Error creating module interface file: %v\n", err)
		return
	}
	defer moduleInterfaceFile.Close()

	// Parse and execute the template
	tmpl, err := template.New("module-interface").Parse(templates["module-interface"])
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		return
	}

	data := templateData{
		Name:          projectName,
		Resource:      strings.ToLower(projectName),
		ResourceURL:   strings.ToLower(projectName),
		ProjectImport: fmt.Sprintf("github.com/%s/%s", username, projectName),
	}

	if err := tmpl.Execute(moduleInterfaceFile, data); err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		return
	}

	// Create main.go
	mainFile, err := os.Create(filepath.Join(projectName, "main.go"))
	if err != nil {
		fmt.Printf("Error creating main.go: %v\n", err)
		return
	}
	defer mainFile.Close()

	// Parse and execute the template
	tmpl, err = template.New("main").Parse(templates["main"])
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		return
	}

	if err := tmpl.Execute(mainFile, data); err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		return
	}

	// Create config.go
	configFile, err := os.Create(filepath.Join(projectName, "config", "config.go"))
	if err != nil {
		fmt.Printf("Error creating config.go: %v\n", err)
		return
	}
	defer configFile.Close()

	// Parse and execute the template
	tmpl, err = template.New("config").Parse(templates["config"])
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		return
	}

	if err := tmpl.Execute(configFile, data); err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		return
	}

	// Create go.mod
	goModContent := fmt.Sprintf(`module %s

go 1.21

require (
	github.com/gofiber/fiber/v2 v2.52.0
	github.com/AlecAivazis/survey/v2 v2.3.7
	github.com/spf13/cobra v1.8.0
)
`, fullProjectPath)

	goModFile, err := os.Create(filepath.Join(projectName, "go.mod"))
	if err != nil {
		fmt.Printf("Error creating go.mod: %v\n", err)
		return
	}
	defer goModFile.Close()

	if _, err := goModFile.WriteString(goModContent); err != nil {
		fmt.Printf("Error writing to go.mod: %v\n", err)
		return
	}

	fmt.Printf("\nProject created successfully! 🚀\n")
	fmt.Printf("\nTo get started:\n")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Printf("  go mod tidy\n")
	fmt.Printf("  go run main.go\n")
	fmt.Printf("\nCreate your first module:\n")
	fmt.Printf("  fiber-cli module users\n")
	fmt.Printf("\nCreate a submodule:\n")
	fmt.Printf("  fiber-cli module profiles --parent /users\n")
}

func createModule(moduleName string, parentModule string) {
	// Convert module name to proper format
	name := toCamelCase(moduleName)
	resource := toSnakeCase(moduleName)
	resourceURL := toKebabCase(moduleName)

	// Get current directory (project root)
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	// Get project name from go.mod
	projectName, err := getProjectName(currentDir)
	if err != nil {
		fmt.Printf("Error getting project name: %v\n", err)
		return
	}

	// Determine if this is a submodule
	isSubmodule := parentModule != ""

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
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(currentDir, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			return
		}
	}

	// Create the module file
	moduleFilePath := filepath.Join(currentDir, modulePath, fmt.Sprintf("%s_module.go", resourceURL))

	moduleFile, err := os.Create(moduleFilePath)
	if err != nil {
		fmt.Printf("Error creating module file: %v\n", err)
		return
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
	}

	// Use appropriate template based on whether this is a submodule
	templateKey := "module"
	if isSubmodule {
		templateKey = "submodule"
	}

	// Parse and execute the template
	tmpl, err := template.New(templateKey).Parse(templates[templateKey])
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		return
	}

	if err := tmpl.Execute(moduleFile, data); err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		return
	}

	// Generate controller, service, and repository
	components := []string{"controller", "service", "repository"}
	for _, component := range components {
		generateComponent(component, moduleName, filepath.Join(modulePath), parentModule)
	}

	if isSubmodule {
		fmt.Printf("Submodule '%s' created successfully under parent '%s'!\n", moduleName, parentModule)
		fmt.Printf("\nTo use this submodule, add it to your parent module:\n")
		fmt.Printf("\n  parentModule := modules.New%sModule()\n", toCamelCase(strings.TrimPrefix(parentModule, "/")))
		fmt.Printf("  submodule := modules.New%sModule(\"%s\")\n", name, parentModule)
		fmt.Printf("  parentModule.AddSubmodule(submodule)\n")
	} else {
		fmt.Printf("Module '%s' created successfully!\n", moduleName)
		fmt.Printf("\nTo use this module, add it to your main.go:\n")
		fmt.Printf("\n  %sModule := modules.New%sModule()\n", resource, name)
		fmt.Printf("  %sModule.Register(app)\n", resource)
	}
}

func generateComponent(componentType, componentName, modulePath, parentModule string) {
	// Convert component name to proper format
	name := toCamelCase(componentName)
	resource := toSnakeCase(componentName)
	resourceURL := toKebabCase(componentName)

	// Check if the component type is supported
	templateContent, ok := templates[componentType]
	if !ok {
		fmt.Printf("Unsupported component type: %s\n", componentType)
		fmt.Println("Supported types: controller, service, model, repository, middleware")
		return
	}

	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	// Get project name from go.mod
	projectName, err := getProjectName(currentDir)
	if err != nil {
		// Not fatal, just use default
		projectName = "yourproject"
	}

	// Determine the directory to create the component in
	var dirPath string

	if modulePath != "" {
		// If module path is specified, use it with plural if repository
		if componentType == "repository" {
			dirPath = filepath.Join(currentDir, modulePath, "repositories")
		} else {
			dirPath = filepath.Join(currentDir, modulePath, componentType+"s")
		}
	} else {
		// Otherwise use the default directory
		dirPath = filepath.Join(currentDir, componentType+"s")
	}

	// Create the directory if it doesn't exist
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dirPath, err)
			return
		}
	}

	// Create the file path
	filePath := filepath.Join(dirPath, resourceURL+"_"+componentType+".go")

	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		var overwrite bool
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("File %s already exists. Overwrite?", filePath),
		}
		survey.AskOne(prompt, &overwrite)

		if !overwrite {
			fmt.Println("Operation cancelled")
			return
		}
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	// Get the module package name
	modulePackage := filepath.Base(modulePath)
	if modulePackage == "." || modulePackage == "" {
		modulePackage = componentType + "s"
	}

	// Parse and execute the template
	tmpl, err := template.New(componentType).Parse(templateContent)
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		return
	}

	data := templateData{
		Name:          name,
		Resource:      resource,
		ResourceURL:   resourceURL,
		ProjectImport: projectName,
		ModulePath:    modulePath,
		ModulePackage: modulePackage,
		ParentModule:  parentModule,
	}

	if err := tmpl.Execute(file, data); err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		return
	}

	fmt.Printf("%s generated successfully: %s\n", componentType, filePath)
}

// Helper functions for case conversion
func toCamelCase(s string) string {
	// Convert from snake_case or kebab-case to CamelCase
	s = strings.TrimSpace(s)
	s = strings.Replace(s, "-", "_", -1)
	parts := strings.Split(s, "_")

	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
		}
	}

	return strings.Join(parts, "")
}

func toSnakeCase(s string) string {
	// Convert from CamelCase or kebab-case to snake_case
	s = strings.TrimSpace(s)
	s = strings.Replace(s, "-", "_", -1)

	if strings.Contains(s, "_") {
		parts := strings.Split(s, "_")
		for i, part := range parts {
			if len(part) > 0 {
				parts[i] = strings.ToLower(part)
			}
		}
		return strings.Join(parts, "_")
	} else {
		var result strings.Builder
		for i, r := range s {
			if i > 0 && 'A' <= r && r <= 'Z' {
				result.WriteRune('_')
			}
			result.WriteRune(r)
		}
		return strings.ToLower(result.String())
	}
}

func toKebabCase(s string) string {
	// Convert from CamelCase or snake_case to kebab-case
	s = toSnakeCase(s)
	return strings.Replace(s, "_", "-", -1)
}

func getProjectName(dir string) (string, error) {
	// Try to read go.mod file to get the module name
	goModPath := filepath.Join(dir, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return "", err
	}

	// Extract module name
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module")), nil
		}
	}

	return "", fmt.Errorf("module name not found in go.mod")
}

func getGitUsername() (string, error) {
	// Try getting from git config
	cmd := exec.Command("git", "config", "--get", "user.username")
	output, err := cmd.Output()
	if err == nil && len(output) > 0 {
		return strings.TrimSpace(string(output)), nil
	}

	// If git config fails, try getting from environment variable
	username := os.Getenv("GITHUB_USERNAME")
	if username != "" {
		return username, nil
	}

	// If both methods fail, prompt the user
	var answer string
	prompt := &survey.Input{
		Message: "Enter your GitHub username:",
	}
	err = survey.AskOne(prompt, &answer)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(answer), nil
}
