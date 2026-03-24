# 更新日志

## [进行中] 项目重构

### 环境配置
- **MySQL**: `127.0.0.1:3306` (root/123456) - 宿主机 Docker
- **Redis**: `127.0.0.1:6379` - 宿主机 Docker
- docker-compose 中不内置 MySQL/Redis，使用宿主机已有的服务

### 后端重构
- [ ] 前后端分离架构
- [ ] JWT 认证 + Redis 缓存 token 实现登录态管理与踢人功能
- [ ] 代码风格规范 + 模块拆分 (handler/service/repository)
- [ ] Swagger 文档生成 (make swagger)
- [ ] 配置外置化 (config.yaml)

### 前端重构
- [ ] Vue.js 框架
- [ ] 组件化拆分
- [ ] 用户信息本地缓存 (localStorage)

### DevOps
- [ ] 前后端分别 Dockerfile
- [ ] docker-compose 统一管理
- [ ] 支持 enable_proxy 代理功能

### 测试
- [x] Docker 部署验证通过
- [ ] 前后端联调测试
