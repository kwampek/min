CREATE SCHEMA IF NOT EXISTS Messenger;

CREATE TABLE IF NOT EXISTS Messenger.MediaFiles (
    media_id SERIAL PRIMARY KEY,
    type SMALLINT NOT NULL,
    file_url TEXT NOT NULL,
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
    avatar_link TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    search_privacy BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS Messenger.Chats (
    chat_id SERIAL PRIMARY KEY,
    type SMALLINT NOT NULL,
    title VARCHAR(32) NOT NULL,
    description TEXT,
    avatar_id INTEGER DEFAULT NULL,
    creator_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    search_privacy BOOLEAN DEFAULT true,

    FOREIGN KEY (creator_id) REFERENCES Messenger.Users(user_id)
);

CREATE TABLE IF NOT EXISTS Messenger.ChatMembers (
    chat_member_id SERIAL PRIMARY KEY,
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
    message_id SERIAL PRIMARY KEY,
    message_text TEXT,
    from_chat_member_id INTEGER NOT NULL,
    media_id INTEGER,
    send_time TIMESTAMP,
    mstatus BOOLEAN DEFAULT false,

    FOREIGN KEY (from_chat_member_id) REFERENCES Messenger.ChatMembers(chat_member_id),
    FOREIGN KEY (media_id) REFERENCES Messenger.MediaFiles(media_id)
);

CREATE TABLE IF NOT EXISTS Messenger.Calls (
    call_id SERIAL PRIMARY KEY,
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
    folder_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    folder_name VARCHAR(32),

    FOREIGN KEY (user_id) REFERENCES Messenger.Users(user_id)
);

CREATE TABLE IF NOT EXISTS Messenger.FolderChats (
    id SERIAL PRIMARY KEY,
    folder_id INTEGER NOT NULL,
    chat_id INTEGER NOT NULL,
    UNIQUE(folder_id, chat_id),

    FOREIGN KEY (folder_id) REFERENCES Messenger.Folders(folder_id),
    FOREIGN KEY (chat_id) REFERENCES Messenger.Chats(chat_id)
);

CREATE TABLE Messenger.Sessions (
    session_id SERIAL PRIMARY KEY,
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
