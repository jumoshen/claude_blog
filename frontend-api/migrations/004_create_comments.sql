-- Migration: Create comments table
-- Created: 2026-04-02

CREATE TABLE IF NOT EXISTS `comments` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `created_at` DATETIME(3),
    `updated_at` DATETIME(3),
    `deleted_at` DATETIME(3),
    `post_slug` VARCHAR(200),
    `user_id` BIGINT,
    `nickname` VARCHAR(50),
    `content` TEXT,
    `ip` VARCHAR(50),
    `device_id` VARCHAR(64),
    `user_agent` VARCHAR(500),
    `status` INT DEFAULT 1 COMMENT '1=正常 0=待审核 -1=违规',
    INDEX `idx_post_slug` (`post_slug`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_ip` (`ip`),
    INDEX `idx_device_id` (`device_id`),
    INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
