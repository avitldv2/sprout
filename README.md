# Sprout

A minimal static site generator written in Go. Sprout converts Markdown files and HTML templates into static websites with fast incremental rebuilds.

## Installation

### From Source

```bash
git clone https://github.com/avitldv2/sprout
cd sprout
make build
sudo make install
```

## Quick Start

```bash
sprout init
sprout build
sprout serve
```

Visit `http://localhost:1313` to see your site.

## Documentation

- **[BUILDING.md](BUILDING.md)** - Complete guide on how Sprout works, building sites, templates, styling, and more

## Commands

- `sprout init` - Initialize a new site
- `sprout build` - Build the site
- `sprout serve` - Start development server

See `sprout <command> --help` for options.

## Features

- Incremental builds (only rebuilds what changed)
- Markdown + HTML templates
- Pretty URLs with automatic directory structure
- Live reload for development
- Single binary, no dependencies
- Minimal and very fast

## License

This project is licensed under the MIT License
