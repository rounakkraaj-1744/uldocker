# Dawker

> A keyboard-first Docker TUI with a command-driven workflow and real-time feedback.

Dawker is a terminal user interface for Docker that merges the visual feedback of a dashboard with the speed of a command-line interface. Built in Go using the Bubbletea framework, it replaces tedious mouse interactions and complex CLI flags with a Vim-inspired command engine.

---

## Why Dawker Exists

Managing local Docker environments is often a fragmented experience. The standard Docker CLI requires excessive typing, context switching, and constantly looking up container IDs. Existing Docker TUIs or GUI tools typically rely on hidden menus, deeply nested popups, or mouse clicks, which break developer flow.

Dawker was built with a single philosophy: **interaction should be immediate.**

By treating the terminal as an integrated environment and utilizing a command buffer (`:`), Dawker allows you to navigate, monitor, and manipulate your Docker resources exactly like you edit code in Vim.

---

## When to Use Dawker

Dawker is best suited for:

- Inspecting running containers quickly
- Debugging services via real-time log streaming
- Performing routine lifecycle operations (start, stop, restart)
- Cleaning up local Docker environments

It is not intended for:

- Complex orchestration workflows
- Building images or managing CI/CD pipelines
- Replacing the Docker CLI entirely

---

## Demo

![Dawker Demo](assets/demo.gif)

---

## Quick Start

Get up and running in under a minute.

### Prerequisites
* Go 1.21+
* A local Docker Daemon running at `/var/run/docker.sock`

### Install (One-liner)

The fastest way to install Dawker on Mac or Linux:

```bash
curl -sSL https://raw.githubusercontent.com/rounakkraaj-1744/dawker/main/install.sh | bash
```

### Install (From Source)

```bash
git clone https://github.com/rounakkraaj-1744/dawker.git
cd dawker
go build -o dawker cmd/main.go
./dawker
```

### Install (Binary)

Download the latest release and install:

```bash
curl -L https://github.com/rounakkraaj-1744/dawker/releases/latest/download/dawker-linux -o dawker
chmod +x dawker
sudo mv dawker /usr/local/bin/
```

### Fast Path (First 60 seconds)
1. **Move:** Use `j` and `k` to scroll through your containers.
2. **View Logs:** Press `Enter` on a running container to stream its logs.
3. **Execute:** Press `:` to open the command buffer. Type `stop` and press `Enter`. The highlighted container stops immediately.

---

## Core Concepts

Dawker's mental model is based on four pillars:

1. **Tabs (Resources)**: Resources are grouped into four primary views: `Containers`, `Images`, `Volumes`, and `Networks`. Switching tabs completely changes your operational context.
2. **Selection (`▸`)**: There is always an active, selected item in your list. Many commands default to targeting this specific item if no other arguments are provided.
3. **Command Mode (`:`)**: The core interaction model. Rather than mapping 50 different keyboard shortcuts for different actions, everything complex is handled through an interactive command buffer at the bottom of the screen.
4. **Details / Logs Panel**: The right half of the screen dynamically updates. If a container is streaming, it shows multiplexed stdout/stderr. Otherwise, it shows static metadata about the selected resource.

---

## Usage Model

Dawker is designed around a simple interaction loop:

1. Navigate to a resource (container, image, volume, or network)
2. Inspect its state (status indicator or live logs)
3. Execute an action using command mode (`:`)

Most operations follow this pattern:

**Select → Inspect → Act**

The command system removes the need to remember exact Docker CLI syntax. Instead, actions are expressed in terms of intent:

```
:stop api       → stop the container matching "api"
:rm exited      → remove all stopped containers
:rmi node       → remove the image matching "node"
```

This keeps interaction consistent across all resource types.

---

## Detailed Usage Guide

### Navigation

Dawker operates in a normal navigation mode by default.

| Key | Action |
| --- | --- |
| `j` / `k` / `↓` / `↑` | Move selection down / up |
| `1` / `2` / `3` / `4` | Switch tab (Containers, Images, Volumes, Networks) |
| `Enter` | Stream logs (for containers) or view static metadata |
| `:` | Enter Command Mode |
| `Tab` | Auto-complete current suggestion (in command mode) |
| `↑` / `↓` | Cycle command history (in command mode) |
| `r` | Force refresh resource data |
| `q` / `Esc` | Quit application or cancel active log stream |

### Working with Containers
When in the **Containers** tab, the TUI provides visual indicators for state (Green for `running`, Red for `exited`). Pressing `Enter` on a container allocates a background `io.Pipe`, separates multiplexed Docker logs via `stdcopy`, and streams them in real-time to the right panel. The log buffer is capped at 500 lines and auto-scrolls to the bottom, with line truncation to prevent horizontal overflow.

### Command Mode

Press `:` to enter command mode. All input is captured by the command buffer until execution or cancellation. A `-- COMMAND --` indicator appears in the footer, and a block cursor marks your position.

* **Autocomplete:** Press `Tab` to auto-complete commands or resource names.
* **History:** Use `Up` and `Down` arrows to cycle through previous successful commands (up to 50 entries).
* **Inline Errors:** If you type a typo (e.g., `:sto`), Dawker suggests the closest match (`Did you mean: stop?`).

### Common Commands

```
:stop api         Stop containers matching "api"
:restart all      Restart all containers
:rm exited        Remove all stopped containers
:prune            Reclaim disk space from unused resources
```

### Command Resolution Rules

Commands are evaluated using the following rules:

