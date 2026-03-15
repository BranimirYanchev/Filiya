# Contributing to Filia Project Backend

Thank you for considering contributing to the Filia Project Backend! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for everyone.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the issue tracker to avoid duplicates. When you create a bug report, include as many details as possible:

- A clear and descriptive title
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- Screenshots if applicable
- Environment details (OS, Go version, etc.)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- Use a clear and descriptive title
- Provide a detailed description of the suggested enhancement
- Explain why this enhancement would be useful
- Include any relevant examples or mockups

### Pull Requests

- Fill in the required template
- Follow the Go coding style
- Include tests for new features
- Update documentation for changes
- End all files with a newline

## Development Setup

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/your-username/filia-project-backend.git
   cd filia-project-backend
   ```
3. Set up the development environment:
   ```bash
   go mod download
   ```
4. Create a `.env` file with your local database configuration

## Coding Guidelines

### Go Style Guide

- Follow the [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `gofmt` to format your code
- Write meaningful comments and documentation
- Keep functions small and focused on a single task
- Use descriptive variable and function names

### Testing

- Write tests for all new features and bug fixes
- Ensure all tests pass before submitting a pull request
- Aim for high test coverage

### Documentation

- Update the README.md if necessary
- Update API documentation for any endpoint changes
- Add comments to explain complex code sections

## Git Workflow

1. Create a new branch for each feature or bugfix
   ```bash
   git checkout -b feature/your-feature-name
   ```
   or
   ```bash
   git checkout -b fix/your-bugfix-name
   ```

2. Make your changes and commit them with clear, descriptive messages
   ```bash
   git commit -m "Add feature: description of the feature"
   ```

3. Push your branch to your fork
   ```bash
   git push origin feature/your-feature-name
   ```

4. Create a Pull Request against the main repository

## Review Process

- All submissions require review
- Changes may be requested before a pull request is merged
- The reviewer may close a pull request if it doesn't meet the guidelines

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Gin Framework Documentation](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)

Thank you for contributing to Filia Project Backend!