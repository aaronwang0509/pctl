# PCTL - PAIC Control CLI Platform

## Project Overview

**PCTL** (PAIC Control) is a modern, comprehensive CLI platform built with Go for managing, testing, and automating Ping Identity Advanced Identity Cloud (PAIC) operations. It serves as the unified command-line interface that consolidates multiple testing and automation tools into a single, powerful, enterprise-ready CLI.

## Vision & Mission

PCTL aims to be the **definitive CLI platform for PAIC operations**, providing:

- **Unified Interface**: Single CLI for all PAIC testing, automation, and management tasks
- **Modern Architecture**: Built with Go for performance, reliability, and easy deployment
- **Extensible Design**: Plugin-based architecture supporting unlimited sub-tools
- **Enterprise Ready**: Battle-tested reliability for production environments
- **Developer First**: Excellent developer experience with comprehensive tooling

## Current Status

- **Version**: 0.1.0 (Active Development)
- **Architecture**: Go-based CLI platform with layered architecture
- **Build System**: Go modules with comprehensive testing (go test, coverage)
- **Token Generation**: ✅ **PRODUCTION READY** - Real PAIC OAuth 2.0 implementation
- **Testing Coverage**: 76.8% with unit, integration, and internal API tests
- **Current State**: Token functionality complete, ready for ELK integration
- **Next Goal**: Integrate ELK functionality with token-based PAIC authentication

## Dependencies & Integration

PCTL builds on top of **Frodo CLI** for core PAIC configuration management:
- Uses Frodo for PAIC connectivity and authentication
- Leverages Frodo's configuration system for environments
- Extends Frodo with additional testing and automation capabilities
- Maintains compatibility with existing Frodo workflows

## Project Architecture

### Core Design Principles

1. **Go-First Architecture**: Modern, performant CLI built with Go
2. **Plugin-Based Design**: Each tool is a plugin/sub-command
3. **Frodo Integration**: Seamless integration with existing Frodo workflows  
4. **Enterprise Reliability**: Production-grade error handling and logging
5. **Developer Experience**: Comprehensive help, validation, and debugging
6. **Single Binary**: Self-contained executable with no runtime dependencies

### Command Structure
```
pctl <command> [subcommand] [options]

Current Migration Targets:
  pctl journey <options>     # 🔄 MIGRATION TARGET - authflow journey → Go
  pctl token <options>       # 🔄 MIGRATION TARGET - authflow token → Go  
  pctl elk <action>          # 🔄 MIGRATION TARGET - plctl.sh → Go
  
Future Expansion:
  pctl frodo <action>        # 🎯 FUTURE - Direct Frodo integration
  pctl config <action>       # 🎯 FUTURE - Configuration management
  pctl env <action>          # 🎯 FUTURE - Environment management
```

### Architecture Design

**Framework**: **Cobra CLI** - Industry standard Go CLI framework
- Perfect for `pctl <tool> <action>` command structure  
- Built-in help, validation, and subcommand management
- Used by kubectl, Docker CLI, Hugo, and major Go CLIs

**Design Pattern**: **Layered Architecture with Clean Separation**
- **CLI Layer** (`cmd/`): Cobra commands handling user interface
- **Public API Layer** (`pkg/`): Reusable business logic for external consumption
- **Internal Layer** (`internal/`): Implementation details, services, and utilities

**Configuration Strategy**: **Config-First with Override Pattern**
- **Primary**: YAML configuration files for complex operations
- **Secondary**: CLI flags for overrides and quick modifications
- **Hierarchy**: CLI flags → Environment vars → Config files → Defaults
- **Tool**: Viper for configuration management with validation

