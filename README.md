# 枝网后端

源项目地址: https://github.com/ASoulCnki/ASoulCnkiBackend

## 系统环境
- Golang 1.16

## 使用方式
- 等待自动爬取数据 或 关闭自动更新，使用 Get 请求 ```localhost:8080/api/v1/update``` 获取数据
- 查重使用 Post 请求 ```localhost:8080/api/v1/check``` 即可

## Todos
- [x] 可在生产环境下使用
- [x] 自动更新评论
- [x] 优化数据库连接
- [x] 优化自动更新
- [x] 枝江作文展
- [x] 兼容源项目接口