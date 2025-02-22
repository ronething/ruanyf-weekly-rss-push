<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [ruanyf-weekly-rss-push](#ruanyf-weekly-rss-push)
  - [usage](#usage)
  - [效果图](#%E6%95%88%E6%9E%9C%E5%9B%BE)
  - [acknowledgement](#acknowledgement)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## ruanyf-weekly-rss-push

基于 actionsflow 实现科学爱好者周刊 rss feed 推送, 目前支持 slack/wecom

### usage

- fork 此仓库
- Secrets Actions 中添加 SLACK_URL 和 WECOM_URL 相对应的值即可
- 定时策略可在 .github/workflows 目录下面进行相应修改

### 效果图

![wecom](./img/wecom.png)_

### acknowledgement

- ruanyf/weekly [issue-190](https://github.com/ruanyf/weekly/issues/2132)
- [actionsflow](https://github.com/actionsflow/actionsflow)
