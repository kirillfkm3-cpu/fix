CREATE TABLE IF NOT EXISTS notifications (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    type TEXT NOT NULL, -- follow_request, group_invite, group_join_request, event_created
    from_user_id TEXT,
    group_id TEXT,
    event_id TEXT,
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (from_user_id) REFERENCES users(id),
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (event_id) REFERENCES events(id)
);