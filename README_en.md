# ApiNest

> Open-source API collaboration platform (v0.0.1)

[中文](README.md) | [English](README_en.md)

## What it does

ApiNest centralizes API documentation and team collaboration. Put your team in a shared workspace, organize APIs by project and folder, configure environments and variables, invite members to edit together, and share selected APIs via links that visitors can view without logging in.

It replaces scattered spreadsheets and private API lists, keeping API definitions, environment switching, and collaboration permissions in one place.

## Tech stack

| Layer | Technologies |
|------|------|
| Backend | Go · Gin · Ent · JWT · Zap · WebSocket |
| Frontend | Vue 3 · TypeScript · Vite · Element Plus · Vue Router · Axios |
| Infrastructure | Docker · PostgreSQL · Valkey · Caddy |

## Features

- **Workspaces & projects**: Create, rename, and delete; organize API assets by project
- **API docs**: Folder hierarchy, API CRUD, request/response field management
- **Environments & variables**: Multi-environment configs for local / staging / production
- **Team collaboration**: Invite members to a workspace with role-based edit permissions
- **Project sharing**: Share records so guests can view selected APIs without login
- **System**: Install wizard, auth, admin console (users & workspace transfer), zh/en UI

## Roadmap

- [ ] Built-in / multi-database support (currently external PostgreSQL only)
- [ ] Desktop client integration
- [ ] Liquid-glass theme and border refinements

## License

ApiNest is licensed under the [Apache License 2.0](LICENSE). Retain the copyright notice and license text when modifying or redistributing.
