CREATE TABLE IF NOT EXISTS messages (
    id TEXT PRIMARY KEY,
    sender_id TEXT NOT NULL,
    receiver_id TEXT,
    group_id TEXT,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES users(id),
    FOREIGN KEY (receiver_id) REFERENCES users(id),
    FOREIGN KEY (group_id) REFERENCES groups(id),
    CHECK ((receiver_id IS NOT NULL AND group_id IS NULL) OR (receiver_id IS NULL AND group_id IS NOT NULL))
);