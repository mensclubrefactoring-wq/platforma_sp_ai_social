# Настройка AI моделей

Платформа спроектирована так, чтобы вы могли легко переключаться между разными поставщиками AI (OpenAI, Gemini, Local LLM).

## Как поменять модель по API-адресу

В файле `main.go` или через переменные окружения (`.env`) вы можете изменить следующие параметры:

1. **`AI_BASE_URL`**: Базовый адрес API. 
   - По умолчанию используется прокси Gemini через OpenAI-совместимый интерфейс: `https://generativelanguage.googleapis.com/v1beta/openai`
   - Для OpenAI: `https://api.openai.com`
   - Для локальной Llama (через Ollama): `http://localhost:11434/v1`

2. **`AI_MODEL`**: Название модели.
   - Примеры: `gpt-4o`, `gemini-1.5-flash`, `llama3`.

3. **`GEMINI_API_KEY`** (или `OPENAI_API_KEY` в коде): Ваш секретный ключ.

### Пример для локальной модели (Ollama)
Если вы хотите запускать всё полностью на своем компьютере без интернета:
1. Установите [Ollama](https://ollama.com/).
2. Запустите модель: `ollama run llama3`.
3. В `.env` укажите:
   ```env
   AI_BASE_URL=http://localhost:11434/v1
   AI_MODEL=llama3
   GEMINI_API_KEY=any_string
   ```

## Где это находится в коде
В `main.go` за это отвечает функция `callAI`. Она использует стандартный формат `POST /v1/chat/completions`, который поддерживают почти все современные AI сервисы.
