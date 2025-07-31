-- Create todos table
CREATE TABLE IF NOT EXISTS todos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    status TEXT DEFAULT 'incomplete' CHECK (status IN ('incomplete', 'complete')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create todo_category junction table for many-to-many relationship
CREATE TABLE IF NOT EXISTS todo_category (
    todo_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    PRIMARY KEY (todo_id, category_id),
    FOREIGN KEY (todo_id) REFERENCES todos(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

-- Insert sample categories
INSERT OR IGNORE INTO categories (title) VALUES 
('Work'),
('Personal'),
('Shopping'),
('Health');

-- Insert sample todos
INSERT OR IGNORE INTO todos (title, status) VALUES 
('Complete project documentation', 'incomplete'),
('Buy groceries', 'incomplete'),
('Go for a run', 'complete'),
('Read a book', 'incomplete');

-- Associate todos with categories
INSERT OR IGNORE INTO todo_category (todo_id, category_id) VALUES 
(1, 1), -- Complete project documentation -> Work
(2, 3), -- Buy groceries -> Shopping
(3, 4), -- Go for a run -> Health
(4, 2); -- Read a book -> Personal