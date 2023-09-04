# rsearch

用于搜索私有库中提交的数据
![raw.png](https://x.imgs.ovh/x/2023/09/04/64f531af6c5c1.png)

# Example

同步数据到本地 `sqlite`

```darcs
rsearch sync --token="xxx"
```

搜索数据

```darcs
rsearch all             # 查看所有数据
rsearch "content"       # 使用模糊查询
rsearch "" tag          # 搜索标签所有数据
rsearch "content" tag   # 搜索标签中匹配数据
```
