# tshell

基于coscli集成cdn功能, 实现无登陆操作cos， cdn功能

## 介绍

腾讯云（cos, cdn）命令行工具，使用 Go 编写，部署方便，且支持跨桶操作。


```
./tshell cdn -h
cdn query toolkit

Usage:
  tshell cdn [flags]
  tshell cdn [command]

Available Commands:
  domain      cdn域名查询
  history     cdn操作历史查询
  prefetch    cdn缓存预热
  refresh     刷新CDN缓存,可访问到最新资源
  top         TOP 数据查询
```


