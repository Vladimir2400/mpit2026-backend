# API Documentation - MVP

## Authentication

### POST /auth/vk
Авторизация через VK Mini App

**Request:**
```json
{
  "vk_params": {
    "vk_user_id": "123456",
    "vk_app_id": "51234567",
    "vk_is_app_user": "1",
    "vk_platform": "mobile_web",
    "sign": "..."
  }
}
```

**Response 200:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": 1735123456,
  "user": {
    "id": 1,
    "vk_id": 123456,
    "gender": "male",
    "birth_date": "2000-01-01",
    "is_verified": false,
    "is_online": true,
    "created_at": "2024-12-04T10:00:00Z"
  },
  "is_new_user": true
}
```

**Response 401:**
```json
{
  "error": "invalid VK signature"
}
```

---

### POST /auth/logout
Выход из системы

**Headers:**
- `Authorization: Bearer <token>`

**Response 200:**
```json
{
  "message": "logged out successfully"
}
```

---

### GET /auth/me
Получить информацию о текущем пользователе

**Headers:**
- `Authorization: Bearer <token>`

**Response 200:**
```json
{
  "user_id": 1
}
```

---

## Profile

### GET /profile/me
Получить свой профиль

**Headers:**
- `Authorization: Bearer <token>`

**Response 200:**
```json
{
  "id": 1,
  "user_id": 1,
  "display_name": "Иван",
  "bio": "Люблю путешествия и музыку",
  "city": "Москва",
  "interests": ["музыка", "спорт", "путешествия"],
  "location_lat": 55.7558,
  "location_lon": 37.6173,
  "location_updated_at": "2024-12-04T10:00:00Z",
  "pref_min_age": 18,
  "pref_max_age": 30,
  "pref_max_distance_km": 50,
  "is_onboarding_complete": true,
  "created_at": "2024-12-04T10:00:00Z",
  "updated_at": "2024-12-04T10:00:00Z"
}
```

**Response 404:**
```json
{
  "error": "profile not found"
}
```

---

### PUT /profile/me
Обновить свой профиль

**Headers:**
- `Authorization: Bearer <token>`

**Request:**
```json
{
  "display_name": "Иван",
  "bio": "Обновленная биография",
  "city": "Москва",
  "interests": ["музыка", "спорт"],
  "location_lat": 55.7558,
  "location_lon": 37.6173,
  "pref_min_age": 20,
  "pref_max_age": 28,
  "pref_max_distance_km": 30
}
```

**Response 200:**
```json
{
  "id": 1,
  "user_id": 1,
  "display_name": "Иван",
  "bio": "Обновленная биография",
  "city": "Москва",
  "interests": ["музыка", "спорт"],
  "location_lat": 55.7558,
  "location_lon": 37.6173,
  "pref_min_age": 20,
  "pref_max_age": 28,
  "pref_max_distance_km": 30,
  "is_onboarding_complete": true,
  "updated_at": "2024-12-04T11:00:00Z"
}
```

---

### POST /profile/complete-onboarding
Завершить онбординг (создать профиль при первом входе)

**Headers:**
- `Authorization: Bearer <token>`

**Request:**
```json
{
  "display_name": "Иван",
  "bio": "Люблю путешествия",
  "city": "Москва",
  "interests": ["музыка", "спорт"],
  "pref_min_age": 18,
  "pref_max_age": 30,
  "pref_max_distance_km": 50
}
```

**Response 201:**
```json
{
  "id": 1,
  "user_id": 1,
  "display_name": "Иван",
  "is_onboarding_complete": true,
  "created_at": "2024-12-04T10:00:00Z"
}
```

---

### GET /profile/:user_id
Получить профиль другого пользователя

**Headers:**
- `Authorization: Bearer <token>`

**Response 200:**
```json
{
  "id": 2,
  "user_id": 2,
  "display_name": "Мария",
  "bio": "Фотограф и путешественница",
  "city": "Санкт-Петербург",
  "interests": ["фотография", "искусство", "путешествия"],
  "age": 25,
  "distance_km": 12.5,
  "created_at": "2024-12-03T10:00:00Z"
}
```

**Response 404:**
```json
{
  "error": "profile not found"
}
```

---

## Big Five Personality Test (TIPI)

### GET /big-five/questions
Получить вопросы TIPI теста

**Response 200:**
```json
{
  "questions": [
    {"id": 1, "text": "Экстраверт, энергичный"},
    {"id": 2, "text": "Критичный, склонный к спорам"},
    {"id": 3, "text": "Надёжный, дисциплинированный"},
    {"id": 4, "text": "Тревожный, легко расстраиваюсь"},
    {"id": 5, "text": "Открытый новому, со сложным внутренним миром"},
    {"id": 6, "text": "Сдержанный, тихий"},
    {"id": 7, "text": "Отзывчивый, тёплый"},
    {"id": 8, "text": "Неорганизованный, беспечный"},
    {"id": 9, "text": "Спокойный, эмоционально стабильный"},
    {"id": 10, "text": "Консервативный, не склонный к творчеству"}
  ],
  "instruction": "Оцените, насколько каждое утверждение описывает вас, по шкале от 1 (совершенно не согласен) до 7 (полностью согласен)"
}
```

---

### POST /big-five/submit
Отправить ответы на TIPI тест

**Headers:**
- `Authorization: Bearer <token>`

**Request:**
```json
{
  "answers": {
    "1": 6,
    "2": 3,
    "3": 7,
    "4": 2,
    "5": 6,
    "6": 2,
    "7": 7,
    "8": 2,
    "9": 6,
    "10": 3
  }
}
```

**Response 201:**
```json
{
  "id": 1,
  "user_id": 1,
  "openness": 0.75,
  "conscientiousness": 0.83,
  "extraversion": 0.67,
  "agreeableness": 0.83,
  "neuroticism": 0.33,
  "completed_at": "2024-12-04T10:00:00Z",
  "created_at": "2024-12-04T10:00:00Z",
  "updated_at": "2024-12-04T10:00:00Z"
}
```

**Response 400:**
```json
{
  "error": "must answer all 10 questions"
}
```

**Response 409:**
```json
{
  "error": "test already completed"
}
```

---

### GET /big-five/my-results
Получить свои результаты Big Five теста

**Headers:**
- `Authorization: Bearer <token>`

**Response 200:**
```json
{
  "id": 1,
  "user_id": 1,
  "openness": 0.75,
  "conscientiousness": 0.82,
  "extraversion": 0.68,
  "agreeableness": 0.90,
  "neuroticism": 0.35,
  "completed_at": "2024-12-04T10:00:00Z",
  "created_at": "2024-12-04T10:00:00Z"
}
```

**Response 404:**
```json
{
  "error": "test not completed yet"
}
```

---

### GET /big-five/user/:user_id
Получить результаты Big Five теста другого пользователя

**Headers:**
- `Authorization: Bearer <token>`

**Response 200:**
```json
{
  "id": 5,
  "user_id": 5,
  "openness": 0.80,
  "conscientiousness": 0.70,
  "extraversion": 0.85,
  "agreeableness": 0.75,
  "neuroticism": 0.40,
  "completed_at": "2024-12-03T10:00:00Z"
}
```

**Response 404:**
```json
{
  "error": "test results not found"
}
```

---

## Feed (Лента пользователей)

### GET /feed/next
Получить следующего пользователя в ленте

**Headers:**
- `Authorization: Bearer <token>`

**Response 200:**
```json
{
  "user": {
    "id": 5,
    "user_id": 5,
    "display_name": "Анна",
    "bio": "Люблю спорт и активный отдых",
    "city": "Москва",
    "age": 24,
    "interests": ["спорт", "йога", "бег"],
    "distance_km": 3.2
  }
}
```

**Response 204 (нет данных):**
```json
{
  "message": "no more users in feed"
}
```

---

### POST /feed/reset-dislikes
Сбросить все дизлайки (обновить ленту)

**Headers:**
- `Authorization: Bearer <token>`

**Response 200:**
```json
{
  "message": "dislikes reset successfully",
  "reset_count": 15
}
```

---

## Swipes (Лайки/Дизлайки)

### POST /swipe
Поставить лайк или дизлайк

**Headers:**
- `Authorization: Bearer <token>`

**Request:**
```json
{
  "swiped_user_id": 5,
  "is_like": true
}
```

**Response 200 (взаимный лайк):**
```json
{
  "is_match": true,
  "match": {
    "id": 3,
    "user1_id": 1,
    "user2_id": 5,
    "created_at": "2024-12-04T12:00:00Z"
  },
  "matched_user": {
    "id": 5,
    "display_name": "Анна",
    "bio": "Люблю спорт и активный отдых",
    "city": "Москва",
    "age": 24
  }
}
```

**Response 200 (обычный лайк/дизлайк):**
```json
{
  "is_match": false,
  "swipe": {
    "id": 10,
    "swiper_id": 1,
    "swiped_id": 5,
    "is_like": true,
    "created_at": "2024-12-04T12:00:00Z"
  }
}
```

**Response 400:**
```json
{
  "error": "cannot swipe yourself"
}
```

**Response 409:**
```json
{
  "error": "swipe already exists"
}
```

---

### GET /swipe/likes-received
Получить список людей, которые поставили мне лайк

**Headers:**
- `Authorization: Bearer <token>`

**Query params:**
- `limit` (optional, default: 20)
- `offset` (optional, default: 0)

**Response 200:**
```json
{
  "likes": [
    {
      "swipe_id": 15,
      "user": {
        "id": 7,
        "user_id": 7,
        "display_name": "Мария",
        "bio": "Художница",
        "city": "Москва",
        "age": 23,
        "interests": ["искусство", "кино"],
        "distance_km": 5.1
      },
      "created_at": "2024-12-04T11:30:00Z"
    }
  ],
  "total": 5
}
```

---

## Matches (Симпатии)

### GET /matches
Получить список всех совпадений с данными пользователей

**Headers:**
- `Authorization: Bearer <token>`

**Query params:**
- `limit` (optional, default: 20)
- `offset` (optional, default: 0)

**Response 200:**
```json
{
  "matches": [
    {
      "match_id": 3,
      "user": {
        "id": 5,
        "user_id": 5,
        "display_name": "Анна",
        "bio": "Люблю спорт",
        "city": "Москва",
        "age": 24,
        "interests": ["спорт", "йога"],
        "is_online": true,
        "last_online_at": "2024-12-04T12:00:00Z"
      },
      "matched_at": "2024-12-04T12:00:00Z",
      "last_message": {
        "content": "Привет! Как дела?",
        "sender_id": 5,
        "created_at": "2024-12-04T12:05:00Z",
        "is_read": false
      },
      "unread_messages": 2
    }
  ],
  "total": 10
}
```

---

### DELETE /matches/:match_id
Удалить совпадение (размачить)

**Headers:**
- `Authorization: Bearer <token>`

**Response 200:**
```json
{
  "message": "match deleted successfully"
}
```

**Response 404:**
```json
{
  "error": "match not found"
}
```

---

## Messages (Чаты)

### GET /messages/conversations
Получить список всех чатов с последними сообщениями

**Headers:**
- `Authorization: Bearer <token>`

**Query params:**
- `limit` (optional, default: 20)
- `offset` (optional, default: 0)

**Response 200:**
```json
{
  "conversations": [
    {
      "match_id": 3,
      "user": {
        "id": 5,
        "user_id": 5,
        "display_name": "Анна",
        "is_online": true,
        "last_online_at": "2024-12-04T12:30:00Z"
      },
      "last_message": {
        "id": 25,
        "content": "Привет! Как дела?",
        "sender_id": 5,
        "is_read": false,
        "created_at": "2024-12-04T12:05:00Z"
      },
      "unread_count": 2
    }
  ],
  "total": 5
}
```

---

### GET /messages/:match_id
Получить все сообщения в чате (polling каждые 5 секунд)

**Headers:**
- `Authorization: Bearer <token>`

**Query params:**
- `limit` (optional, default: 50)
- `offset` (optional, default: 0)
- `since` (optional) - ISO 8601 timestamp, вернет только новые сообщения после указанной даты

**Response 200:**
```json
{
  "messages": [
    {
      "id": 25,
      "match_id": 3,
      "sender_id": 5,
      "content": "Привет! Как дела?",
      "is_read": false,
      "created_at": "2024-12-04T12:05:00Z"
    },
    {
      "id": 24,
      "match_id": 3,
      "sender_id": 1,
      "content": "Привет!",
      "is_read": true,
      "created_at": "2024-12-04T12:00:00Z"
    }
  ],
  "total": 15,
  "has_new_messages": true
}
```

**Response 403:**
```json
{
  "error": "unauthorized to access this conversation"
}
```

**Пример использования с polling:**
```javascript
// Первый запрос - получить все сообщения
GET /messages/3?limit=50&offset=0

