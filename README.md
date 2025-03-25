# ArikawaX

ArikawaX is an extension package for the [Arikawa Discord bot framework](https://github.com/diamondburned/arikawa), providing additional utilities and middleware to enhance Discord bot development in Go.

## Features

### Middleware
This package provides interaction handler middlewares:

- **Timeout Detection**: Logs interactions that takes too long to return
- **Panic Recovery**: General panic recovery
- **Logging**: Http server like in and out logging

## Requirements

- github.com/diamondburned/arikawa/v3

## Usage

### Importing the Package

```go
import "github.com/No3371/arikawax/middleware"
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [Arikawa Discord Framework](https://github.com/diamondburned/arikawa) for the base Discord bot framework 