### Planned Architecture
```
pctl/
├── main.go                  # Application entry point
├── cmd/                     # CLI Layer (Cobra commands)
│   ├── root.go             # Main CLI entry point and global config
│   ├── journey.go          # pctl journey - Journey testing CLI
│   ├── token.go            # pctl token - Token management CLI
│   └── elk.go              # pctl elk - ELK stack management CLI
├── pkg/                    # Public API Layer (importable)
│   ├── journey/            # Journey testing public APIs
│   │   ├── client.go       # Public client: NewClient(), Test()
│   │   ├── types.go        # Public types and interfaces
│   │   └── config.go       # Configuration structures
│   ├── token/              # Token management public APIs
│   │   ├── generator.go    # Public generator: NewGenerator(), Generate()
│   │   ├── types.go        # Public types and interfaces
│   │   └── config.go       # Configuration structures
│   ├── elk/                # ELK management public APIs  
│   │   ├── manager.go      # Public manager: Start(), Stop(), Logs()
│   │   ├── types.go        # Public types and interfaces
│   │   └── config.go       # Configuration structures
│   └── common/             # Shared public utilities
├── internal/               # Internal Implementation Layer
│   ├── journey/            # Journey testing business logic
│   │   ├── service.go      # Core authentication journey logic
│   │   ├── api.go          # PAIC API communication
│   │   └── validator.go    # Configuration validation
│   ├── token/              # Token management business logic
│   │   ├── generator.go    # JWT token generation logic
│   │   ├── service.go      # Service account token logic
│   │   └── validator.go    # Configuration validation
│   ├── elk/                # ELK management business logic
│   │   ├── docker.go       # Docker Compose orchestration
│   │   ├── logs.go         # Log streaming implementation  
│   │   └── monitoring.go   # Health checks and status
│   ├── config/             # Global configuration management
│   │   ├── loader.go       # Viper-based config loading with hierarchy
│   │   ├── validator.go    # Configuration validation
│   │   └── types.go        # Global configuration types
│   ├── logger/             # Structured logging system
│   └── utils/              # Internal utilities and helpers
├── configs/                # Configuration templates and examples
│   ├── journey/            # Journey YAML config examples
│   ├── token/              # Token YAML config examples
│   └── elk/                # ELK YAML config examples
└── go.mod                  # Go modules and dependencies
```

## Existing Tools Analysis

### 1. AuthFlow (TypeScript CLI)
**Current State**: Mature TypeScript CLI with comprehensive PAIC testing capabilities

**Key Features**:
- Authentication journey testing with YAML configs
- JWT token generation and management
- Service account token creation
- Cross-platform binary distribution
- Topic-based command structure (`authflow journey`, `authflow token`)

**Migration Priority**: HIGH - Core PAIC testing functionality
**Migration Complexity**: Medium - Well-structured codebase with clear interfaces

### 2. ELK Local (Shell/Python Management Tool)  
**Current State**: Comprehensive Docker-based ELK stack management

**Key Features**:
- Docker-composed Elasticsearch + Kibana setup
- Real-time PAIC log streaming via Frodo
- Automated log rotation and lifecycle management
- Platform detection (Linux/macOS with architecture support)
- Background daemon log streaming

**Migration Priority**: HIGH - Critical for log analysis workflows
**Migration Complexity**: Medium-High - Shell scripting + Python + Docker orchestration

## Build System & Distribution

### Target Build Configuration
```bash
# Cross-platform binary generation
go build -o bin/pctl-linux-amd64
go build -o bin/pctl-linux-arm64  
go build -o bin/pctl-darwin-amd64
go build -o bin/pctl-darwin-arm64
go build -o bin/pctl-windows-amd64.exe
```

### Dependencies
- **Go 1.21+**: Primary development language
- **Frodo CLI**: External dependency for PAIC operations
- **Docker**: For ELK stack management (elk tool)
- **No Runtime Dependencies**: Single binary distribution

## Migration Strategy

### Phase 1: Foundation ✅ **COMPLETED**
- [x] Project structure setup
- [x] CLAUDE.md documentation  
- [x] Go module initialization
- [x] Basic CLI framework with Cobra
- [x] Core package structure

