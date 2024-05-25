


## 图形验证码
图形验证码主要功能：生成、验证、速率限制、锁定。

- 使用redis存储验证码，并实现错误次数锁定限制
- 使用redis的score有序集合实现速率限制
- 使用 `github.com/mojocn/base64Captcha` 库生成图形验证码