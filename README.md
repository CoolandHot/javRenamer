# javRenamer
rename multiple jav files into a more manageable categorical way

# usage
drap and drop files into `javRename拖放文件.bat`

# modify
## proxy
change the proxy on line 12.

PROXY = {"http": "socks5://127.0.0.1:1099", "https": "socks5://127.0.0.1:1099"}

Delete it if you don't need one.

## rename rules
**default rule**:
actress-[avid]-[title]-[publishDate].suffix

Change order on line 10.

The first part is not covered by [ ], which could be changed  on line 53.