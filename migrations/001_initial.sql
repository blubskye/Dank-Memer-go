-- Dank Memer Bot Database Schema
-- MariaDB/MySQL compatible

-- Guild settings
CREATE TABLE IF NOT EXISTS guilds (
    id VARCHAR(20) PRIMARY KEY COMMENT 'Discord guild ID (snowflake)',
    prefix VARCHAR(32) NOT NULL DEFAULT 'pls',
    disabled_commands JSON DEFAULT '[]' COMMENT 'Array of disabled command names',
    premium BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- User coin balances
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(20) PRIMARY KEY COMMENT 'Discord user ID',
    coins BIGINT UNSIGNED DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Command cooldowns
CREATE TABLE IF NOT EXISTS cooldowns (
    user_id VARCHAR(20) NOT NULL,
    command VARCHAR(64) NOT NULL,
    expires_at BIGINT NOT NULL COMMENT 'Unix timestamp in milliseconds',
    PRIMARY KEY (user_id, command),
    INDEX idx_expires (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Blocked users and guilds
CREATE TABLE IF NOT EXISTS blocked (
    id VARCHAR(20) PRIMARY KEY COMMENT 'User ID or Guild ID',
    type ENUM('user', 'guild') NOT NULL,
    reason TEXT,
    blocked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Donators/Premium users
CREATE TABLE IF NOT EXISTS donators (
    id VARCHAR(20) PRIMARY KEY COMMENT 'Discord user ID',
    level INT UNSIGNED DEFAULT 1 COMMENT 'Donation tier',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Bot statistics (for multi-shard aggregation)
CREATE TABLE IF NOT EXISTS stats (
    id INT PRIMARY KEY DEFAULT 1,
    guilds INT UNSIGNED DEFAULT 0,
    users INT UNSIGNED DEFAULT 0,
    channels INT UNSIGNED DEFAULT 0,
    shards INT UNSIGNED DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert initial stats row
INSERT INTO stats (id, guilds, users, channels, shards) VALUES (1, 0, 0, 0, 0)
ON DUPLICATE KEY UPDATE id=id;

-- Scheduled event to cleanup expired cooldowns (runs every hour)
-- Note: This requires EVENT_SCHEDULER to be ON in MariaDB
-- Run: SET GLOBAL event_scheduler = ON;
DELIMITER //
CREATE EVENT IF NOT EXISTS cleanup_expired_cooldowns
ON SCHEDULE EVERY 1 HOUR
DO
BEGIN
    DELETE FROM cooldowns WHERE expires_at < (UNIX_TIMESTAMP() * 1000);
END //
DELIMITER ;
