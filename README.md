# 超级简单的构建项目的工具
尚未完全开发完
```yaml
version: "1.0"
namespace: default
import:
  - path: ./a/b
    name: utils
default_command: build
build:
  desc: build command
  shell:
    - go build -o main.exe ./
```