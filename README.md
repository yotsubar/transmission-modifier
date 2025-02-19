# Transmission Tracker 修改器

一个用于批量修改 Transmission 中种子 Tracker 地址的命令行工具。

## 功能特点

- 支持使用正则表达式匹配并替换 Tracker 地址
- 支持通过配置文件或命令行参数进行配置
- 支持批量处理多个种子
- 提供详细的操作日志和统计信息

## 安装

1. 确保已安装 Go 1.21 或更高版本
2. 克隆仓库并进入项目目录
3. 执行以下命令安装依赖并编译：

```bash
go mod download
go build
```

## 配置

### 配置文件

可以通过 YAML 格式的配置文件进行配置。创建一个配置文件（例如 `config.yaml`），参考以下示例：

```yaml
# 正则表达式，用于匹配要替换的tracker地址
# 注意：被双引号括起时，反斜杠要进行转义
pattern: "example\\.com"
# 或
#pattern: example\.com

# 新的tracker地址
tracker: "https://new.example.com/announce"

# Transmission RPC 配置
host: "localhost"
port: 9091
username: "your-username"
password: "your-password"
```

### 命令行参数

也可以通过命令行参数进行配置，支持的参数如下：

- `-pattern`: 正则表达式，用于匹配要替换的 tracker 地址
- `-tracker`: 新的 tracker 地址
- `-host`: Transmission RPC 主机地址（默认：localhost）
- `-port`: Transmission RPC 端口（默认：9091）
- `-username`: Transmission RPC 用户名
- `-password`: Transmission RPC 密码
- `-config`: 配置文件路径

注意：命令行参数的优先级高于配置文件。

## 使用示例

### 使用配置文件

```bash
./transmission-modifier -config config.yaml
```

### 使用命令行参数

```bash
./transmission-modifier -pattern "old-tracker\.com" -tracker "https://new-tracker.com/announce" -username "admin" -password "123456"
```

### 混合使用

可以同时使用配置文件和命令行参数，命令行参数会覆盖配置文件中的相应设置：

```bash
./transmission-modifier -config config.yaml -tracker "https://another-tracker.com/announce"
```

## 输出说明

程序运行时会显示详细的操作日志，包括：

- 匹配到的 Tracker 地址
- 更新操作的结果
- 每个种子的处理状态
- 最终的统计信息（成功和失败数量）

示例输出：

```console
[种子名称] 找到匹配的tracker：http://old-tracker.com/announce
[种子名称] 已替换为新的tracker：https://new-tracker.com/announce

处理完成！成功：1，失败：0
```

## 注意事项

1. 确保 Transmission 的 RPC 接口已启用
2. 正确配置 Transmission 的访问凭据
3. 正则表达式中的特殊字符需要正确转义
4. 建议在进行批量修改前先备份重要数据
