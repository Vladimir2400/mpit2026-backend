# Статус разработки Dating App Backend

**Версия**: MVP 0.1
**Обновлено**: 2025-12-04
**Цель**: Хакатон MPIT 2026 (3-5 дней разработки)

---

## 📊 Общий прогресс

```
День 1: Инфраструктура                    ████████████████████ 100%
День 2: Профили и рекомендации            ░░░░░░░░░░░░░░░░░░░░   0%
День 3: Свайпы и матчи                    ░░░░░░░░░░░░░░░░░░░░   0%
День 4: Чат и подписки                    ░░░░░░░░░░░░░░░░░░░░   0%
День 5: Полировка и тесты                 ░░░░░░░░░░░░░░░░░░░░   0%

Общий прогресс MVP:                       ████░░░░░░░░░░░░░░░░  20%
```

---

## ✅ Реализованный функционал

### 🏗️ Инфраструктура (100%)

#### ✅ Инициализация проекта
- [x] Go модуль (`go.mod`)
- [x] Clean Architecture структура папок
- [x] Зависимости установлены (Gin, GORM, JWT, Redis, Viper, etc.)
- [x] `.gitignore` настроен
- [x] `.env.example` и `.env` созданы

#### ✅ Конфигурация
- [x] `internal/config/config.go` - Загрузка настроек через Viper
- [x] Валидация критических параметров
- [x] Поддержка environment variables
- [x] Структурированная конфигурация (Server, DB, Redis, JWT, Encryption)

#### ✅ Docker окружение
- [x] `docker-compose.yml` - PostgreSQL 15 + Redis 7
- [x] Health checks для сервисов
- [x] Volumes для персистентности данных
- [x] Сетевая изоляция

#### ✅ Domain модели (11 файлов)
- [x] `user.go` - Пользователи (валидация 18+, гендер, роли)
- [x] `profile.go` - Профили (MBTI, интересы, геолокация, предпочтения)
- [x] `swipe.go` - Свайпы (left/right/super)
- [x] `match.go` - Матчи (взаимные лайки)
- [x] `message.go` - Сообщения и Conversations
- [x] `subscription.go` - Подписки (планы, статусы)
- [x] `verification.go` - KYC верификация (мок)
- [x] `consent.go` - Согласия пользователей (ФЗ-152)
- [x] `recommendation.go` - Логи AI рекомендаций
- [x] `errors.go` - Доменные ошибки

#### ✅ PKG утилиты (100%)
- [x] `pkg/crypto/hash.go` - bcrypt хеширование (cost 12)
- [x] `pkg/crypto/encrypt.go` - AES-256-GCM шифрование/дешифрование
- [x] `pkg/jwt/jwt.go` - JWT генерация и валидация (Access + Refresh)
- [x] `pkg/validator/validator.go` - Валидация (email, phone, password, age, MBTI)

#### ✅ База данных
- [x] `migrations/000001_init_schema.up.sql` - Все 11 таблиц
- [x] `migrations/000001_init_schema.down.sql` - Rollback
- [x] Индексы для производительности (GIST, GIN, B-tree)
- [x] Constraints (CHECK, UNIQUE, FK)
- [x] Поддержка UUID, массивов, JSONB, геолокации

**Таблицы БД:**
1. `users` - Пользователи
2. `profiles` - Профили
3. `photos` - Фотографии
4. `verification_data` - Верификация
5. `user_consents` - Согласия
6. `swipes` - Свайпы
7. `matches` - Матчи
8. `conversations` - Беседы
9. `messages` - Сообщения
10. `subscriptions` - Подписки
11. `recommendation_logs` - Логи рекомендаций

#### ✅ Инструменты разработки
- [x] `Makefile` - Команды для разработки
- [x] `README.md` - Полная документация
- [x] `STATUS.md` - Этот файл

---

## 🔄 В процессе реализации

### 🚧 День 1 (продолжение) - Аутентификация

**Текущие задачи:**
- [ ] Репозитории (PostgreSQL + Redis)
  - [ ] `repository/postgres/user.go`
  - [ ] `repository/postgres/profile.go`
  - [ ] `repository/postgres/consent.go`
  - [ ] `repository/redis/cache.go`
  - [ ] `repository/redis/session.go`

