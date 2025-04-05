# PHILOSOPHOS AI GraphRAG Server

GraphRAG-сервер на Go с использованием Gin, PostgreSQL, pgvector, Hugging Face и модели Qwen для загрузки документов, поиска по ним и генерации ответов.

## Возможности

- Загрузка документов
- Векторизация через Qwen (через Hugging Face API)
- Поиск похожих документов
- Генерация ответа на запрос

---

## Запуск проекта

### 1. Клонировать репозиторий

```bash
git clone https://github.com/yourusername/PHILOSOPHOS-AI-GRAPH-RAG-SYSTEM-1.git

cd PHILOSOPHOS-AI-GRAPH-RAG-SYSTEM-1

2. Установить зависимости

go mod tidy

3. Создать БД

4. Заполнить .env
Создайте файл .env в корне проекта и добавьте:

HF_TOKEN=your_huggingface_token
DB_PASSWORD=yourpassword

4. Запустить сервер

go run cmd/main.go

   Используемые технологии

Go + Gin

PostgreSQL + pgvector

Hugging Face API (Qwen2.5-7B-Instruct)

GORM

    Важно
Убедитесь, что токен Hugging Face (HF_TOKEN) действителен.

Используйте модель с поддержкой text-embedding и chat-completion.