### Phase 2: Token Generation ✅ **COMPLETED**
- [x] Real JWT token generation with PAIC OAuth 2.0 flows
- [x] JWK JSON string processing and RSA private key conversion
- [x] Service account token support (authflow token compatibility)
- [x] Multiple output formats (text, JSON, YAML)
- [x] Complete layered architecture (CLI → Public API → Internal)
- [x] Comprehensive testing suite with 76.8% coverage
- [x] Internal API for cross-command integration

### Phase 3: ELK Migration (NEXT)
- [ ] Migrate Docker Compose orchestration
- [ ] Integrate token generation for authenticated PAIC API calls
- [ ] Migrate Python log streaming to Go
- [ ] Maintain shell script compatibility during transition
- [ ] Add improved status monitoring and health checks

### Phase 4: Journey Migration
- [ ] Migrate authentication journey functionality
- [ ] Maintain YAML configuration compatibility
- [ ] Preserve all existing authflow journey features
- [ ] Add Go-specific improvements (better error handling, performance)

### Phase 5: Integration & Enhancement
- [ ] Unified configuration system
- [ ] Cross-tool integration (journey + token + elk workflows)
- [ ] Enhanced error handling and logging
- [ ] Performance optimizations
- [ ] Plugin architecture for future tools

### Phase 6: Future Expansion
- [ ] Native Frodo integration commands
- [ ] Configuration management tools
- [ ] Environment management utilities
- [ ] Additional PAIC testing capabilities

## Established Patterns & Rules

### **✅ Proven Architecture Patterns**

**1. Layered Architecture (Successfully Implemented)**
```
CLI Layer (cmd/)     →  Public API (pkg/)     →  Internal Logic (internal/)
     ↓                        ↓                         ↓  
pctl token           →  token.NewClient()     →  ServiceAccountGenerator
Command-line UX      →  Importable Library    →  Business Logic + PAIC API
```

**2. Configuration Strategy (Authflow Compatible)**
- **Primary**: YAML configuration files (compatible with existing authflow configs)
- **Secondary**: CLI flag overrides for quick modifications  
- **Hierarchy**: CLI flags > Environment vars > Config files > Defaults
- **Format**: snake_case YAML fields (e.g., `service_account_id`, `exp_seconds`)

**3. Real Implementation Pattern (Token ✅ Complete)**
- **JWK Processing**: Parse JWK JSON strings → Convert to RSA private keys
- **JWT Creation**: RS256 algorithm with proper claims (iss, sub, aud, exp, jti)
- **OAuth 2.0 Flow**: JWT Bearer Token grant type to PAIC endpoints
- **Response Handling**: Parse PAIC responses with comprehensive error handling

### **✅ Testing Standards (76.8% Coverage)**

**4. Comprehensive Test Strategy**
```
Unit Tests (pkg/*_test.go)        →  Test public API interfaces
Integration Tests (internal/)     →  Test business logic components  
Internal API Tests (test/)        →  Test cross-command usage patterns
```

**5. Test Patterns**
- **Table-driven tests** with subtests for comprehensive coverage
- **Temporary file creation** for config testing  
- **Mock data** for testing without external PAIC dependencies
- **Error scenario validation** with meaningful error messages
- **Multiple output format testing** (text, JSON, YAML)

### **✅ Development Rules**

**6. Code Organization Standards**
1. **cmd/**: CLI command definitions using Cobra (thin wrappers)
2. **pkg/**: Public APIs for external consumption and internal integration
3. **internal/**: Private business logic and implementation details  
4. **configs/**: YAML configuration templates (examples/ and real/)
5. **test/**: Integration tests and internal API usage patterns

**7. Security & Configuration Rules**
- **Sensitive Data**: `configs/*/real/` directories in .gitignore
- **JWK Handling**: Treat JWK as opaque JSON strings (no crypto knowledge needed)
- **Token Security**: Never log private keys or sensitive token details
- **Error Handling**: Comprehensive validation with user-friendly messages

### **✅ Command Integration Patterns**

**8. Internal API Usage** (For ELK Integration)
```go
// Load configuration
config, err := pkgtoken.LoadConfig("path/to/config.yaml")

