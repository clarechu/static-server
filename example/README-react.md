# vue 动态路径

如果我们使用动态路径的时候我们需要更改哪些配置

## 在react 的情况下

...待开发


1. 在`index`页面添加base 标签。

在`src/pages/document.ejs`中添加以下内容

```html
  <head>
    <base href="/" data-inject-target="BASE_URL">
  </head>  
```

2. 修改 react中webpack 打包路径

`config/config.ts` 中添加

```ts
export default defineConfig({
    publicPath: "./",
})
```