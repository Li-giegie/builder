# 超级简单的构建项目的工具
配置文件
```yaml
version: "1.0"        // 版本
namespace: default    // 命名空间
import:               // 导入其他文件
  - path: './b1.yaml' // 引用builder配置文件路径，这里仅作为展示，实际是真实的文件路径
    name: 'b1'        // 引用名称
default_command: hello // 默认执行命令
hello:                 // 定义的命令
  desc: hello command  // 命令描述
  shell:               // 执行的脚本
    - echo "hello world"
```
执行hello命令演示
```
builder.exe hello
或
builder.exe // default_command 定义的默认命令是hello这里可以省略
```

命令行
```
Usage:
  builder [flags] commands...
  builder [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        init a builder config file
  list        list builder commands

Flags:
  -c, --config string   config file path (default "./.builder.yaml")
  -h, --help            help for builder

Use "builder [command] --help" for more information about a command.
```
