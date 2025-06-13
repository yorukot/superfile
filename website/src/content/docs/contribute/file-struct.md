---
title: superfile Project Structure Guide
description: A detailed guide to understanding superfile's codebase organization
head:
  - tag: title
    content: superfile Project Structure Guide | superfile
---

# superfile Project Structure Guide

The project follows a standard Go project layout with clear separation of concerns. Here's a detailed breakdown of the main directories and their purposes:

## Core Directories

### `src/` - Main Source Code

The main source code is organized into several key directories:

#### `cmd/` - Entry Point

- `main.go` - The main entry point of the application that handles:
  - CLI argument parsing
  - Configuration initialization
  - Application startup

#### `config/` - Configuration Management

- `fixed_variable.go` - Contains constant values and configuration paths
- `icon/` - Icon-related configuration
  - `function.go` - Icon initialization and management functions
  - `icon.go` - Icon definitions and mappings

#### `internal/` - Core Application Logic

Contains the main business logic of the application, organized by functionality:

**Configuration & Types:**

- `config_function.go` - Configuration loading and management
- `config_type.go` - Configuration-related type definitions
- `default_config.go` - Default configuration values
- `type.go` - Core type definitions

**File Operations:**

- `file_operations.go` - Basic file operation functions
- `file_operations_compress.go` - File compression functionality
- `file_operations_extract.go` - File extraction functionality
- `handle_file_operations.go` - File operation handlers

**UI & Interaction:**

- `handle_modal.go` - Modal dialog management
- `handle_panel_movement.go` - Panel navigation logic
- `handle_panel_navigation.go` - Panel focus management
- `handle_pinned_operations.go` - Pinned items functionality
- `key_function.go` - Keyboard input handling
- `model.go` - Core application model
- `model_render.go` - UI rendering logic

**Utilities:**

- `function.go` - General utility functions
- `get_data.go` - Data retrieval functions
- `string_function.go` - String manipulation utilities
- `string_function_test.go` - String utility tests
- `style.go` - UI styling definitions
- `style_function.go` - UI styling functions
- `string_function_test.go` - String utility tests
- `style.go` - UI styling definitions
- `style_function.go` - UI styling functions

### `testsuite/` - superfile's testsuite written in Python

- Automatically tests superfile's functionality.
- See `testsuite/ReadMe.md` for more info

## Code Organization Principles

1. **Separation of Concerns:**

   - Configuration management is isolated in the `config/` directory
   - Core business logic lives in `internal/`
   - UI-related code is separated from business logic

2. **Modular Design:**

   - Each file has a specific responsibility
   - Related functionality is grouped together
   - Clear dependencies between components

3. **Testing:**
   - Test files are placed alongside the code they test
   - Example: `string_function_test.go` tests `string_function.go`

## Contributing Guidelines

When contributing to superfile:

1. **Adding New Features:**

   - Place new business logic in appropriate `internal/` subdirectories
   - Keep UI-related code separate from business logic
   - Follow existing naming conventions

2. **Making Changes:**

   - Maintain the existing file structure
   - Add tests for new functionality
   - Update configuration files if needed

3. **Code Style:**
   - Follow Go best practices
   - Maintain consistent formatting
   - Add appropriate documentation

This structure helps maintain code organization and makes it easier for new contributors to understand where to make changes.
