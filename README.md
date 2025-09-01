# Swagger Module for Tinh Tinh

<div align="center">
<img alt="GitHub Release" src="https://img.shields.io/github/v/release/tinh-tinh/swagger">
<img alt="GitHub License" src="https://img.shields.io/github/license/tinh-tinh/swagger">
<a href="https://codecov.io/gh/tinh-tinh/swagger" > 
 <img src="https://codecov.io/gh/tinh-tinh/swagger/graph/badge.svg?token=WTLU1P7INE"/> 
 </a>
<a href="https://pkg.go.dev/github.com/tinh-tinh/swagger"><img src="https://pkg.go.dev/badge/github.com/tinh-tinh/swagger.svg" alt="Go Reference"></a>
</div>

<div align="center">
    <img src="https://avatars.githubusercontent.com/u/178628733?s=400&u=2a8230486a43595a03a6f9f204e54a0046ce0cc4&v=4" width="200" alt="Tinh Tinh Logo">
</div>

## Overview

The Swagger module for Tinh Tinh provides automatic OpenAPI 3.0 documentation and a built-in Swagger UI for your Tinh Tinh-based APIs. It introspects your controllers, DTOs, routes, and metadata to generate a standards-compliant OpenAPI spec.

## Features

- Generates OpenAPI 3.0 spec from your Tinh Tinh app
- Auto-discovers routes, query/body/path parameters, responses, and security
- Adds detailed schemas from your DTOs and validation tags
- Supports file upload, security schemes, tags, and custom consumers
- Ready-to-use Swagger UI, served directly from your app
- Easily extensible and customizable via metadata

## Installation

```bash
go get -u github.com/tinh-tinh/swagger/v2
```

## Quick Start Example

```go
import (
    "github.com/tinh-tinh/swagger/v2"
    "github.com/tinh-tinh/tinhtinh/v2/core"
)

server := core.CreateFactory(AppModule)
server.SetGlobalPrefix("/api")

// Build the OpenAPI spec
spec := swagger.NewSpecBuilder().
    SetTitle("My API").
    SetServer(&swagger.ServerObject{
        Url: "http://localhost:3000/api",
    }).
    SetDescription("API documentation for my service").
    SetVersion("1.0.0").
    AddSecurity(&swagger.SecuritySchemeObject{
        Type:         "http",
        Scheme:       "Bearer",
        BearerFormat: "JWT",
        Name:         "bearerAuth",
    }).
    Build()

// Parse all app routes and generate spec
spec.ParsePaths(server)

// Mount Swagger UI and OpenAPI JSON
swagger.SetUp("/swagger", server, spec)
```

- Swagger UI will be available at: `http://localhost:3000/swagger`
- OpenAPI JSON will be at: `http://localhost:3000/openapi.json`

## Usage Patterns

### Controller and DTO Example

```go
type SignUpUser struct {
    Name     string    `validate:"isAlpha" example:"John"`
    Email    string    `validate:"required,isEmail" example:"john@gmail.com"`
    Password string    `validate:"required,isStrongPassword" example:"12345678@Tc"`
    Birth    time.Time `validate:"required" example:"2024-12-12"`
    Avatar   string    `example:"https://cdn.tinhtinh.dev/images/avatar.png"`
}

func authController(module core.Module) core.Controller {
    return module.NewController("Auth").
        Metadata(swagger.ApiTag("Auth")).Registry().
        Pipe(core.Body(SignUpUser{})).
        Post("", func(ctx core.Ctx) error {
            payload := ctx.Body().(*SignUpUser)
            return ctx.JSON(core.Map{"data": payload})
        })
}
```

### Add Route Metadata

#### Tags
```go
swagger.ApiTag("User")
```

#### Security
```go
swagger.ApiSecurity("bearerAuth")
```

#### Consumers (e.g., multipart)
```go
swagger.ApiConsumer("multipart/form-data")
```

#### File Upload
```go
swagger.ApiFile(swagger.FileOptions{
    Name: "image",
    Description: "user avatar (image file)",
    Required: true,
})
```

Example controller for image upload:
```go
type UploadFile struct {
    Image storage.File `example:"image"`
}

func fileController(module core.Module) core.Controller {
    ctrl := module.NewController("Files").
        Metadata(swagger.ApiTag("File")).
        Metadata(swagger.ApiConsumer("multipart/form-data")).
        Metadata(swagger.ApiFile(swagger.FileOptions{
            Name: "image",
            Description: "User image file",
            Required: true,
        }))

    ctrl.Pipe(core.Body(UploadFile{})).Post("/upload", func(ctx core.Ctx) error {
        file := ctx.Body().(*UploadFile)
        // handle file.Image
        return ctx.JSON(core.Map{"url": "https://cdn.tinhtinh.dev/images/example-upload.png"})
    })

    return ctrl
}
```

#### Custom Response Schema
```go
swagger.ApiOkResponse(&Response{
    Title: "Acrane",
    Image: "https://cdn.tinhtinh.dev/images/sample-image.png",
})
```

## How It Works

- Parses all controllers and their routes for HTTP methods, paths, and DTOs
- Uses struct tags (`query`, `path`, `example`, `validate`, etc.) to generate detailed parameter and schema info
- Automatically creates OpenAPI-compliant docs with proper reference linking
- Renders a modern Swagger UI using CDN assets

## Customize and Extend

- Use your own metadata (tags, security, consumers, file options) on controllers/routes
- Add/override OpenAPI info in the `SpecBuilder`
- Mix multiple modules/controllers and all routes will be scanned

## Contributing

We welcome contributions! Please feel free to submit a Pull Request.

## Support

If you encounter any issues or need help, you can:
- Open an issue in the GitHub repository
- Check our documentation
- Join our community discussions
