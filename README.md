# Builder - 轻量级命令编排工具
Builder 是一款简化命令行工作流的工具，通过 YAML 配置文件定义和组织命令，支持模块化复用、命名空间隔离，让复杂的命令序列更易管理和执行。

## 核心功能
- YAML 配置驱动：用简洁的 YAML 格式定义命令，无需学习复杂语法
- 命名空间隔离：通过 namespace 划分命令集，避免命名冲突
- 跨文件复用：通过 import 导入其他配置文件的命令，实现模块化管理
- 默认命令：设置 default_command，支持无参数直接执行预设命令

## 安装方法
### Go install
```
go install github.com/Li-giegie/builder@latest
```
### 源码编译（需 Go 1.23+）
```shell
# 克隆仓库
git clone https://github.com/Li-giegie/builder.git
cd builder

# 编译可执行文件
go build -o builder main.go

# 安装到系统路径（全局可用）
# Linux/macOS
sudo mv builder /usr/local/bin/
# Windows（管理员权限，移动到 PATH 目录如 C:\Windows\System32）
```
## 快速开始
### 1. 初始化配置文件
生成默认的命令配置文件（.builder.yaml）：
```shell
builder init
```
默认配置文件内容如下：
```yaml
version: "1.0"        # 配置版本（固定值）
namespace: default    # 命令集命名空间
import:               # 导入其他配置文件（可选）
  - path: "{{file path}}"  # 目标配置文件路径
    name: "{{reNamespace}}" # 导入后的别名
default_command: hello # 无参数时默认执行的命令
hello:                 # 自定义命令
  desc: hello command  # 命令描述
  shell:               # 执行的命令列表
    - echo "hello world"
```
### 2. 查看可用命令
列出配置文件中定义的所有命令：
```shell
builder list
```
输出示例：
```
builder commands:
  hello   	 hello command
```
### 3. 执行命令
- 执行指定命令：
```shell
builder hello
```
输出：`hello world`
- 执行默认命令（由 default_command 定义）：
```shell
builder
```
等价于执行 `builder hello`
## 配置文件详解
配置文件是 Builder 的核心，用于定义命令集、依赖关系和复用规则。基本结构如下：
```yaml
version: "1.0"          # 固定版本号，当前仅支持 "1.0"
namespace: <命名空间>    # 命令集标识，用于区分不同模块的命令
import:                 # 导入外部配置（可选）
  - path: <文件路径>     # 外部配置文件的路径（相对/绝对路径）
    name: <别名>         # 导入后用于引用的名称
default_command: <命令名> # 无参数时默认执行的命令
<命令名>:               # 自定义命令
  desc: <描述>          # 命令说明（builder list 时显示）
  shell:                # 命令执行的脚本列表（按顺序执行）
    - <命令1>
    - <命令2>
```
### 字段说明
```
version	            配置文件版本，固定为 "1.0"
namespace	        命名空间名称，用于隔离不同模块的命令（如 build、deploy）
import	            导入其他配置文件中的命令集，实现跨文件复用
default_command	    当执行 builder 无参数时，默认执行的命令名称
<命令名>	            自定义命令的标识（如 start、clean）
<命令名>.desc	    命令的描述信息，用于 builder list 展示
<命令名>.shell	    命令实际执行的脚本列表，支持系统原生命令（如 echo、mkdir 等）
```
## 命令复用与跨文件调用
通过 import 导入其他配置文件的命令，并用 $<别名>.<命令名> 格式引用，实现命令的模块化复用。

### 示例：多文件协作
1. 创建 deploy.yaml（部署相关命令）：
```yaml
version: "1.0"
namespace: deploy
prepare:
  desc: 初始化部署目录
  shell:
    - mkdir -p /tmp/deploy_dir
package:
  desc: 打包应用
  shell:
    - zip -r app.zip ./dist
```
2. 在主配置 .builder.yaml 中导入并使用：
```yaml
version: "1.0"
namespace: main
import:
  - path: './deploy.yaml'  # 导入部署配置
    name: 'deploy'         # 别名设为 deploy
deploy:
  desc: 完整部署流程
  shell:
    - $deploy.prepare  # 调用 deploy 命名空间的 prepare 命令
    - $deploy.package  # 调用 deploy 命名空间的 package 命令
    - echo "部署完成"
```
3. 执行命令：
```shell
builder deploy  # 依次执行 prepare → package → 打印部署完成
```
## 命令行参数
```
Usage:
  builder [flags] commands...  # 执行指定命令
  builder [command]           # 内置命令（如 init、list）

Available Commands:
  completion  生成指定 shell 的自动补全脚本
  help        查看命令帮助信息
  init        初始化配置文件
  list        列出所有可用命令

Flags:
  -c, --config string   指定配置文件路径（默认 "./.builder.yaml"）
  -h, --help            显示帮助信息
```
## 示例场景
### 1. 项目构建与清理
```yaml
version: "1.0"
namespace: project
default_command: build  # 默认执行 build 命令

build:
  desc: 构建项目
  shell:
    - go build -o app main.go
    - echo "构建成功"

clean:
  desc: 清理构建产物
  shell:
    - rm -f app
    - echo "清理完成"

rebuild:
  desc: 重建项目（先清理再构建）
  shell:
    - $project.clean  # 调用当前命名空间的 clean 命令
    - $project.build  # 调用当前命名空间的 build 命令
```
使用：
```shell
builder        # 执行 build（默认命令）
builder clean  # 执行清理
builder rebuild  # 先清理再构建
```
### 2. 多环境部署
```yaml
version: "1.0"
namespace: env
import:
  - path: './common.yaml'
    name: 'common'  # 导入公共命令

deploy_dev:
  desc: 部署到开发环境
  shell:
    - $common.package  # 复用公共打包命令
    - scp app.zip dev@dev-server:/app

deploy_prod:
  desc: 部署到生产环境
  shell:
    - $common.package  # 复用公共打包命令
    - scp app.zip prod@prod-server:/app
```
使用：
```shell
builder deploy_dev  # 部署到开发环境
builder deploy_prod  # 部署到生产环境
```
## 注意事项
- 执行顺序：shell 列表中的命令按顺序执行，前序命令失败则终止执行
- 路径处理：配置文件中的路径以执行 builder 命令的当前目录为基准
- 循环引用：导入的配置文件不可相互引用（如 A 导入 B，B 导入 A），会导致执行失败
- 系统兼容性：shell 中的命令需适配当前操作系统（如 Windows 用 del，Linux 用 rm）
