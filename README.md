# appimaged

Simple AppImage installer and desktop integrator for Linux.

`appimaged` automatically:

- Copies AppImages into `~/Applications`
- Makes them executable
- Installs icons
- Creates `.desktop` launchers
- Adds applications to your desktop dashboard/menu

Works well on:

- Fedora
- Ubuntu
- Debian
- Arch Linux
- Most Linux desktop environments using Freedesktop `.desktop` entries

---

# Features

- Minimal and lightweight
- No dependencies
- Interactive CLI
- Written in Go
- User-local installation (no root required)

---

# Build

## Requirements

- Go 1.20+

## Compile

```bash
go build -o appimaged# appimaged
