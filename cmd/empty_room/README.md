# 空教室模块

**此模块用于研究生获取空教室信息**  create by 李梓玄 2023.10.17

此模块基于kitex框架(微服务框架):https://github.com/cloudwego/kitex

此模块基于jwch库(教务处的信息爬取库):https://github.com/west2-online/jwch

## 接口文档

TODOOOOOOOOOOOOOOOOOOO!

## 接口的传入参数与返回值

参阅项目根目录下的`idl/empty_room.thrift`文件

如需修改接口的传入参数与返回值，参照thrift语法修改`empty_room.thrift`后，在`cmd/empty_room`执行

```shell
# 更多信息参照 cmd/empty_room/Makefile
make gen
```

其中请求中的参数`building`需要传入一个教学楼Map的key，代码如下

```go
var (
	buildingMap = map[string]string{
		"x3":  "公共教学楼西3",
		"x2":  "公共教学楼西2",
		"x1":  "公共教学楼西1",
		"zl":  "公共教学楼中楼",
		"d1":  "公共教学楼东1",
		"d2":  "公共教学楼东2",
		"d3":  "公共教学楼东3",
		"wkl": "公共教学楼文科楼",
		"wk":  "公共教学楼文科楼",
	}
)

```

该Map定义于jwch库中: https://github.com/west2-online/jwch

## 接口逻辑

**因为研究生的教学管理系统没有获取空教室功能，所以本接口基于一个默认的本科生账号密码去本科教学管理系统获取**

此模块暂时只能获取旗山校区的空教室信息

### `GetEmptyRoom`

1. 验证token   handler
2. 当未传入账号密码时,使用默认账号密码 handler
3. 根据request信息拼接key,从`redis`中寻找缓存
4. 如果命中缓存，则返回空教室信息
5. 如果未命中缓存,调用`jwch`库进行登录、获取空教室信息
6. 返回空教室信息，同时开启一个Goroutine 将空教室信息写入

