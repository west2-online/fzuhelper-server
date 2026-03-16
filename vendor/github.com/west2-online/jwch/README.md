# JWCH

This is an interface encapsulation class of Fuzhou University Academic Affairs Office implemented by Golang, which supports simulated users to conduct personal academic affairs operations.

# Docs

| Name  | Location                        |
| ----- | ------------------------------- |
| api   | [./docs/api.md](./docs/api.md)     |
| error | [./docs/error.md](./docs/error.md) |
| model | [./docs/model.md](./docs/model.md) |

# How to use

We use this repo as an 

```bash
❯ go get github.com/west2-online/jwch
```

Then we just need to modify **main.go** to test any func.

For more detail, plz visit API docs.

# Current progress

- [X] User login
- [X] Get course selections for each semester
- [X] Get marks
- [X] Get user info
- [X] Session check
- [X] Automatic code identification
- [X] Set any apis but not implement
- [X] Empty Rooms
- [X] Using Currency
- [ ] Complete all apis
- [ ] Benchmark test
- [ ] Bug check & fix
- [ ] ...

# File tree

```
.
├── README.md			// 文档
├── cookies.txt
├── errno			// 错误处理
│   ├── code.go
│   ├── default.go
│   └── errno.go
├── docs			// 文档
│   ├── api.md			// API接口
│   ├── error.md		// 错误定义
├── go.mod
├── go.sum
├── jwch			// 教务处类
│   ├── course.go		// 课程
│   ├── jwch.go			// 类主函数
│   ├── mark.go			// 成绩
│   ├── model.go		// 自定义结构体
│   ├── user.go			// 用户
│   └── xpath.go		// xpath优化函数
├── main.go
└── utils			// 通用函数
    └── utils.go
```

### Lark Docs
https://west2-online.feishu.cn/wiki/HAZUwmgkWiRq4zkBkAecu5eGn9n

### Local Action Test

we can use act to local github action test.

1. Create `act_secret_file` in root folder
2. Insert items:
```env
USERNAME_23="YourJwchAccount"
PASSWORD_23="YourJwchSecret"
```
3. use `act push --secret-file ./act_secret_file`