# creditget
discuz 通用论坛用户积分爬取。

可以爬去论坛内所有用户积分等各类信息，并生成 markdown 表格展示。

使用 SQLite 储存，Go 编写，多协程爬取。

## 使用
编译后，复制项目内的 config.json，修改其中的各项配置。

`{
    "start": 0, //起始 uid 左闭右开
    "end": 3709554, //结束 uid 
    "thread": 8, //使用的协程数
    "sleepTime" : 500, //某个协程爬取一次后睡眠多少毫秒
    "disucuzApiAddress": "https://www.mcbbs.net/api/mobile/index.php?version=4&module=profile&uid=", // api 地址
    "extcredits1": "人气",
    "extcredits2": "金粒",
    "extcredits3": "金锭", 
    "extcredits4": "绿宝石", 
    "extcredits5": "下界之星",
    "extcredits6": "贡献",
    "extcredits7": "爱心",
    "extcredits8": "钻石"
}`

extcredits1 等和积分名字的对应关系可来形如 https://www.mcbbs.net/api/mobile/index.php?version=4&module=profile&uid=1 底部查看

未使用的 extcredits 可直接删除这个字段。

如
`{
    "start": 0, 
    "end": 3709554, 
    "thread": 8, 
    "sleepTime" : 500, 
    "disucuzApiAddress": "https://www.mcbbs.net/api/mobile/index.php?version=4&module=profile&uid=",
    "extcredits1": "人气",
    "extcredits2": "金粒",
    "extcredits3": "金锭"
}`

形如 https://www.mcbbs.net/api/mobile/index.php?version=4&module=check 中的 totalmembers 数值为论坛总人数。

然后直接执行程序即可爬取，爬取完成后会自动退出。中途若意外退出，可按上次的进度继续爬取。

执行时带上任意参数即可生成展示表格，如 `./main a`