1. **Explicit target takes priority**
   - `:stop api` → matches containers containing "api" in their name

2. **No arguments → use current selection**
   - `:stop` → applies to the item highlighted by `▸` in the active tab

3. **Keyword targets**
   - `all` → all resources in the current tab
   - `running` → containers currently running
   - `exited` → containers that have stopped

4. **Multiple matches**
   - Commands are applied to all matching resources sequentially

5. **No matches**
   - Command fails with inline error feedback (no crash, no modal)

### Smart Behavior

The command system is designed for fast, low-friction interaction.

- Commands are resolved using fuzzy matching against resource names
- Targets are inferred from the active tab and selection when omitted
- Batch operations are supported through keywords (`all`, `running`, `exited`)

**Context-Aware Execution**: If you type `:rm` without arguments, Dawker inspects your current tab and active selection. If you are on the Images tab, it removes the selected image. If you are on the Volumes tab, it removes the selected volume.

**Fuzzy Matching**: You do not need to type full container hashes or names. If you have a container named `production-api-server-1`, typing `:stop api` matches it based on partial name.

**Batch Keywords**: Typing `:rm exited` removes all containers that are no longer running. `:stop all` halts every container in the list.

**Live Previews**: As you type, the UI renders a preview line above the input (e.g., *Will stop: api-server, cache-db*). You see exactly what the command will affect before pressing `Enter`.

---

## Command Reference

All commands execute against the Docker SDK directly. `[target]` can be a partial name, a keyword (`all`, `running`, `exited`), or omitted to use the current selection.

### Lifecycle Management (Containers)
| Command | Description |
| --- | --- |
| `:start [target]` | Starts the target container(s) |
| `:stop [target]` | Gracefully stops the target container(s) |
| `:restart [target]` | Restarts the target container(s) |

### Cleanup & Removal
| Command | Description |
| --- | --- |
| `:rm [target]` | Removes containers (force removal) |
| `:rmi [target]` | Removes images |
| `:rmv [target]` | Removes volumes |
| `:rmn [target]` | Removes networks |
| `:prune` | System-wide prune of unused containers, images, volumes, and networks |

### Inspection
| Command | Description |
| --- | --- |
| `:logs [target]` | Trigger log stream (WIP) |
| `:stats [target]` | View resource utilization (WIP) |

---

## Practical Workflows

### Debugging a Failing Service
1. Notice a red `(exited)` status on your API container.
2. Press `Enter` to view the crash logs in the right panel.
3. Fix the code in your editor.
4. Press `:`, type `restart`, press `Enter`. The container restarts and the UI refreshes to show `(running)`.

### Cleaning a Dirty Environment
1. Your local Docker environment is cluttered with stopped containers and dangling images.
2. Press `:`, type `rm exited`, press `Enter` to clear dead containers.
3. Press `:`, type `prune`, press `Enter` to reclaim disk space from all unused resources.

### Managing Multiple Services
1. You have 5 microservices running (`auth-api`, `billing-api`, `user-api`, etc.).
2. You need to stop all of them.
3. Press `:`, type `stop api`.
4. The preview confirms: *Will stop: auth-api, billing-api, user-api*.
5. Press `Enter`.

---

## Architecture

Dawker is strictly layered to ensure the UI never blocks and the Docker daemon is never overwhelmed.

* **TUI Layer (`internal/app`)**: Built on `charmbracelet/bubbletea` and `lipgloss`. The `Update` loop is decomposed into `handleNormalModeKey` and `handleCommandModeKey` to prevent state collision. UI rendering happens concurrently without blocking Docker operations.
* **Command System (`internal/command`)**:
  * `parser.go`: Tokenizes raw strings into command structs.
  * `selector.go`: Fuzzy-matching logic and keyword filters (`all`, `running`, `exited`).
  * `executor.go`: A registry pattern mapping parsed commands to Docker SDK actions.
  * `suggest.go`: Evaluates the input buffer in real-time for autocomplete and typo corrections.
* **Docker Layer (`internal/docker`)**: Wraps the official Go Docker SDK. Uses a singleton `*client.Client` via `sync.Once` to eliminate repeated API version negotiation and connection overhead.

**Execution Flow:**

```
Input string → Parser → Fuzzy Target Resolution → Executor Registry → Docker SDK Action → UI State Reload
```

---

## Design Decisions

* **Keyboard-First vs Modals**: Traditional TUIs rely on pop-up modals (e.g., "Are you sure? [Y/n]"). Dawker replaces modals with a command preview line. Seeing *Will stop: database* is faster and more transparent than a generic confirmation prompt.
* **Single Connection**: Wrapping the Docker client in a singleton ensures commands execute without connection setup delay. This is critical for responsive command mode interaction.
* **Minimal UI**: There are no heavy borders or excessive colors. Color is used strictly for semantics (Green = running, Red = exited, Orange = Command Mode active).

---

## Limitations

- Only supports local Docker daemon via Unix socket (`/var/run/docker.sock`)
- Remote Docker hosts are not configurable
- Log streaming is capped at 500 lines and does not support filtering or search
- Does not support Docker Compose, Swarm, or Kubernetes

---

## Contributing

Contributions are welcome, particularly for expanding the `executor.go` registry with more Docker SDK actions.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/new-handler`)
3. Commit your changes (`git commit -m 'feat: add stats handler'`)
4. Push to the branch (`git push origin feature/new-handler`)
5. Open a Pull Request

---

## License

MIT License
