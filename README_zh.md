<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [科学爱好者周刊推送](#%E7%A7%91%E5%AD%A6%E7%88%B1%E5%A5%BD%E8%80%85%E5%91%A8%E5%88%8A%E6%8E%A8%E9%80%81)
  - [简介](#%E7%AE%80%E4%BB%8B)
  - [特性](#%E7%89%B9%E6%80%A7)
  - [使用方法](#%E4%BD%BF%E7%94%A8%E6%96%B9%E6%B3%95)
  - [AI 内容分析](#ai-%E5%86%85%E5%AE%B9%E5%88%86%E6%9E%90)
  - [效果预览](#%E6%95%88%E6%9E%9C%E9%A2%84%E8%A7%88)
  - [致谢](#%E8%87%B4%E8%B0%A2)
  - [开源协议](#%E5%BC%80%E6%BA%90%E5%8D%8F%E8%AE%AE)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 科学爱好者周刊推送

[English](README.md) | [中文](README_zh.md)

## 简介

这是一个基于 Go 语言开发的自动化服务，用于将阮一峰的科学爱好者周刊 RSS 推送到 Slack、企业微信和飞书。该服务使用 GitHub Actions 实现自动化运行。

## 特性

- 自动监控 RSS 订阅源
- 支持推送到：
  - Slack
  - 企业微信
  - 飞书
- AI 驱动的内容分析和推荐提取
- 基于 GitHub Actions 运行
- 使用 Go 语言重写，性能更好

## 使用方法

1. Fork 此仓库到你的账户
2. 在仓库的 Settings -> Secrets and Variables -> Actions 中添加以下密钥：
   - `SLACK_URL`: Slack Webhook URL（用于 Slack 推送）
   - `WECOM_URL`: 企业微信 Webhook URL（用于企业微信推送）
   - `FEISHU_URL`: 飞书机器人 Webhook URL（用于飞书推送）
   - `ARK_BASE_URL`: （可选）Ark API 的基础 URL
   - `ARK_API_KEY`: Ark AI API 密钥（用于 AI 内容分析）
   - `ARK_MODEL_ID`: Ark API 使用的模型 ID
3. GitHub Actions 会按照预设时间自动运行检查并推送更新

## AI 内容分析

本服务包含 AI 驱动的内容分析功能：

- 从每期周刊中提取 5 个最有价值的推荐项目
- 根据实用价值、新颖性和技术趋势分析内容
- 将这些推荐添加到推送通知中
- 使用 Ark API（类似于 OpenAI API）

如果不想使用 AI 功能，只需不设置 `ARK_API_KEY` 环境变量即可。

## 效果预览

![wecom](./img/wecom.png)

## 致谢

- [ruanyf/weekly](https://github.com/ruanyf/weekly)
- 原始讨论: [issue-190](https://github.com/ruanyf/weekly/issues/2132)

## 开源协议

MIT License
