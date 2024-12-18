-- +goose Up
CREATE TABLE Cards (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status INTEGER NOT NULL,
    completedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    estimatedMins INTEGER NOT NULL DEFAULT 0, 
    trackedMins INTEGER NOT NULL DEFAULT 0,
    isActive BOOLEAN NOT NULL DEFAULT FALSE,
    projectId INTEGER NOT NULL,
    FOREIGN KEY (projectId) REFERENCES Projects(id) ON DELETE CASCADE
);

CREATE TABLE TimeEntries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    cardId INTEGER NOT NULL,
    startTime TIMESTAMP NOT NULL,
    endTime TIMESTAMP NOT NULL,
    duration INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (cardId) REFERENCES Cards(id) ON DELETE CASCADE
);

CREATE TABLE Projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE UserProfile (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    settings TEXT NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO Projects ('id', 'name') VALUES (1, 'Inbox');
INSERT INTO UserProfile ('id', 'name', 'settings') VALUES (1, 'default', '');

-- +goose Down
DROP TABLE IF EXISTS Cards;
DROP TABLE IF EXISTS TimeEntries;
DROP TABLE IF EXISTS Projects;
DROP TABLE IF EXISTS UserProfile;
