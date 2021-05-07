<div align=center>
<img src="https://images.liusha.me/common/logo.png" width=200" height="200" />
</div>

<h1 align="center">Echo-Admin</h1>

<div align="center">
 Echo + Gorm + Casbin + Uber-FX  based scaffolding for front and back separation management system
 
<br/>
<br/>

<div align=center>
<img src="https://img.shields.io/badge/golang-1.12-blue"/>
<img src="https://img.shields.io/badge/echo-4.2.2-lightBlue"/>
<img src="https://img.shields.io/badge/gorm-1.21.8-red"/>
<img src="https://img.shields.io/badge/casbin-2.26.0-brightgreen"/>
<img src="https://img.shields.io/badge/vue-2.6.12-green"/>
</div>

<br/>
</div>

English | [简体中文](https://github.com/RealLiuSha/echo-admin/blob/main/README.md)

## Feature

* Follow RESTful API design specifications
* Provides rich middleware support based on Echo Api framework (jwt, authority, request-level transaction, access log, cors...)
* Casbin based RBAC access control model
* GORM based database storage that can expand many types of databases
* Uber/fx based implements dependency injection
* Support Swagger documentation (based on swaggo)
* Configuration, modularization

## Synopsis

welcome PR and Issue.

[Online Demo](https://admin.srelab.cn)
 
```
# readonly account
username: test
password: 123123
```

## Getting started

```
node >= 12.22.1
golang >= 1.12 
```

**Use git to clone this project**

```
git clone https://github.com/RealLiuSha/echo-admin
```

**API docs generation**

```
make swagger
```

**Initialize the database and start the service** 

```
make migrate # create tables
make setup # setup menu data
make # start
```