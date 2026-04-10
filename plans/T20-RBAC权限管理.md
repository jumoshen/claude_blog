# T20: 管理后台 RBAC 权限管理

## 需求分析
实现基于角色的访问控制（Role-Based Access Control）。

## 实现方案

### 1. 数据库设计
```sql
-- 管理员角色表
CREATE TABLE IF NOT EXISTS `admin_roles` (
  `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(50) NOT NULL COMMENT '角色名称',
  `permissions` TEXT COMMENT '权限列表(JSON)',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员角色表';

-- 管理员表新增角色字段
ALTER TABLE users ADD COLUMN role_id BIGINT UNSIGNED DEFAULT NULL COMMENT '角色ID';
```

### 2. 权限定义
- `post:create` - 创建文章
- `post:edit` - 编辑文章
- `post:delete` - 删除文章
- `post:publish` - 发布/下架文章
- `category:*` - 分类管理
- `tag:*` - 标签管理
- `user:*` - 用户管理
- `log:*` - 日志查看

### 3. 中间件实现
- JWT token 中包含角色和权限信息
- 路由中间件检查权限

## 单元测试范围
- 权限检查逻辑测试
