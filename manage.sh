#!/bin/bash

# Скрипт для управления Docker-окружением

function show_help {
    echo "Использование: ./manage.sh [команда]"
    echo "Команды:"
    echo "  start       - Запустить контейнеры (--profile fetcher_only)"
    echo "  stop        - Остановить контейнеры (сохраняя данные)"
    echo "  restart     - Перезапустить контейнеры"
    echo "  rebuild     - Пересобрать и запустить контейнеры"
    echo "  clean       - Остановить и удалить все контейнеры и сети"
    echo "  fullclean   - Остановить и удалить все контейнеры, сети и тома"
    echo "  status      - Показать статус контейнеров"
    echo "  logs [service] - Показать логи контейнера (postgres, fetcher или main)"
    echo "  follow [service] - Показать и следить за логами контейнера"
    echo "  help        - Показать эту справку"
}

case "$1" in
    start)
        docker-compose --profile fetcher_only up -d
        ;;
    stop)
        docker-compose --profile fetcher_only down
        ;;
    restart)
        docker-compose --profile fetcher_only down
        docker-compose --profile fetcher_only up -d
        ;;
    rebuild)
        docker-compose --profile fetcher_only down
        docker-compose --profile fetcher_only up --build -d
        ;;
    clean)
        docker-compose down --remove-orphans
        docker system prune -f
        ;;
    fullclean)
        docker-compose down --volumes --remove-orphans
        docker system prune -af --volumes
        ;;
    status)
        docker-compose ps
        ;;
    logs)
        if [ -z "$2" ]; then
            docker-compose logs
        else
            docker-compose logs "$2"
        fi
        ;;
    follow)
        if [ -z "$2" ]; then
            docker-compose logs -f
        else
            docker-compose logs -f "$2"
        fi
        ;;
    help|*)
        show_help
        ;;
esac
