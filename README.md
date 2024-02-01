# Restful api

## Usage

- Docker
  `make docker-build && make docker-run`

- Develop on local
  [Please install gin to reload server automatically](https://github.com/codegangsta/gin)
  `make live`

## Goal

implement a restful task API application, which includes the following endpoints:

- GET `/tasks`
- POST` /tasks`
- PUT `/tasks/{id}`
- DELETE `/tasks/{id}`

A task should contain at least the following fields:

- `name`
  - type: string
  - description:task name
- `status`
  - type: integer
  - enum:[0,1]
  - description:0 represents an incomplete task, while 1 represents a completed task

## DOD

- unit tests
- Provides Dockerfile to run API in Docker

## Principle

- 根據 https://github.com/bxcodec/go-clean-arch 專案的 clear architecture 架構，不過實作上很少會抽換 usecase，所以這邊實作並沒有多抽一層
- 透過 clear architecture 來實作，可以讓程式碼更好維護，並且可以更好的測試
- 自己維護事務操作，需要考慮的情境 https://chat.openai.com/share/8ce9e2a0-787f-4b0d-a2e1-abc7a613051f
