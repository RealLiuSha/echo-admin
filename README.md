<div align=center>
<img src="https://images.liusha.me/common/logo.png" width=200" height="200" />
</div>

<h1 align="center">Echo-Admin</h1>

<div align="center">
 基于 Echo + Gorm + Casbin + Uber-FX 实现的 RBAC 权限管理脚手架，致力于提供一套尽可能轻量且优雅的中后台解决方案。

<br/>
<br/>

<div align=center>
<img src="https://img.shields.io/badge/golang-1.16-blue"/>
<img src="https://img.shields.io/badge/echo-4.3.0-lightBlue"/>
<img src="https://img.shields.io/badge/gorm-1.21.9-red"/>
<img src="https://img.shields.io/badge/casbin-2.30.1-brightgreen"/>
<img src="https://img.shields.io/badge/vue-2.6.12-green"/>
</div>

<br/>
</div>

[English](https://github.com/RealLiuSha/echo-admin/blob/main/README.en.md) | 简体中文

## 特性
* 遵循 `RESTful API` 设计规范
* 基于 `Echo API` 框架，提供了丰富的中间件支持 (JWT 认证、鉴权、请求级事务、访问日志、跨域等)
* 基于 `Casbin` 的 `RBAC` 访问控制模型
* 基于 `Gorm V2` 的数据库存储，可扩展多种类型数据库
* 基于 `uber/fx` 实现依赖注入
* 支持 `Swagger` 文档 (基于 `swaggo`)
* 配置化、模块化

## 简介

`echo-admin` 是基于 vue 和 go 整合了优秀的开源框架和工具实现的中后台管理系统，集成了用户认证、角色管理、动态菜单和权限控制，让任何可能的使用者把时间专注在业务开发上。

[在线预览](https://admin.liusha.me)
 
```
# 只读账号
用户名: test
密码: 123123
```

[Swagger 文档](https://admin.liusha.me/swagger/index.html)

[前端项目源码](https://github.com/RealLiuSha/echo-admin-ui)

## 使用说明

欢迎 PR 和 Issue，理想情况下，我都会尽快处理和回复，感谢你关注甚至使用 `echo-admin`。

**开发语言推荐版本**

```
node >= 12.22.1
golang >= 1.16 
```

**下载代码**

```
git clone https://github.com/RealLiuSha/echo-admin
```

**生成文档**

当你完善了项目中的 swagger 文档需要重新生成，执行以下指令

```
make swagger
```

**项目初始化**

`echo-admin` 通过 `makefile` 预设了一些指令，详情可自行[查阅](https://github.com/RealLiuSha/echo-admin/blob/main/Makefile)

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

- [ ] 异步任务
- [ ] 实现日志审计
- [ ] 全配置化的工作流
- [ ] 个人中心
- [ ] 系统状态展示
- [ ] 生产级的项目质量

## 互动交流

| 微信 |
|  :---:  | 
| <img width="150" src="https://images.liusha.me/20210507/20210507183345.jpg"> 
