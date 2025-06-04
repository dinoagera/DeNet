CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    pass_hash VARCHAR(255) NOT NULL,
    balance INTEGER NOT NULL DEFAULT 0,
    referrer_id INTEGER NULL,
    CONSTRAINT fk_referrer FOREIGN KEY (referrer_id) REFERENCES users(id)
);
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    reward INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE TABLE user_tasks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    completed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, task_id)
);
-- Для тестирование API
INSERT INTO tasks (name, description, reward) VALUES
    ('Регистрация', 'Завершите регистрацию в системе', 100),
    ('Подтверждение email', 'Подтвердите свой email адрес', 150),
    ('Заполнение профиля', 'Заполните информацию в своем профиле', 200),
    ('Первая задача', 'Выполните свою первую задачу', 250);