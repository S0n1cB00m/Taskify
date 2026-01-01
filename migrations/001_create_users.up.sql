CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL, -- используем хеширование с Cost Factor = 12
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Добавляем временное создание пользователя с id=1, чтобы реализовывать boards
INSERT INTO users (id, email, username, password_hash)
VALUES (
           1,
           'test@example.com',
           'testuser',
           '$2a$12$hh.jUXyfU9tsJInwcI90iuJ/kP3.VaA40wggmtRd5Zj5jb1FQi4jG'
       )
    ON CONFLICT (id) DO NOTHING;