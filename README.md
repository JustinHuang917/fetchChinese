fetchChinese
===

遍历文件目录，获取中文字符串信息

**使用示例**:

``` shell
./fetchChinese -action fetch -dir ./targerDir -filter *.cs* > output.txt
```

**输出结果**：
``` shell
1, filePath:../Index.cshtml,row:64,col:63 length:12,word:申请查询
2, filePath:../Create.cshtml,row:66,col:65 length:12,word:申请创建
```

**ToDo**：
1. 修改输出文件后作为输入，反向替代中文字符（action:reverse）
2. 从svn/git等获取项目文件后，探测文件内容
