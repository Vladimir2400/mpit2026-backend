#!/bin/bash

# SSH Reverse Tunnel для локальной разработки
# Пробрасывает локальный порт 8080 на сервер

SERVER="root@138.124.72.73"
LOCAL_PORT=8080
REMOTE_PORT=8080

echo "🚇 Создаю SSH туннель..."
echo "   Локальный порт: $LOCAL_PORT"
echo "   Удаленный порт: $REMOTE_PORT"
echo "   Сервер: $SERVER"
echo ""
echo "💡 Чтобы остановить туннель, нажмите Ctrl+C"
echo ""

# Запуск туннеля с автоматическим переподключением
while true; do
    echo "⏳ Подключаюсь к серверу..."

    ssh -R $REMOTE_PORT:localhost:$LOCAL_PORT $SERVER \
        -N \
        -o ServerAliveInterval=60 \
        -o ServerAliveCountMax=3 \
        -o ExitOnForwardFailure=yes \
        -o StrictHostKeyChecking=no

    EXIT_CODE=$?

    if [ $EXIT_CODE -eq 0 ]; then
        echo "✅ Туннель закрыт корректно"
        break
    else
        echo "❌ Туннель оборвался (код: $EXIT_CODE). Переподключаюсь через 5 секунд..."
        sleep 5
    fi
done
