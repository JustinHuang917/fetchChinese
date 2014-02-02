fetchChinese
===

遍历文件目录，获取中文字符串信息

**使用示例**:

``` shell
./fetchChinese -action fetch -dir ./test -filter 1.cs > output.txt
```

**输出结果**：
``` shell
0, filePath:test/1.cs,row:9,col:19 length:12,word:输出中文
```

**ToDo**：

1. 修改输出文件后作为输入，反向替代中文字符（action:reverse）

2. 从svn/git等获取项目文件后，探测文件内容
