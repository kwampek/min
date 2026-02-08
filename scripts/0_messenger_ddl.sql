CREATE SCHEMA IF NOT EXISTS Messenger;

CREATE TABLE IF NOT EXISTS Messenger.MediaFiles (
    media_id SERIAL PRIMARY KEY,
    type INTEGER NOT NULL,
    file_url VARCHAR(200),
    created_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS Messenger.Users (
    user_id SERIAL PRIMARY KEY,
    login VARCHAR(200) UNIQUE NOT NULL,
    password VARCHAR(200) NOT NULL,
    phone_number VARCHAR(20),
    email VARCHAR(50),
    birthday TIMESTAMP,
    sex BOOLEAN,
    avatar_link VARCHAR(200),
    created_at TIMESTAMP NOT NULL,
    search_privacy BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS Messenger.Chats (
    chat_id SERIAL PRIMARY KEY,
    type INTEGER NOT NULL,
    title VARCHAR(200) NOT NULL,
    description VARCHAR(200),
    avatar_link VARCHAR(200),
    creator_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    search_privacy BOOLEAN DEFAULT true,

    FOREIGN KEY (creator_id) REFERENCES Messenger.Users(user_id)
);

CREATE TABLE IF NOT EXISTS Messenger.ChatMembers (
    chat_member_id SERIAL PRIMARY KEY,
    chat_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    member_from TIMESTAMP NOT NULL,
    role VARCHAR(200),
    access_rights BOOLEAN DEFAULT true,

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
    folder_name VARCHAR(200),

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

CREATE TABLE IF NOT EXISTS Messenger.Tokens (
    user_id INTEGER NOT NULL REFERENCES Messenger.Users(user_id),
    token TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
