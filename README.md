fetchChinese
===

遍历文件目录，获取中文字符串信息

**使用示例**:

``` shell
./getCh -action fetch -dir ./targerDir -filter *.cs* > output.txt
```

**输出结果**：
``` shell
1, filePath:../Index.cshtml,row:64,col:63 length:12,word:申请查询
2, filePath:../Create.cshtml,row:66,col:65 length:12,word:申请创建
```

**ToDo**：
修改输出文件后，反向需该替代中文字符（action:reverse）