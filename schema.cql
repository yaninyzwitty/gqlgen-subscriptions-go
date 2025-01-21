CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY,      -- Corresponds to User.id
    name TEXT,                  -- Corresponds to User.name
    email TEXT                  -- Corresponds to User.email
);

CREATE TABLE IF NOT EXISTS chat_rooms (
    id BIGINT PRIMARY KEY,               -- Corresponds to Room.id
    name TEXT,                           -- Corresponds to Room.name
    created_at TIMESTAMP,                -- Corresponds to Room.createdAt
    participants SET<BIGINT>             -- Stores IDs of users (User.id)
);


CREATE TABLE IF NOT EXISTS message_by_room (
    room_id BIGINT,                      -- Corresponds to Message.roomId
    chat_id BIGINT,                      -- Corresponds to Message.id
    sender_id BIGINT,                    -- Corresponds to User.id (Message.sender)
    text TEXT,                           -- Corresponds to Message.content
    timestamp TIMESTAMP,                 -- Corresponds to Message.timestamp
    PRIMARY KEY ((room_id), chat_id)     -- Allows fetching messages by room
) WITH CLUSTERING ORDER BY (chat_id DESC);
