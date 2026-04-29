# Запуск

1. cp .env.example .env и заполнить WORDSTAT_*, CURRENCY_API_KEY (для compose подстановки положите .env рядом с docker-compose.yml — Compose подхватывает его для ${…}).
2. docker compose up --build -d postgres app
3. Загрузка CSV: docker compose --profile import run --rm importer
4. Пример: curl "http://localhost:8081/tools/api/wordstat?query=iPhone%16"