// Create client options
options := pkgtoken.GeneratorOptions{
    Config:       *config,
    OutputFormat: pkgtoken.OutputFormatJSON,
    Verbose:      false,
}

// Generate token for internal use
client := pkgtoken.NewClient(options)
result, err := client.Generate()

// Access token for authenticated API calls
accessToken := result.AccessToken
```

**9. Cross-Command Integration**
- **Token → ELK**: ELK commands use token generation for authenticated PAIC API calls
- **Configuration Sharing**: Common YAML patterns across all commands
- **Error Consistency**: Standardized error handling and user messaging

## Development Guidelines

### CLI Design Standards  
1. **Flattened Commands**: `pctl journey`, `pctl token`, `pctl elk` (no nesting)
2. **Config-First Design**: YAML configs with CLI flag overrides
3. **Help System**: Comprehensive help with examples for all commands
4. **Validation**: Input validation with helpful, actionable error messages
5. **Verbose Mode**: Detailed logging when requested (`-v` flag)

## Tool Integration Requirements

### Token Command (`pctl token`) ✅ **COMPLETED**
**Migration from**: `authflow token`
- ✅ **Implemented**: Real JWT token generation with PAIC OAuth 2.0 flows
- ✅ **Implemented**: Service account tokens using JWK JSON strings
- ✅ **Implemented**: Full authflow token config compatibility (YAML)
- ✅ **Implemented**: Multiple output formats (text, JSON, YAML)
- ✅ **Implemented**: Comprehensive testing suite (76.8% coverage)
- ✅ **Implemented**: Internal API for cross-command integration
- 🚀 **Enhanced**: Better error handling, validation, and performance than original

### ELK Command (`pctl elk`) **NEXT PRIORITY**
**Migration from**: `plctl.sh`
- **Must Have**: Docker orchestration, log streaming, lifecycle management
- **Must Have**: Integration with `pctl token` for authenticated PAIC API calls  
- **Should Have**: Platform detection, background modes, all plctl.sh features
- **Nice to Have**: Native Go log processing, improved monitoring

### Journey Command (`pctl journey`) **FUTURE**
**Migration from**: `authflow journey`
- **Must Have**: Authentication journey testing with YAML configs
- **Must Have**: Integration with `pctl token` for authentication flows
- **Should Have**: All existing authflow journey features and options
- **Nice to Have**: Performance improvements, better error messages

### Integration Status
- ✅ **Token Generation**: Production-ready with real PAIC OAuth 2.0 implementation
- 🔄 **ELK Integration**: Ready for token-authenticated PAIC API calls
- 📋 **Journey Integration**: Will leverage existing token generation
- 🎯 **Cross-Command**: Established patterns for internal API usage

## Configuration Management

### Configuration Hierarchy
1. **System Config**: `/etc/pctl/config.yaml`
2. **User Config**: `~/.pctl/config.yaml` 
3. **Project Config**: `./pctl.yaml`
4. **Environment Variables**: `PCTL_*` prefix
5. **Command Line**: Highest priority

### Configuration Format (YAML)
```yaml
# Core PCTL configuration
version: "1.0"
log_level: "info"
output_format: "text"

# Tool-specific configurations
tools:
  authflow:
    default_config_path: "./configs/authflow"
    timeout: 30s
    
  elk:
    docker_compose_timeout: 300s
    default_retention: "7d"

# Environment management
environments:
  default: "development"
  profiles:
    development:
      frodo_profile: "dev"
    production:
      frodo_profile: "prod"
