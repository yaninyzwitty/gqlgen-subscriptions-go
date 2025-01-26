CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY,      -- Corresponds to User.id
    name TEXT,                  -- Corresponds to User.name
    email TEXT                  -- Corresponds to User.email
);

CREATE TABLE IF NOT EXISTS users_outbox (
    user_id BIGINT,                      -- Partition key: groups events by user
    event_id UUID,                       -- Clustering key: uniquely identifies events
    published BOOLEAN,                   -- Clustering key: helps sort by published status
    event_type TEXT,                     -- Type of event (e.g., "USER_CREATED", "USER_UPDATED")
    payload TEXT,                        -- JSON payload with event details
    created_at TIMESTAMP,                -- Timestamp when the event was recorded
    PRIMARY KEY ((user_id), published, event_id)
) WITH CLUSTERING ORDER BY (published ASC, event_id DESC); -- Sorts first by `published`, then by `event_id`



CREATE TABLE IF NOT EXISTS chat_rooms (
    id BIGINT PRIMARY KEY,               -- Corresponds to Room.id
    name TEXT,                           -- Corresponds to Room.name
    created_at TIMESTAMP,                -- Corresponds to Room.createdAt
    participants SET<BIGINT>             -- Stores IDs of users (User.id)
);

CREATE TABLE IF NOT EXISTS chat_rooms_outbox (
    room_id BIGINT,                      -- Partition key: groups events by chat room
    event_id UUID,                       -- Clustering key: uniquely identifies events
    published BOOLEAN,                   -- Clustering key: helps sort by published status
    event_type TEXT,                     -- Type of event (e.g., "ROOM_CREATED", "ROOM_UPDATED")
    payload TEXT,                        -- JSON payload with event details
    created_at TIMESTAMP,                -- Timestamp when the event was recorded
    PRIMARY KEY ((room_id), published, event_id)
) WITH CLUSTERING ORDER BY (published ASC, event_id DESC); -- Sorts first by `published`, then by `event_id`






CREATE TABLE IF NOT EXISTS message_by_room (
    room_id BIGINT,                      -- Corresponds to Message.roomId
    chat_id BIGINT,                      -- Corresponds to Message.id
    sender_id BIGINT,                    -- Corresponds to User.id (Message.sender)
    text TEXT,                           -- Corresponds to Message.content
    timestamp TIMESTAMP,                 -- Corresponds to Message.timestamp
    PRIMARY KEY ((room_id), chat_id)     -- Allows fetching messages by room
) WITH CLUSTERING ORDER BY (chat_id DESC);

CREATE TABLE IF NOT EXISTS message_by_room_outbox (
    room_id BIGINT,                      -- Partition key: groups events by room
    published BOOLEAN,                   -- Clustering key: helps sort by published status
    event_id UUID,                       -- Clustering key: uniquely identifies events
    chat_id BIGINT,                      -- Corresponds to `message_by_room.chat_id`
    event_type TEXT,                     -- Type of event (e.g., "MESSAGE_SENT")
    payload TEXT,                        -- JSON payload with event details
    created_at TIMESTAMP,                -- Timestamp when the event was recorded
    PRIMARY KEY ((room_id), event_id)
) WITH CLUSTERING ORDER BY (room_id ASC, event_id DESC); -- Sorts by `published` and `event_id`


CREATE CUSTOM INDEX message_by_room_outbox_published ON message_by_room_outbox (published) using 'StorageAttachedIndex';
