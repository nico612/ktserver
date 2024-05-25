package captcha

import "github.com/mojocn/base64Captcha"

// Base64Generator 验证码生成器, 这里使用了 base64Captcha 库 的 Driver生成图形验证码
type Base64Generator struct {
	driver base64Captcha.Driver
}

func NewBase64Generator(driver base64Captcha.Driver) *Base64Generator {
	return &Base64Generator{
		driver: driver,
	}
}

func DefaultBase64Generator() *Base64Generator {
	return NewBase64Generator(base64Captcha.DefaultDriverDigit)
}

func (g *Base64Generator) Generate() (id string, b64s string, answer string, err error) {
	id, content, answer := g.driver.GenerateIdQuestionAnswer()
	item, err := g.driver.DrawCaptcha(content)
	if err != nil {
		return "", "", "", err
	}
	return id, item.EncodeB64string(), answer, nil
}
