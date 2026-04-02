-- Migration: Add is_fake field to comments table
-- Created: 2026-04-02

ALTER TABLE `comments` ADD COLUMN `is_fake` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否假数据' AFTER `status`;
