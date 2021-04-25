# go-camp homeworks / Go训练营作业

## Homework 1
运行作业并查看抛出的error，请
```
cd homework1
```
然后，运行
```
go run homework1.go
```

与作业最相关的代码在此处：
https://github.com/ChenxiSu/go-camp/blob/main/homework1/dao/user.go#L28
主要是使用了github.com/pkg/errors package来打包错误传递过程中的堆栈信息，以便于调用者拿到错误之后对错误的root cause进行定位和处理。
