# Dawker
A keyboard-first, highly responsive Docker TUI for power users.

## Overview
Dawker is a terminal user interface (TUI) designed to manage Docker resources without leaving the keyboard. Built in Go using the Bubbletea framework, it provides a fast, Vim-inspired command system that replaces tedious mouse clicks and long CLI strings. It is built for developers and sysadmins who need immediate access to container logs, states, and lifecycle management through an interface that feels like a modern development environment.

## Key Features
* **Container Dashboard**: Real-time status monitoring for containers, images, volumes, and networks organized in clean, responsive tabs.
* **Logs Streaming**: Instant, multiplexed (stdout/stderr) log streaming directly within the right-hand panel.
* **Command Mode**: A Vim-like execution buffer (`:`) allowing direct action against Docker resources without switching contexts.
* **Smart Commands**: Type-ahead autocomplete, argument suggestions, and inline previews showing exactly what resources will be affected before execution.
* **Keyboard-First UX**: Zero mouse dependency. Navigate, execute, and monitor using standard `j/k` bindings and simple keystrokes.

## Demo
*(Placeholder for Demo GIF/Screenshot showing command mode and log streaming)*

## Installation

### Prerequisites
* Go 1.21 or higher
* Docker Daemon running locally (`/var/run/docker.sock`)

### Build and Run
Clone the repository and run the application directly:
```bash
git clone https://github.com/yourusername/uldocker.git
cd uldocker
go build -o dawker cmd/main.go
./dawker
```

## Usage Guide

### Navigation
Dawker operates primarily in a normal navigation mode. Use these keys to move through your resources:

| Key | Action |
| --- | --- |
| `j` / `k` / `↓` / `↑` | Move selection down / up |
| `enter` | Stream logs (for containers) or view details |
| `1` - `4` | Switch between Containers, Images, Volumes, and Networks tabs |
| `:` | Enter Command Mode |
| `r` | Refresh resource data |
| `q` / `esc` | Quit application or cancel current stream/input |

### Command Mode
Press `:` to enter the command buffer at the bottom of the screen. Commands execute against the Docker SDK directly.

Supported commands: `start`, `stop`, `restart`, `rm`, `rmi`, `rmv`, `rmn`, `prune`, `logs`, `stats`.

**Examples:**
* `:stop api` — Stops the container fuzzily matching "api"
* `:restart all` — Restarts all containers in the current list
* `:rmi node` — Removes the image fuzzily matching "node"
* `:prune` — Reclaims space by removing unused containers, images, volumes, and networks

### Smart Behavior
* **Fuzzy Matching**: You do not need to type the full container ID or name. Typing `:stop ap` will match the `api-server` container automatically.
* **Context-Aware Execution**: If you type a command without arguments (e.g., `:start`), it will execute against the currently highlighted item in your active tab.
* **Batch Execution**: Passing keywords like `all`, `running`, or `exited` (e.g., `:rm exited`) will apply the action to all matching resources.
* **Command Preview**: The UI renders an italicized preview line above the input (e.g., *Will stop: api-server, db-worker*) so you know exactly what the command will hit before pressing Enter.

## How It Works (Architecture)
Dawker is structured to separate presentation from business logic:
* **TUI Layer (`internal/app`)**: Built on `charmbracelet/bubbletea`. `update.go` handles state transitions and keystrokes, cleanly separated into normal mode and command mode handling. `view.go` handles UI rendering using `lipgloss`.
* **Command System (`internal/command`)**: Handles input parsing (`parser.go`), autocomplete logic (`suggest.go`), target resolution via fuzzy matching (`selector.go`), and action dispatch (`executor.go`).
* **Docker Integration (`internal/docker`)**: Wraps the official Docker Go SDK. Utilizes a singleton client pattern (`sync.Once`) to eliminate connection latency during command execution and data polling.
* **Data Flow**: `Keystroke → Update State → Parse Command → Resolve Targets → Execute SDK Action → Return Result → Render UI`.

## Design Decisions
* **Keyboard-First**: Context switching between a keyboard and mouse breaks flow. Dawker maps standard Vim motions to ensure hands stay on the home row.
* **Command Mode over Modals**: Instead of popping up complex modal dialogues to confirm actions, a CLI-style input buffer combined with live previews provides faster, more transparent execution.
* **Not a Docker Replacement**: Dawker focuses on the 90% use case—managing running environments, checking logs, and cleaning up. It explicitly does not attempt to handle complex image building or compose orchestration.
* **Minimal UI**: We chose a high signal-to-noise ratio. Information density is prioritized over excessive borders or colors, ensuring the data is readable at a glance.

## Limitations
* Requires a local Docker daemon (unix socket). Remote TCP Docker hosts are not currently supported out-of-the-box.
* Log streaming is currently capped and does not support advanced filtering (like filtering by string matches within the stream).
* Does not support Docker Swarm or Kubernetes clusters.

## Contributing
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License
MIT License
