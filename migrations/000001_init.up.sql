CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() UNIQUE,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    allow_comments BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4() UNIQUE,
    post_id INT REFERENCES posts(id) ON DELETE CASCADE,
    parent_id INT REFERENCES comments(id) ON DELETE CASCADE,
    content TEXT NOT NULL CHECK (char_length(content) <= 2000),
    created_at TIMESTAMP DEFAULT now()
);