- [ ] Auth Usecase
  - [ ] `usecase/auth/service.go`
  - [ ] Register (с валидацией возраста)
  - [ ] Login (JWT генерация)
  - [ ] Refresh token
  - [ ] Logout

- [ ] HTTP Layer
  - [ ] `delivery/http/handler/auth.go`
  - [ ] `delivery/http/middleware/auth.go`
  - [ ] `delivery/http/middleware/logger.go`
  - [ ] `delivery/http/middleware/cors.go`
  - [ ] `delivery/http/router.go`

- [ ] Main
  - [ ] `cmd/api/main.go`
  - [ ] Инициализация БД
  - [ ] Запуск сервера

---

## 📅 Предстоящий функционал

### День 2: Профили и рекомендации

#### Профили
- [ ] `usecase/profile/service.go` - Управление профилями
- [ ] `repository/postgres/profile.go` - CRUD профилей
- [ ] `repository/postgres/photo.go` - Управление фото
- [ ] `delivery/http/handler/profile.go` - HTTP handlers

**Endpoints:**
- `GET /api/v1/profile` - Получить свой профиль
- `PUT /api/v1/profile` - Обновить профиль
- `POST /api/v1/profile/photos` - Загрузить фото
- `DELETE /api/v1/profile/photos/:id` - Удалить фото
- `PUT /api/v1/profile/preferences` - Обновить предпочтения поиска

#### Верификация (МОК)
- [ ] `usecase/verification/service.go` - Мок-верификация
- [ ] `repository/postgres/verification.go`
- [ ] `delivery/http/handler/verification.go`

**Endpoints:**
- `POST /api/v1/verification/submit` - Загрузить документы (автоверификация)
- `GET /api/v1/verification/status` - Статус верификации

#### AI Рекомендации
- [ ] `usecase/recommendation/service.go` - Сервис рекомендаций
- [ ] `usecase/recommendation/algorithm.go` - Rule-based алгоритм
- [ ] `usecase/recommendation/scoring.go` - Scoring логика
- [ ] `repository/postgres/recommendation.go` - Логирование
- [ ] `delivery/http/handler/recommendation.go`

**Алгоритм:**
- Cold start (новые пользователи)
- Персонализированные рекомендации
- MBTI совместимость (30%)
- Общие интересы (40%)
- Расстояние (20%)
- Активность (10%)
- Балансировка популярности

**Endpoints:**
- `GET /api/v1/recommendations` - Лента рекомендаций (limit, offset)

### День 3: Свайпы и матчи

#### Свайпы
- [ ] `usecase/swipe/service.go` - Обработка свайпов
- [ ] `repository/postgres/swipe.go` - История свайпов
- [ ] `delivery/http/handler/swipe.go`

**Бизнес-логика:**
- Проверка активной подписки (мужчины)
- Детекция взаимного лайка
- Создание матча при совпадении
- Обновление popularity_score
- Логирование для AI

**Endpoints:**
- `POST /api/v1/swipes` - Свайпнуть (body: {swiped_user_id, direction})
- `GET /api/v1/swipes/history` - История свайпов

#### Матчи
- [ ] `usecase/match/service.go` - Управление матчами
- [ ] `repository/postgres/match.go` - CRUD матчей
- [ ] `delivery/http/handler/match.go`

**Бизнес-логика:**
- Создание при взаимном right swipe
- Автоматическое создание conversation
- Отправка уведомлений (mock для MVP)
- Размэтч (unmatch)

**Endpoints:**
- `GET /api/v1/matches` - Все активные матчи
- `GET /api/v1/matches/:id` - Детали матча
- `DELETE /api/v1/matches/:id` - Размэтчить

### День 4: Чат и подписки

#### Чат
- [ ] `usecase/chat/service.go` - Сервис чата
- [ ] `usecase/chat/encryption.go` - Шифрование сообщений
- [ ] `repository/postgres/message.go` - Хранение сообщений
- [ ] `repository/postgres/conversation.go` - Управление беседами
- [ ] `delivery/http/handler/chat.go`

**Бизнес-логика:**
- AES-256-GCM шифрование
- Проверка участия в conversation
- Обновление last_message_preview
- Mark as read
- Пагинация сообщений

