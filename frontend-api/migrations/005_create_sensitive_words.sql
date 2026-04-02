-- Migration: Create sensitive_words table
-- Created: 2026-04-02
-- Updated: 2026-04-02

CREATE TABLE IF NOT EXISTS `sensitive_words` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `word` VARCHAR(100) NOT NULL COMMENT '敏感词',
    `level` INT DEFAULT 1 COMMENT '级别：1=一般 2=严重',
    `created_at` DATETIME(3),
    `updated_at` DATETIME(3),
    `deleted_at` DATETIME(3),
    UNIQUE INDEX `idx_word` (`word`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
