# inspirai-devkit

Cross-platform CLI tools for inspirai monorepo management.

## Installation

```bash
go install github.com/inspirai-store/inspirai-devkit/cmd/sm@latest
```

## Commands

| Command | Description |
|---------|-------------|
| `sm init` | Initialize all submodules and create symlinks |
| `sm sync` | Sync all submodules (git pull) |
| `sm status` | Show status of all submodules |
| `sm links` | Rebuild all symlinks |

## Development

```bash
cd .submodules/inspirai-devkit
go build -o sm ./cmd/sm
./sm --help
```