**Endpoints:**
- `GET /api/v1/conversations` - Список чатов (с превью)
- `GET /api/v1/conversations/:id` - Детали беседы
- `GET /api/v1/conversations/:id/messages` - Сообщения (пагинация)
- `POST /api/v1/conversations/:id/messages` - Отправить сообщение
- `PUT /api/v1/conversations/:id/read` - Отметить прочитанным
- `DELETE /api/v1/messages/:id` - Удалить сообщение (soft delete)

#### Подписки (МОК)
- [ ] `usecase/subscription/service.go` - Управление подписками
- [ ] `repository/postgres/subscription.go`
- [ ] `delivery/http/handler/subscription.go`

**Бизнес-логика:**
- Мужчины: платная подписка (мок - бесплатно)
- Женщины: автоматическая бесплатная подписка
- Согласие на автопродление (ФЗ-69)
- Проверка активности подписки

**Endpoints:**
- `GET /api/v1/subscriptions/plans` - Доступные тарифы
- `GET /api/v1/subscriptions/current` - Текущая подписка
- `POST /api/v1/subscriptions` - Создать подписку (мок)
- `DELETE /api/v1/subscriptions/:id` - Отменить подписку

### День 5: Полировка и тесты

#### Доработки
- [ ] Обработка ошибок (unified error responses)
- [ ] Валидация всех входных данных
- [ ] Seed данные для демо
- [ ] Swagger документация (опционально)
- [ ] Логирование (структурированное)

#### Тестирование
- [ ] Unit тесты (critical usecase logic)
  - [ ] Recommendation algorithm
  - [ ] Match creation
  - [ ] Message encryption
  - [ ] JWT validation
- [ ] Integration тесты
  - [ ] Auth flow
  - [ ] Swipe → Match flow
  - [ ] Chat flow

#### Performance
- [ ] Проверка индексов БД
- [ ] Redis кэширование (profiles, recommendations)
- [ ] Query optimization

---

## 🎯 Полный список API Endpoints

### 🔐 Аутентификация (`/api/v1/auth`)

| Метод | Endpoint | Описание | Статус | Body |
|-------|----------|----------|--------|------|
| POST | `/register` | Регистрация пользователя | ⏳ Planned | `{email, phone, password, gender, birth_date}` |
| POST | `/login` | Вход в систему | ⏳ Planned | `{email/phone, password}` |
| POST | `/refresh` | Обновление access token | ⏳ Planned | `{refresh_token}` |
| POST | `/logout` | Выход из системы | ⏳ Planned | - |

