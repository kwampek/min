CREATE SCHEMA IF NOT EXISTS Messenger;

CREATE TABLE IF NOT EXISTS Messenger.MediaFiles (
    media_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    preview TEXT DEFAULT NULL,
    storage_key TEXT NOT NULL,
    mime_type VARCHAR(64),
    size_bytes BIGINT,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS Messenger.Users (
    user_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    login VARCHAR(32) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20),
    email VARCHAR(254),
    birthday DATE,
    sex BOOLEAN,
    avatar_id INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    search_privacy BOOLEAN NOT NULL DEFAULT TRUE,

    FOREIGN KEY (avatar_id) REFERENCES Messenger.MediaFiles(media_id)
);

CREATE TABLE IF NOT EXISTS Messenger.Chats (
    chat_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    type SMALLINT NOT NULL,
    title VARCHAR(32),
    description TEXT,
    avatar_id INTEGER DEFAULT NULL,
    creator_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    search_privacy BOOLEAN DEFAULT true,

    FOREIGN KEY (creator_id) REFERENCES Messenger.Users(user_id),
    FOREIGN KEY (avatar_id) REFERENCES Messenger.MediaFiles(media_id)
);

CREATE TABLE IF NOT EXISTS Messenger.ChatMembers (
    chat_member_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    chat_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    custom_title VARCHAR(32) DEFAULT NULL,
    member_from TIMESTAMP NOT NULL,
    role VARCHAR(32),
    access_rights BIGINT NOT NULL DEFAULT 0,

    FOREIGN KEY (chat_id) REFERENCES Messenger.Chats(chat_id),
    FOREIGN KEY (user_id) REFERENCES Messenger.Users(user_id)
);

CREATE TABLE IF NOT EXISTS Messenger.Messages (
    message_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    message_text TEXT,
    from_chat_member_id INTEGER NOT NULL,
    media_id INTEGER,
    send_time TIMESTAMP NOT NULL DEFAULT now(),
    mstatus BOOLEAN DEFAULT false,

    FOREIGN KEY (from_chat_member_id) REFERENCES Messenger.ChatMembers(chat_member_id),
    FOREIGN KEY (media_id) REFERENCES Messenger.MediaFiles(media_id)
);

CREATE TABLE IF NOT EXISTS Messenger.Calls (
    call_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP
);

CREATE TABLE IF NOT EXISTS Messenger.CallMembers (
    user_id INTEGER NOT NULL,
    call_id INTEGER NOT NULL,
    member_from TIMESTAMP NOT NULL,
    member_to TIMESTAMP,

    PRIMARY KEY (user_id, call_id),
    FOREIGN KEY (user_id) REFERENCES Messenger.Users(user_id),
    FOREIGN KEY (call_id) REFERENCES Messenger.Calls(call_id)
);

CREATE TABLE IF NOT EXISTS Messenger.Folders (
    folder_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INTEGER NOT NULL,
    folder_name VARCHAR(32),

    FOREIGN KEY (user_id) REFERENCES Messenger.Users(user_id)
);

CREATE TABLE IF NOT EXISTS Messenger.FolderChats (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    folder_id INTEGER NOT NULL,
    chat_id INTEGER NOT NULL,
    UNIQUE(folder_id, chat_id),

    FOREIGN KEY (folder_id) REFERENCES Messenger.Folders(folder_id),
    FOREIGN KEY (chat_id) REFERENCES Messenger.Chats(chat_id)
);

CREATE TABLE IF NOT EXISTS Messenger.Sessions (
    session_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INTEGER NOT NULL,
    token_hash TEXT NOT NULL,
    device_name VARCHAR(128),
    ip_address VARCHAR(64),
    created_at TIMESTAMP DEFAULT now(),
    last_activity TIMESTAMP DEFAULT now(),
    expires_at TIMESTAMP NOT NULL,
    revoked BOOLEAN DEFAULT FALSE,

    FOREIGN KEY(user_id) REFERENCES Messenger.Users(user_id)
);

CREATE UNIQUE INDEX idx_users_login_unique ON Messenger.Users(LOWER(login));
CREATE UNIQUE INDEX idx_users_email_unique ON Messenger.Users(LOWER(email)) WHERE email IS NOT NULL AND email != '';
CREATE UNIQUE INDEX idx_users_phone_unique ON Messenger.Users(phone_number) WHERE phone_number IS NOT NULL AND phone_number != '';
CREATE UNIQUE INDEX idx_chatmembers_chat_user ON Messenger.ChatMembers(chat_id, user_id);
CREATE INDEX idx_chatmembers_user ON Messenger.ChatMembers(user_id);
CREATE INDEX idx_chatmembers_chat ON Messenger.ChatMembers(chat_id);
-- CREATE INDEX idx_messages_member_time ON Messenger.Messages(from_chat_member_id, send_time DESC);
CREATE INDEX idx_chatmembers_chat_member ON Messenger.ChatMembers(chat_id, chat_member_id);
CREATE INDEX idx_folderchats_chat ON Messenger.FolderChats(chat_id);
CREATE INDEX idx_folders_user ON Messenger.Folders(user_id);
CREATE UNIQUE INDEX idx_sessions_token ON Messenger.Sessions(token_hash);
CREATE INDEX idx_sessions_user ON Messenger.Sessions(user_id);
