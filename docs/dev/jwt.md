
基于 JWT + Refresh Token 的用户认证实践

HTTP 是一个无状态的协议，一次请求结束后，下次在发送服务器就不知道这个请求是谁发来的了（同一个 IP 不代表同一个用户），在 Web 应用中，用户的认证和鉴权是非常重要的一环，实践中有多种可用方案，并且各有千秋。

## 基于 Session 的会话管理

在 Web 应用发展的初期，大部分采用基于 Session 的会话管理方式，逻辑如下。

- 客户端使用用户名密码进行认证
- 服务端生成并存储 Session，将 SessionID 通过 Cookie 返回给客户端
- 客户端访问需要认证的接口时在 Cookie 中携带 SessionID
- 服务端通过 SessionID 查找 Session 并进行鉴权，返回给客户端需要的数据



![img](https://pic2.zhimg.com/80/v2-f5bea5679c0bece2575c5ad5bbefb04d_720w.webp)



基于 Session 的方式存在多种问题。

- 服务端需要存储 Session，并且由于 Session 需要经常快速查找，通常存储在内存或内存数据库中，同时在线用户较多时需要占用大量的服务器资源。
- 当需要扩展时，创建 Session 的服务器可能不是验证 Session 的服务器，所以还需要将所有 Session 单独存储并共享。
- 由于客户端使用 Cookie 存储 SessionID，在跨域场景下需要进行兼容性处理，同时这种方式也难以防范 CSRF 攻击。

## 基于 Token 的会话管理

鉴于基于 Session 的会话管理方式存在上述多个缺点，无状态的基于 Token 的会话管理方式诞生了，所谓无状态，就是服务端不再存储信息，甚至是不再存储 Session，逻辑如下。

- 客户端使用用户名密码进行认证
- 服务端验证用户名密码，通过后生成 Token 返回给客户端
- 客户端保存 Token，访问需要认证的接口时在 URL 参数或 HTTP Header 中加入 Token
- 服务端通过解码 Token 进行鉴权，返回给客户端需要的数据



![img](https://pic4.zhimg.com/80/v2-d1047db8e2e79446cae4af6d73697833_720w.webp)



基于 Token 的会话管理方式有效解决了基于 Session 的会话管理方式带来的问题。

- 服务端不需要存储和用户鉴权有关的信息，鉴权信息会被加密到 Token 中，服务端只需要读取 Token 中包含的鉴权信息即可
- 避免了共享 Session 导致的不易扩展问题
- 不需要依赖 Cookie，有效避免 Cookie 带来的 CSRF 攻击问题
- 使用 CORS 可以快速解决跨域问题

## JWT 介绍

JWT 是 JSON Web Token 的缩写，JWT 本身没有定义任何技术实现，它只是定义了一种基于 Token 的会话管理的规则，涵盖 Token 需要包含的标准内容和 Token 的生成过程。

一个 JWT Token 长这样。

```text
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NDQ1MTE3NDMsImp0aSI6IjYxYmVmNjkyLTE4M2ItNGYxYy1hZjE1LWUwMDM0MTczNzkxOSJ9.CZzB2-JI1oPRFxNMaoFz9-9cKGTYVXkOC2INMoEYNNA
```

仔细辨别会发现它由 `A.B.C` 三部分组成，这三部分依次是头部（Header）、负载（Payload）、签名（Signature），头部和负载以 JSON 形式存在，这就是 JWT 中的 JSON，三部分的内容都分别单独经过了 Base64 编码，以 `.` 拼接成一个 JWT Token。

JWT 的 Header 中存储了所使用的加密算法和 Token 类型。

```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

Payload 是负载，JWT 规范规定了一些字段，并推荐使用，开发者也可以自己指定字段和内容，例如下面的内容。

```json
{
  username: 'yage',
  email: 'sa@simpleapples.com',
  role: 'user',
  exp: 1544602234
}
```

需要注意的是，Payload的内容只经过了 Base64 编码，对客户端来说当于明文存储，所以不要放置敏感信息。

Signature 部分用来验证 JWT Token 是否被篡改，所以这部分会使用一个 Secret 将前两部分加密，逻辑如下。

```text
HMACSHA256(base64UrlEncode(header) + "." + base64UrlEncode(payload), secret)
```

## JWT 优势 & 问题

JWT 拥有基于 Token 的会话管理方式所拥有的一切优势，不依赖 Cookie，使得其可以防止 CSRF 攻击，也能在禁用 Cookie 的浏览器环境中正常运行。

而 JWT 的最大优势是服务端不再需要存储 Session，使得服务端认证鉴权业务可以方便扩展，避免存储 Session 所需要引入的 Redis 等组件，降低了系统架构复杂度。但这也是 JWT 最大的劣势，由于有效期存储在 Token 中，JWT Token 一旦签发，就会在有效期内一直可用，无法在服务端废止，当用户进行登出操作，只能依赖客户端删除掉本地存储的 JWT Token，如果需要禁用用户，单纯使用 JWT 就无法做到了。

## 基于 JWT 的实践

既然 JWT 依然存在诸多问题，甚至无法满足一些业务上的需求，但是我们依然可以基于 JWT 在实践中进行一些改进，来形成一个折中的方案，毕竟，在用户会话管理场景下，没有银弹。

前面讲的 Token，都是 Access Token，也就是访问资源接口时所需要的 Token，还有另外一种 Token，Refresh Token，通常情况下，Refresh Token 的有效期会比较长，而 Access Token 的有效期比较短，当 Access Token 由于过期而失效时，使用 Refresh Token 就可以获取到新的 Access Token，如果 Refresh Token 也失效了，用户就只能重新登录了。

在 JWT 的实践中，引入 Refresh Token，将会话管理流程改进如下。

- 客户端使用用户名密码进行认证
- 服务端生成有效时间较短的 Access Token（例如 10 分钟），和有效时间较长的 Refresh Token（例如 7 天）
- 客户端访问需要认证的接口时，携带 Access Token
- 如果 Access Token 没有过期，服务端鉴权后返回给客户端需要的数据
- 如果携带 Access Token 访问需要认证的接口时鉴权失败（例如返回 401 错误），则客户端使用 Refresh Token 向刷新接口申请新的 Access Token
- 如果 Refresh Token 没有过期，服务端向客户端下发新的 Access Token
- 客户端使用新的 Access Token 访问需要认证的接口



![img](https://pic1.zhimg.com/80/v2-8f29f24dd291ddf46abda5d5ab7bec6c_720w.webp)



将生成的 Refresh Token 以及过期时间存储在服务端的数据库中，由于 Refresh Token 不会在客户端请求业务接口时验证，只有在申请新的 Access Token 时才会验证，所以将 Refresh Token 存储在数据库中，不会对业务接口的响应时间造成影响，也不需要像 Session 一样一直保持在内存中以应对大量的请求。

上述的架构，提供了服务端禁用用户 Token 的方式，当用户需要登出或禁用用户时，只需要将服务端的 Refresh Token 禁用或删除，用户就会在 Access Token 过期后，由于无法获取到新的 Access Token 而再也无法访问需要认证的接口。这样的方式虽然会有一定的窗口期（取决于 Access Token 的失效时间），但是结合用户登出时客户端删除 Access Token 的操作，基本上可以适应常规情况下对用户认证鉴权的精度要求。

## 总结

JWT 的使用，提高了开发者开发用户认证鉴权功能的效率，降低了系统架构复杂度，避免了大量的数据库和缓存查询，降低了业务接口的响应延迟。然而 JWT 的这些优点也增加了 Token 管理上的难度，通过引入 Refresh Token，既能继续使用 JWT 所带来的优势，又能使得 Token 管理的精度符合业务的需求。
