# javRenamer

rename multiple jav files into a more manageable categorical way

# usage

drap and drop files into `javRename拖放文件.bat`

# modify

## Python

### proxy

**line 14**: PROXY = {"http": "socks5://127.0.0.1:1099", "https": "socks5://127.0.0.1:1099"}

If you don't need one, set it to None: PROXY = None

### rename rules

**line 9 default rule**:
actress-[avid]-[title]-[publishDate].suffix

**line 45**: The first part is not covered by [ ] by default, change it yourself.

### download cover

**line 10**: downimg = True or False

## Golang

### Proxy

It's fixed. No intension to change.

### rename rules

**line 33 default rule**:
actress-[avid]-[title]-[publishDate].suffix

### Select the site

variable `siteX` is controlled by command line parameter `-site_num`,

- 0: javbus
- 1: avmoo
- 2: javLibrary
