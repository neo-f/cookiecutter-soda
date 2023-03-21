<!-- markdownlint-configure-file { "MD013": { "line_length": 180 } } -->

# {{ cookiecutter.project_slug }}服务说明

- Build Time: {build_time}
- Go Version: {go_version}

> **NOTE**

- 本服务所有接口的返回值均遵循Rcrai前端组件规范，如

```json
{
  "code": 0,
  "data": {},
  "message": "",
  "details": {}
}
```

但为使文档更加简洁，所有该文档中涉及的JSON Response均只展示data内部的数据结构(特殊情况除外).

如果有疑问请联系本服务维护者 :D