**Response:**
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "gender": "male",
    "role": "user"
  }
}
```

---

### 👤 Профиль (`/api/v1/profile`)

| Метод | Endpoint | Описание | Статус | Auth |
|-------|----------|----------|--------|------|
| GET | `/profile` | Получить свой профиль | ⏳ Planned | ✅ Required |
| PUT | `/profile` | Обновить профиль | ⏳ Planned | ✅ Required |
| GET | `/profile/:user_id` | Получить чужой профиль | ⏳ Planned | ✅ Required |
| GET | `/preferences` | Получить настройки поиска | ⏳ Planned | ✅ Required |
| PUT | `/preferences` | Обновить настройки | ⏳ Planned | ✅ Required |

**Update Profile Body:**
```json
{
  "display_name": "Иван",
  "bio": "Люблю путешествия",
  "city": "Москва",
  "mbti_type": "INTJ",
  "interests": ["путешествия", "кино", "спорт"],
  "location_lat": 55.7558,
  "location_lon": 37.6173,
  "pref_min_age": 25,
  "pref_max_age": 35,
  "pref_max_distance_km": 50,
  "pref_gender": "female"
}
```

---

### 📸 Фотографии (`/api/v1/profile/photos`)

| Метод | Endpoint | Описание | Статус | Auth |
|-------|----------|----------|--------|------|
| POST | `/photos` | Загрузить фото | ⏳ Planned | ✅ Required |
| GET | `/photos` | Список фото профиля | ⏳ Planned | ✅ Required |
| PUT | `/photos/:id/primary` | Сделать основным | ⏳ Planned | ✅ Required |
| PUT | `/photos/:id/order` | Изменить порядок | ⏳ Planned | ✅ Required |
| DELETE | `/photos/:id` | Удалить фото | ⏳ Planned | ✅ Required |

**Upload Body:**
```
Content-Type: multipart/form-data
file: [binary]
order_index: 0
is_primary: true
```

---

### ✅ Верификация (`/api/v1/verification`)

| Метод | Endpoint | Описание | Статус | Auth |
|-------|----------|----------|--------|------|
| POST | `/submit` | Загрузить документы | ⏳ Planned | ✅ Required |
| GET | `/status` | Статус верификации | ⏳ Planned | ✅ Required |
| POST | `/consent` | Дать согласие на биометрию | ⏳ Planned | ✅ Required |

**Submit Body (МОК - автоверификация):**
```json
{
  "selfie_photo_url": "/uploads/selfie.jpg",
  "passport_photo_url": "/uploads/passport.jpg",
  "biometric_consent_given": true
}
```

---

### 🎯 Рекомендации (`/api/v1/recommendations`)

| Метод | Endpoint | Описание | Статус | Auth |
|-------|----------|----------|--------|------|
| GET | `/` | Лента рекомендаций | ⏳ Planned | ✅ Required |

**Query params:**
- `limit` - Количество (default: 20)
- `offset` - Смещение (pagination)

**Response:**
```json
{
  "recommendations": [
    {
      "user_id": "uuid",
      "profile": {
        "display_name": "Анна",
        "age": 28,
        "city": "Москва",
        "bio": "...",
        "photos": [...],
        "mbti_type": "ENFP",
        "interests": ["музыка", "йога"]
      },
      "compatibility_score": 0.85,
      "reasoning": {
        "mbti_match": true,
        "shared_interests": 3,
        "distance_km": 5.2
      },
      "distance_km": 5.2
    }
  ],
  "total": 150,
  "has_more": true
}
```

---

### 👆 Свайпы (`/api/v1/swipes`)

| Метод | Endpoint | Описание | Статус | Auth |
|-------|----------|----------|--------|------|
| POST | `/` | Свайпнуть на профиль | ⏳ Planned | ✅ Required |
| GET | `/history` | История свайпов | ⏳ Planned | ✅ Required |
| GET | `/stats` | Статистика свайпов | ⏳ Planned | ✅ Required |

**Swipe Body:**
```json
{
  "swiped_user_id": "uuid",
  "direction": "right"  // "left", "right", "super"
}
```

**Response при матче:**
```json
{
  "success": true,
  "match": {
    "id": "uuid",
    "matched_user": {...},
    "matched_at": "2025-12-04T12:00:00Z"
  }
}
```

---

### 💕 Матчи (`/api/v1/matches`)

| Метод | Endpoint | Описание | Статус | Auth |
|-------|----------|----------|--------|------|
| GET | `/` | Все активные матчи | ⏳ Planned | ✅ Required |
| GET | `/:id` | Детали матча | ⏳ Planned | ✅ Required |
| DELETE | `/:id` | Размэтчить | ⏳ Planned | ✅ Required |

**Response:**
```json
{
  "matches": [
    {
      "id": "uuid",
      "matched_user": {
        "id": "uuid",
        "display_name": "Мария",
        "age": 27,
        "photos": [...]
      },
      "matched_at": "2025-12-04T10:30:00Z",
      "conversation_id": "uuid",
      "last_message": "Привет!",
      "unread_count": 2
    }
  ],
  "total": 15
}
```

---

### 💬 Чат (`/api/v1/conversations`)

| Метод | Endpoint | Описание | Статус | Auth |
|-------|----------|----------|--------|------|
| GET | `/` | Список чатов | ⏳ Planned | ✅ Required |
| GET | `/:id` | Детали беседы | ⏳ Planned | ✅ Required |
| GET | `/:id/messages` | Сообщения | ⏳ Planned | ✅ Required |
| POST | `/:id/messages` | Отправить сообщение | ⏳ Planned | ✅ Required |
| PUT | `/:id/read` | Отметить прочитанным | ⏳ Planned | ✅ Required |
| DELETE | `/messages/:id` | Удалить сообщение | ⏳ Planned | ✅ Required |

**Send Message Body:**
```json
{
  "content": "Привет! Как дела?"
}
```

**Messages Response:**
```json
{
  "messages": [
    {
      "id": "uuid",
      "sender_id": "uuid",
      "content": "Привет!",  // Дешифровано
      "is_read": true,
      "created_at": "2025-12-04T12:00:00Z"
    }
  ],
  "has_more": true,
  "total": 48
}
```

---

### 💳 Подписки (`/api/v1/subscriptions`)

| Метод | Endpoint | Описание | Статус | Auth |
|-------|----------|----------|--------|------|
| GET | `/plans` | Доступные тарифы | ⏳ Planned | - |
| GET | `/current` | Текущая подписка | ⏳ Planned | ✅ Required |
| POST | `/` | Создать подписку (мок) | ⏳ Planned | ✅ Required |
| DELETE | `/:id` | Отменить подписку | ⏳ Planned | ✅ Required |
| POST | `/consent` | Согласие на автопродление | ⏳ Planned | ✅ Required |

**Plans Response:**
```json
{
  "plans": [
    {
      "type": "monthly",
      "price_cents": 0,  // МОК: бесплатно
      "duration_days": 30,
      "features": [
        "Unlimited swipes",
        "See who liked you",
        "Super likes",
        "Compatibility score"
      ]
    }
  ]
}
```

---

### 🔧 Служебные (`/`)

| Метод | Endpoint | Описание | Статус | Auth |
|-------|----------|----------|--------|------|
| GET | `/health` | Health check | ⏳ Planned | - |
| GET | `/api/v1/privacy-policy` | Политика конфиденциальности | ⏳ Planned | - |

---

## 🎨 Структура Response

### Успешный ответ
```json
{
  "success": true,
  "data": { ... },
  "message": "Operation successful"
}
```

### Ошибка
```json
{
  "success": false,
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "Invalid email or password",
    "details": { ... }
  }
}
```

---

## 🔒 Безопасность

### Реализовано
- ✅ JWT токены (Access: 15 мин, Refresh: 30 дней)
- ✅ bcrypt хеширование (cost 12)
- ✅ AES-256-GCM шифрование
- ✅ Валидация входных данных
- ✅ SQL Injection защита (GORM)
- ✅ Проверка возраста 18+ (БД constraint)

### Планируется
- ⏳ CORS middleware
- ⏳ Rate limiting (Redis)
- ⏳ Request logging
- ⏳ XSS защита
- ⏳ CSRF токены (опционально)

---

## 📈 Метрики и мониторинг (опционально)

### Планируется после MVP
- [ ] Prometheus metrics
- [ ] Grafana dashboards
- [ ] Error tracking (Sentry)
- [ ] Performance monitoring
- [ ] Database query analytics

---

## 🚀 Команды для разработки

```bash
# Установка
make setup              # Полная установка (docker + migrations)

# Разработка
make dev                # Режим разработки (docker + run)
make run                # Запустить приложение
make docker-up          # Запустить PostgreSQL + Redis
make migrate-up         # Применить миграции

# Тестирование
make test               # Запустить тесты

# Утилиты
make seed               # Заполнить тестовыми данными
make clean              # Очистить артефакты
```

---

## 📝 Примечания

### MVP Ограничения
- 🔸 Чат через REST API (не WebSocket)
- 🔸 Мок-подписки (без реальных платежей)
- 🔸 Мок-верификация (автоматическое подтверждение)
- 🔸 Без push уведомлений
- 🔸 Без admin панели

### Расширения после MVP
- 🚀 WebSocket для real-time чата
- 💰 Интеграция ЮKassa
- 🧠 ML рекомендации (Python + gRPC)
- 📱 Push notifications (FCM)
- 🛡️ Admin panel
- 🔍 PostGIS для геопоиска
- ⚡ Advanced rate limiting

---

## 📞 Контакты

**Проект**: MPIT 2026 Hackathon
**Backend**: Dating App MVP
**Архитектура**: Clean Architecture
**Язык**: Go 1.21+

---

**Последнее обновление**: 2025-12-04
**Следующая цель**: Завершить День 1 - Аутентификация ✅
