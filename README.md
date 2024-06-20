# OneAPI 数据库迁移工具

这是一个`songquanpeng/one-api`数据迁移到`MartialBE/one-api`的工具。支持的数据库包括 MySQL、Postgres 和 SQLite。

## 功能

- 从旧数据库读取数据并插入到新数据库中
- 支持的表包括：`abilities`、`channels`、`logs`、`options`、`redemptions`、`tokens`、`users`
- 自动检测数据库驱动
- 支持通过环境变量配置数据库连接

## 使用方法

### 环境变量配置

首先，需要通过环境变量配置旧数据库和新数据库的连接字符串：

- `ONEAPI_OLD_SQL_DSN`: songquanpeng/one-api数据库的连接字符串
- `ONEAPI_NEW_SQL_DSN`: MartialBE/one-api数据库的连接字符串

例如，对于 MySQL 数据库，可以设置以下环境变量：

```bash
export ONEAPI_OLD_SQL_DSN="mysql://user:password@tcp(oldhost:3306)/olddb"
export ONEAPI_NEW_SQL_DSN="mysql://user:password@tcp(newhost:3306)/newdb"
```
例如，对于 PostgreSQL 数据库，可以设置以下环境变量：

```bash
export ONEAPI_OLD_SQL_DSN="postgres://user:password@tcp(oldhost:3306)/olddb"
export ONEAPI_NEW_SQL_DSN="postgres://user:password@tcp(newhost:3306)/newdb"
```

例如，对于 Sqlite 数据库，可以设置以下环境变量：

```bash
export ONEAPI_OLD_SQL_DSN="旧数据库数据文件绝对地址"  // 比如 /olddb.sqlite3 视你项目具体实施情况而定
export ONEAPI_NEW_SQL_DSN="新数据库数据文件绝对地址"  // 比如 /newdb.sqlite3 视你项目具体实施情况而定
```

### 使用 Docker

我们提供了一个 Docker 镜像 `zerodeng/oneapi-dbtransfer:latest`，可以直接使用该镜像进行数据迁移。

#### 运行 Docker 容器

```bash
docker run --rm \
  -e ONEAPI_OLD_SQL_DSN="mysql://user:password@tcp(oldhost:3306)/olddb" \
  -e ONEAPI_NEW_SQL_DSN="mysql://user:password@tcp(newhost:3306)/newdb" \
  zerodeng/oneapi-dbtransfer:latest
```

### 使用 Docker Compose

我们也提供了一个 Docker Compose 配置文件，您可以使用 Docker Compose 来运行数据库迁移工具。

#### 创建 `docker-compose.yml` 文件

将以下内容复制到 `docker-compose.yml` 文件中：

```yaml
version: '3.8'

services:
  dbtransfer:
    image: zerodeng/oneapi-dbtransfer:latest
    environment:
      ONEAPI_OLD_SQL_DSN: "mysql://user:password@tcp(oldhost:3306)/olddb"
      ONEAPI_NEW_SQL_DSN: "mysql://user:password@tcp(newhost:3306)/newdb"
```

#### 启动服务

在终端中导航到 `docker-compose.yml` 文件所在的目录，然后运行以下命令启动服务：

```bash
docker-compose up
```

#### 停止服务

数据迁移完成后，您可以使用以下命令停止并移除服务：

```bash
docker-compose down
```

### 使用二进制程序

我们也提供了二进制程序，您可以从[发布页面](https://github.com/ZeroDeng01/one-api-db-transfer/releases)下载并运行。

#### 下载二进制程序

请下载适合您操作系统的版本，并赋予执行权限。

```bash
chmod +x oneapi-dbtransfer-对应系统版本
```

#### 运行二进制程序
*以下示例仅供，具体以对应平台为准*
##### Linux
```bash
./db-transfer-linux-amd64 user:password@tcp(oldhost:3306)/olddb user:password@tcp(oldhost:3306)/newdb
```
##### Windows
```bash
./db-transfer-windows-amd64.exe user:password@tcp(oldhost:3306)/olddb user:password@tcp(oldhost:3306)/newdb
```
##### MacOs
```bash
./db-transfer-darwin-amd64 user:password@tcp(oldhost:3306)/olddb user:password@tcp(oldhost:3306)/newdb
```
## 截图
![运行截图](http://img.qiniu.zerodeng.com/img/202406110247982.png)


## 注意事项

- 由于两个项目数据结构差异比较大，所以迁移后部分数据需要手动调整，比如部分渠道的密钥使用`|`分割，但是两个项目里面密钥填写顺序不一样。
- 确保在迁移过程中，旧数据库和新数据库的连接稳定。
- 迁移过程中会输出进度信息，请关注控制台输出以了解迁移进度。

## 常见问题

### Q: 迁移过程中遇到错误怎么办？

A: 请查看控制台输出的错误信息，根据错误信息进行调试。如果遇到无法解决的问题，请提交 issue 并附上错误信息。

### Q: 如何处理新库中缺少的字段？

A: 工具会自动检测旧库中存在但新库中缺少的字段，并在控制台输出警告信息。这些字段的数据将不会被迁移。

### Q: 如何提高迁移速度？

A: 迁移速度受限于数据库性能和网络带宽。可以尝试优化数据库配置和网络环境以提高迁移速度。

## 声明
⚠️数据无价，数据迁移操作需要您有一定的技术基础并提前对相关重要数据进行备份。本程序不对您的数据安全负责。

## 贡献

欢迎提交 issue 和 pull request 来改进本项目。

## 支持项目

如果你觉得本项目对您有帮助，可以请作者喝杯咖啡！🎉
![支持](http://img.qiniu.zerodeng.com/img/202406112038424.jpg)
