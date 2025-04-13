# GraphRAG API

## О проекте

GraphRAG API — это сервер для работы с текстами. Он позволяет загружать текстовые файлы, добавлять тексты и искать ответы на вопросы. Вы можете загрузить файл или отправить текст, а затем задать вопрос — сервер найдёт подходящий ответ.

## Как это работает

- **Загрузка**: Отправьте файл `.txt`, чтобы сохранить его содержимое.
- **Добавление**: Отправьте тексты в формате JSON.
- **Поиск**: Задайте вопрос, и сервер вернёт подходящий текст.

## Установка

### Что нужно

- Go (версия 1.21 или выше)
- PostgreSQL с расширением `pgvector`
- Python 3.12 с библиотекой `sentence-transformers`

### Клонирование

Склонируйте проект с GitHub и переключитесь на ветку `GraphRAG_LLM-Project`:

```bash
git clone https://github.com/Gybernaty-Neurocore/PHILOSOPHOS-AI-GRAPH-RAG-SYSTEM.git
cd PHILOSOPHOS-AI-GRAPH-RAG-SYSTEM
git checkout GraphRAG_LLM-Project

### Установка зависимостей
Установите зависимости для Go:

bash

- go mod tidy
- pip install sentence-transformers==3.2.1 transformers==4.46.0
- docker run -d --name ... -e POSTGRES_PASSWORD=... -p ...:5432 ankane/pgvector

### .env
PG_PASSWORD=...

### Запусти сервер:
- go run cmd/main.go