// Каждые 5 секунд запрашивать новые сообщения
setInterval(() => {
  const lastMessageTime = "2024-12-04T12:05:00Z"; // timestamp последнего сообщения
  GET /messages/3?since=2024-12-04T12:05:00Z
}, 5000);
```

---

### POST /messages/:match_id
Отправить сообщение

**Headers:**
- `Authorization: Bearer <token>`

**Request:**
```json
{
  "content": "Привет! Как дела?"
}
```

**Response 201:**
```json
{
  "id": 26,
  "match_id": 3,
  "sender_id": 1,
  "content": "Привет! Как дела?",
  "is_read": false,
  "created_at": "2024-12-04T12:10:00Z"
}
```

**Response 403:**
```json
{
  "error": "unauthorized to send message to this match"
}
```

---

### PUT /messages/:message_id/read
Отметить сообщение как прочитанное

**Headers:**
- `Authorization: Bearer <token>`

**Response 200:**
```json
{
  "message": "message marked as read"
}
```

---

### PUT /messages/:match_id/read-all
Отметить все сообщения в чате как прочитанные

**Headers:**
- `Authorization: Bearer <token>`

**Response 200:**
```json
{
  "message": "all messages marked as read",
  "count": 5
}
```

---

## Notifications (Уведомления)

### GET /notifications
Получить список уведомлений

**Headers:**
- `Authorization: Bearer <token>`

**Query params:**
- `limit` (optional, default: 20)
- `offset` (optional, default: 0)
- `unread_only` (optional, default: false)

**Response 200:**
```json
{
  "notifications": [
    {
      "id": 10,
      "user_id": 1,
      "content": "У вас новое совпадение с Анной!",
      "is_read": false,
      "created_at": "2024-12-04T12:00:00Z"
    },
    {
      "id": 9,
      "user_id": 1,
      "content": "Мария поставила вам лайк",
      "is_read": true,
      "created_at": "2024-12-04T11:30:00Z"
    }
  ],
  "total": 15,
  "unread_count": 3
}
```

---

### PUT /notifications/:notification_id/read
Отметить уведомление как прочитанное

**Headers:**
- `Authorization: Bearer <token>`

**Response 200:**
```json
{
  "message": "notification marked as read"
}
```

---

### PUT /notifications/read-all
Отметить все уведомления как прочитанные

**Headers:**
- `Authorization: Bearer <token>`

**Response 200:**
```json
{
  "message": "all notifications marked as read",
  "count": 5
}
```

---

### DELETE /notifications/:notification_id
Удалить уведомление

**Headers:**
- `Authorization: Bearer <token>`

**Response 200:**
```json
{
  "message": "notification deleted successfully"
}
```

---

## Dashboard (/me)

### GET /me
Получить сводную информацию о пользователе (уведомления, сообщения, новые симпатии)

**Headers:**
- `Authorization: Bearer <token>`

**Response 200:**
```json
{
  "user": {
    "id": 1,
    "vk_id": 123456,
    "is_online": true
  },
  "profile": {
    "id": 1,
    "display_name": "Иван",
    "is_onboarding_complete": true
  },
  "counters": {
    "unread_notifications": 3,
    "unread_messages": 5,
    "new_likes": 2,
    "new_matches": 1
  },
  "recent_activity": {
    "last_match": {
      "match_id": 3,
      "user": {
        "id": 5,
        "display_name": "Анна",
        "age": 24
      },
      "matched_at": "2024-12-04T12:00:00Z"
    },
    "last_message": {
      "match_id": 2,
      "sender": {
        "id": 4,
        "display_name": "Мария"
      },
      "content": "Привет!",
      "created_at": "2024-12-04T11:45:00Z"
    }
  }
}
```

---

## Error Responses

### 400 Bad Request
```json
{
  "error": "invalid request body"
}
```

### 401 Unauthorized
```json
{
  "error": "missing authorization header"
}
```

### 403 Forbidden
```json
{
  "error": "forbidden"
}
```

### 404 Not Found
```json
{
  "error": "resource not found"
}
```

### 409 Conflict
```json
{
  "error": "resource already exists"
}
```

### 500 Internal Server Error
```json
{
  "error": "internal server error"
}
```

---

## Polling Strategy

Для получения обновлений в реальном времени используется polling:

### Чаты (Messages)
- Каждые **5 секунд** запрашивать `GET /messages/:match_id?since=<last_message_timestamp>`
- Получать только новые сообщения после последнего известного

### Список чатов (Conversations)
- Каждые **10 секунд** запрашивать `GET /messages/conversations`
- Обновлять счетчики непрочитанных сообщений

### Dashboard (/me)
- Каждые **15 секунд** запрашивать `GET /me`
- Обновлять счетчики уведомлений, новых лайков, матчей

### Уведомления
- Каждые **20 секунд** запрашивать `GET /notifications?unread_only=true`
- Показывать новые уведомления пользователю

**Оптимизация:** При открытии экрана чата - polling каждые 5 сек, при свертывании - каждые 15 сек

---

## Notes

1. Все endpoints (кроме `/auth/vk`) требуют JWT токен в заголовке `Authorization: Bearer <token>`
2. Все timestamps в формате ISO 8601 (UTC)
3. Пагинация: используйте `limit` и `offset` query параметры
4. Расстояние `distance_km` рассчитывается от координат текущего пользователя
5. Используется polling для получения обновлений (см. раздел "Polling Strategy")
6. После успешного свайпа с `is_match: true` создается уведомление обоим пользователям
7. При сбросе дизлайков (`/feed/reset-dislikes`) все дизлайки удаляются, лента обновляется
8. Возраст пользователя рассчитывается автоматически из `birth_date`
9. Параметр `since` в запросах позволяет получать только новые данные после указанного timestamp
