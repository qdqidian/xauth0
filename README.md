# xauth0
caddy auth0插件，负责在指定的url中认证授权

## auth0
完整的认证授权解决方案，除常规操作权限外，集成数据权限及字段级数据权限

- [x] 配置了auth0指令的所有请求都会进行操作权限校验(操作权限支持:匿名/必须登录/必须有指定url perm)
- [x] url配置了数据权限校验(支持位置: param/header/body)，则进行单点数据验证
- [x] 未配置数据权限校验的，可通过join auth0远程数据表做数据过滤(如查询列表时,join远程表做数据过滤)
- [x] url配置了字段级权限的，则根据当前人字段级权限，对返回的json进行过滤
- [x] 完成操作权限校验后，当前用户信息可选择注入进请求头中，带进后端服务(header_up  x-global-user-id {auth0.user.id});
当前支持{auth0.user.id}/{auth0.user.name}/{auth0.user.code}

## order

auth0排序到basicauth后面
order auth0 after basicauth