# Changelog

All notable changes to this project will be documented in this file.

## [1.1.0] - 2026-03-25

### Added
- 访问日志功能 (Visit Log)
  - 新增 Visit 模型记录文章访问
  - GetPost 接口新增 RecordVisit 异步记录访问日志
  - 新增 CreateVisit Repository 方法

### Fixed
- 修复 Dockerfile templates 文件夹不存在的问题

### Changed
- 简化 Dockerfile，移除不存在的 static 和 templates 文件夹
- 优化容器网络配置

## [1.0.0] - 2026-03-24

### Added
- 基础博客 API 服务
- 文章管理 (CRUD)
- 用户认证 (JWT + GitHub OAuth)
- 标签和分类管理
- 归档功能
