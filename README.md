Для запуска проекта имеется Makefile:
 Команда: make run-server для запуска сервера
 
Весь сервер помещен в контейнер для удобства.

Сервер имеет следующие эндпоинты:
-GET /users/{id}/status - вся доступная информация о пользователе
-GET /users/leaderboard - топ пользователей с самым большим балансом
-POST /users/{id}/task/complete - выполнение задания 
-POST /users/{id}/referrer - ввод реферального кода (может быть id другого пользователя)
-POST /register - для регистрации пользователя
-POST /login - для аутентификации пользователя

Протестированные эндпоинты:
-POST /register - для регистрации пользователя
![Снимок экрана 2025-06-04 224832](https://github.com/user-attachments/assets/24723ea4-c14c-4ff2-8af6-c211ea5456f9)

-POST /login - для аутентификации пользователя
![Снимок экрана 2025-06-04 225316](https://github.com/user-attachments/assets/1aab5dea-c472-4a8f-ad3a-8c5830b75132)


-POST /users/{id}/task/complete - выполнение задания 
![Снимок экрана 2025-06-04 225713](https://github.com/user-attachments/assets/37fc88a2-298b-4d78-a050-aed0a14edbad)

-GET /users/leaderboard - топ пользователей с самым большим балансом
![Снимок экрана 2025-06-04 225549](https://github.com/user-attachments/assets/e289ed24-27b6-4a3f-a2b7-a8f54688af12)

-GET /users/{id}/status - вся доступная информация о пользователе
![Снимок экрана 2025-06-04 225739](https://github.com/user-attachments/assets/addba015-0799-4081-8e1f-c3d6cde6fb40)

