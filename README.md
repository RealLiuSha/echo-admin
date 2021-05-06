<div align=center>
<img src="https://images.liusha.me/common/logo.png" width=200" height="200" />
</div>

<h1 align="center">Echo-Admin</h1>

<div align="center">
 基于 Echo + Gorm + Casbin + Uber-FX 实现的 RBAC 权限管理脚手架，致力于提供一套尽可能轻量且优雅的中后台解决方案。
 
<br/>
<br/>

<div align=center>
<img src="https://img.shields.io/badge/golang-1.12-blue"/>
<img src="https://img.shields.io/badge/echo-4.2.2-lightBlue"/>
<img src="https://img.shields.io/badge/gorm-1.21.8-red"/>
<img src="https://img.shields.io/badge/casbin-2.26.0-brightgreen"/>
<img src="https://img.shields.io/badge/vue-2.6.12-green"/>
</div>

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

## 基本介绍

`echo-admin` 是基于 vue 和 go 整合了优秀的开源框架和工具包实现的前后端分离的管理系统，集成了基本的用户认证、角色管理、动态菜单和权限控制，让任何可能的使用者把时间专注在业务开发上。

[在线预览](https://admin.srelab.cn)

[前端项目源码](https://github.com/RealLiuSha/echo-admin-ui)

## 使用说明

**开发语言推荐版本**

```
node >= 12.22.1
golang >= 1.12 
```

**下载代码**

```
git clone https://github.com/RealLiuSha/echo-admin
```

**生成文档**

当你完善了项目中的  文档需要重新生成，只需要执行以下指令

```
make swagger
```

**项目初始化**

`echo-admin` 通过 `makefile` 预设了一些指令，详情可自行查阅

首次启动本项目前需要相对应的修改配置文件 `config/config.yaml`， 你至少需要保证 `mysql` 和 `redis` 的相关配置正确，随后你可以通过以下指令完成表的新建和数据的初始化 

```
make migrate # 创建表
make setup # 初始化菜单数据
```

**启动**

```
make
```

## 计划任务

- [ ] 实现日志审计
- [ ] 模块化工作流
- [ ] 个人中心
- [ ] 系统状态展示
