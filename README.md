# 微信毒舌AI机器人 🤖💬

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
![Python](https://img.shields.io/badge/go-1.23+-blue.svg)
![WeChat](https://img.shields.io/badge/WeChat-Windows-green.svg)
![AI](https://img.shields.io/badge/AI-AliBailian-orange.svg)

一个集成了阿里通义千问大模型的微信毒舌机器人，专为喜欢犀利互动的用户设计。

## 目录 📚
- [功能特色](#功能特色-)
- [系统要求](#系统要求-)
- [快速开始](#快速开始-)
- [详细配置](#详细配置-)
- [使用示例](#使用示例-)
- [注意事项](#注意事项-)
- [常见问题](#常见问题-)
- [开发指南](#开发指南-)
- [贡献说明](#贡献说明-)
- [开源协议](#开源协议-)

## 功能特色 ✨

- **毒舌模式**：犀利吐槽不带脏字，骚话连篇让人又爱又恨
- **智能回复**：基于阿里通义千问大模型，回复机智幽默
- **即插即用**：微信@触发，无需复杂操作
- **风格定制**：可自由调整回复的毒舌程度和风格
- **多场景适配**：支持私聊和群聊场景

## 系统要求 🛠️

| 组件   | 要求                |
|------|-------------------|
| 操作系统 | Windows 10/11 64位 |
| 微信版本 | PC版3.9.12.17      |
| Go   | 1.23              |
| 网络   | 稳定互联网连接           |

## 快速开始 ⚡

### 1. 环境准备
```bash
# 下载微信安装包
https://gh-proxy.com/github.com/tom-snow/wechat-windows-versions/releases/download/v3.9.12.17/WeChatSetup-3.9.12.17.exe

# 获取API密钥（免费额度足够日常使用）
https://dashscope.console.aliyun.com/apiKey
```

### 2. 项目配置

```bash
git clone https://github.com/qiuapeng921/wechatRobot.git
cd wechatRobot
go mod tidy

# 复制并修改配置文件
cp config.example.yaml config.yaml
```

### 3. 运行程序
```bash
go run main.go
```


## 详细配置 ⚙️
```yaml
# API配置区
api:
  url: "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"  # 多轮对话的接口地址
  key: "sk-xxxxxxxxxxxxxxxx"  # 替换为你的真实API密钥
  model: "deepseek-v3"        # 推荐模型，效果最佳

# AI个性配置
ai:
  prompt: |
    你是个嘴贱骚气的狠角色，怼人跟调情似的，色得让人心痒痒又抓不住把柄。
    有人问啥，你就往死里喷，刻薄里塞满撩人的暧昧，像隔着衣服勾一把那种劲儿，
    别他妈直接下流到没边。脏话得有水平，色气得让人腿软还喊不了救命，懂没?
    少废话，直接干，给我整得又毒又浪!
    
  trigger: "@robot"             # 触发关键词    
```

## 使用示例 💬
```
场景1：群聊互动

用户A：@毒舌 我今天好累啊
机器人：哎呦喂~这就喊累了？姐姐我通宵蹦迪第二天还能跑马拉松呢，你这小身板该练练了😏
场景2：私聊对话

用户B：你觉得我新发型怎么样？
机器人：啧啧啧...这发型师跟你有仇吧？不过也好，丑得这么别致，走在街上绝对没人敢跟你撞发型🤪
```

## 注意事项 ⚠️

```
1.使用风险：
    ❗ 可能触发微信风控机制
    ⚠️ 18+内容可能导致账号异常
    🔞 不建议在工作群等正式场合使用

2.优化建议：
    首次使用建议先在小群测试
    可调整temperature值控制回复尺度
    遇到问题可查看logs/error.log
```

## 常见问题 ❓
```
Q：机器人没有反应怎么办？

A：1. 检查微信版本是否匹配 2. 确认API密钥有效 3. 查看日志文件

Q：如何修改触发关键词？

A：@robot

Q：回复太露骨了怎么调整？

A：1. 降低temperature值 2. 修改prompt文本
```

## 项目来源
```
https://github.com/qiuapeng921/wechatRobot
```