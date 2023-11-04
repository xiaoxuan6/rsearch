# rsearch

用于搜索私有库中提交的数据

<img src="https://x.imgs.ovh/x/2023/09/04/64f531af6c5c1.png"  style="width: 50%;" />

> [!NOTE]
> 该项目自用，如果想使用该工具。
>
> 1、 fork 该项目
> 
> 2、 修改 `common/constants.go` 中的 `Owner` 和 `Repo` 为自己的 `github` 账号和仓库
> 
> 3、 `go install github.com/github账号/rsearch@latest`

# Install

```bash
go install github.com/xiaoxuan6/rsearch@latest
```

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
rsearch tags            # 查看所有标签
```