```

## Security & Compliance

### Security Requirements
1. **Credential Management**: Never log or expose credentials
2. **Configuration Security**: Secure handling of sensitive configs
3. **Network Security**: Proper TLS/SSL validation
4. **Access Control**: Respect system permissions and access controls

### Compliance Features
- **Audit Logging**: Track all tool operations
- **Configuration Validation**: Prevent insecure configurations
- **Secure Defaults**: Security-first default settings
- **Credential Rotation**: Support for rotating credentials

## Development Environment

### Prerequisites
- Go 1.21 or later
- Docker and Docker Compose (for ELK functionality)
- Frodo CLI installed and configured
- Git for version control

### Development Commands
```bash
# Initialize project
go mod init github.com/your-org/pctl
go mod tidy

# Build development version
go build -o bin/pctl

# Run tests
go test ./...

# Build all platforms
make build-all

# Development mode
go run main.go <command>

# Example usage:
go run main.go journey -c configs/journey/example.yaml
go run main.go token -c configs/token/service-account.yaml  
go run main.go elk start
```

### Testing Strategy ✅ **IMPLEMENTED**
- ✅ **Unit Tests**: All packages have comprehensive tests (76.8% coverage)
- ✅ **Integration Tests**: Internal API usage patterns and cross-command integration
- ✅ **Configuration Tests**: YAML config loading, validation, and compatibility
- ✅ **Error Handling Tests**: Comprehensive error scenario validation
- ✅ **Output Format Tests**: Multiple format validation (text, JSON, YAML)
- 🔄 **End-to-End Tests**: Full workflow testing with real PAIC environments (token ✅)
- 📋 **Performance Tests**: Ensure Go version matches or improves performance

### Build & Development Commands ✅ **ESTABLISHED**
```bash
# Development and testing
go test -v ./...                    # Run all tests with verbose output
go test -cover ./...                # Run tests with coverage reporting  
go build -o bin/pctl                # Build development binary
go run main.go <command>            # Run in development mode

# Token functionality (production-ready)
./bin/pctl token -c configs/token/real/service-account.yaml -v
./bin/pctl token -c configs/token/examples/service-account.yaml -o json
```

## Notes for Claude

### ✅ **Current Achievement Status**
- 🎉 **Token Generation**: PRODUCTION READY with real PAIC OAuth 2.0 flows
- 🎯 **Architecture**: Proven layered design with 76.8% test coverage
- 🔧 **Internal APIs**: Established patterns for cross-command integration
- 📋 **Next Target**: ELK functionality with token-based authentication

### Development Priorities (UPDATED)
1. ✅ **Token Foundation**: Complete with real PAIC integration
2. 🔄 **ELK Integration**: Use token generation for authenticated PAIC API calls
3. 📋 **Journey Migration**: Leverage existing token and config patterns
4. **Cross-Platform**: Must work consistently across Linux/macOS/Windows
5. **Enterprise Focus**: Production-grade reliability and comprehensive testing

### Established Patterns (USE THESE)
- **Configuration**: YAML-first with CLI overrides, authflow-compatible
- **Architecture**: CLI → Public API → Internal (proven with token)
- **Testing**: Unit + Integration + Internal API usage patterns
- **Error Handling**: Comprehensive validation with user-friendly messages
- **Internal Integration**: `pkg/` APIs for cross-command usage
- **Security**: Sensitive configs in `configs/*/real/` (gitignored)

### Key Architectural Decisions
- **Flattened Commands**: `pctl token`, `pctl elk`, `pctl journey` (no nesting)
- **Real Implementation**: No mocks - actual PAIC OAuth 2.0 flows
- **JWK as Strings**: Treat JWK JSON as opaque strings (successful pattern)
- **Layered Testing**: Unit (pkg), Integration (internal), Usage (test/)
- **Config Compatibility**: Support existing authflow YAML formats

### Future Vision (REFINED)
PCTL is becoming the comprehensive CLI platform that:
- ✅ Generates real PAIC tokens for authentication
- 🔄 Integrates ELK functionality with token-based auth
- 📋 Consolidates journey testing with unified token generation
- 🎯 Provides superior developer experience with comprehensive testing
- 🚀 Scales to enterprise requirements with production-ready reliability