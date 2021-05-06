<h1 align="center">Echo-Admin</h1>

<div align="center">
 基于 Echo + Gorm + Casbin + Uber-FX 实现的 RBAC 权限管理脚手架，致力于提供一套尽可能轻量且优雅的中后台解决方案。
<br/>

[![ReportCard][reportcard-image]][reportcard-url] [![GoDoc][godoc-image]][godoc-url] [![License][license-image]][license-url]

</div>

## 特性
* 尽量遵循 `RESTful API` 设计规范 & 基于接口的编程规范
* 基于 `Echo` 框架，提供了丰富的中间件支持 (内置 JWTAuth，Casbin，GormThx&Recovery，ZapLogger)
* 基于 `Casbin` 的 RBAC 权限访问控制模型 (细化到 Api & Method)
* 基于 `Gorm 2.0` 的数据库存储
* 基于 `Uber-FX` 框架实现依赖注入
* 基于 `Zap` 实现日志输出并自动切割存档
* 基于 `JWT` 的用户认证
* 支持 `Swagger` 文档 (基于 `swaggo`)

## 快速开始

```
$ git clone https://github.com/RealLiuSha/echo-admin
$ cd echo-admin
$ make
```