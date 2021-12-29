# vue 动态路径

如果我们使用动态路径的时候我们需要更改哪些配置

## 在vue2.x 的情况下


1. 添加prefix npm 包

```bash
# yarn 

yarn add prefix-uri

# npm 

npm i prefix-uri
```

2. 更改vue 配置`vue.config.js`

```editorconfig
module.exports = {
    // 选项...
    publicPath: './',
    assetsDir: "static",
}
```
3. 在index.html 添加base 标签

```html
    <base href="/" data-inject-target="BASE_URL">
```


4. 跳转路径的时候使用相对路径

例如:

```bash
# 跳转路由 使用相对路径不要使用绝对路径

import prefixUrl from 'prefix-uri/prefix-url'
    goAbout() {
      var baseUrl = window.location.origin + prefixUrl('/about')
      console.log("baseUrl: ", baseUrl)
      this.$router.push(baseUrl)
    }
    
 # 使用   axios 请求
 import prefixUrl from 'prefix-uri/prefix-url'

     health() {
      var baseUrl = window.location.origin + prefixUrl('/')
      this.baseUrl = baseUrl
      console.log("baseUrl: ", baseUrl)
      this.axios.get(baseUrl).then((response) => {
        console.log(response.code)
        this.response = response
      }).catch((error) => {
        console.log(error)
        this.response = error
      })
    }
```