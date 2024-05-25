package i18n

import (
	"embed"
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type I18n struct {
	ops       *Options
	bundle    *i18n.Bundle    // Bundle 是一个包含所有翻译的集合
	localizer *i18n.Localizer // Localizer 用于根据语言标签查找翻译
	lang      language.Tag    // 语言 如 language.English、language.Chinese、 language.TraditionalChinese（繁体中文）
}

func New(options ...func(*Options)) (rp *I18n) {
	ops := getOptionsOrSetDefault(nil)
	for _, f := range options {
		f(ops)
	}
	// 创建一个新的 Bundle 实例
	bundle := i18n.NewBundle(ops.language)
	// 创建一个新的 Localizer 实例
	localizer := i18n.NewLocalizer(bundle, ops.language.String())

	// 根据文件格式注册 解码器， 默认为 yaml
	switch ops.format {
	case "toml":
		bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	case "json":
		bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	default:
		bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	}

	rp = &I18n{
		ops:       ops,
		bundle:    bundle,
		localizer: localizer,
		lang:      ops.language,
	}

	for _, item := range ops.files {
		rp.Add(item)
	}
	rp.AddFS(ops.fs)
	return
}

// Add 添加翻译文件
func (i *I18n) Add(f string) {
	info, err := os.Stat(f)
	if err != nil {
		return
	}

	if info.IsDir() {
		// 遍历目录下的所有文件
		filepath.Walk(f, func(path string, fi os.FileInfo, err error) error {
			if !fi.IsDir() {
				// 从文件中加载翻译信息
				i.bundle.LoadMessageFile(path)
			}
			return nil
		})
	} else {
		i.bundle.LoadMessageFile(f)
	}
}

// AddFS 添加 embed.FS 文件系统
func (i *I18n) AddFS(fs embed.FS) {
	files := readFS(fs, ".")
	for _, name := range files {
		i.bundle.LoadMessageFileFS(fs, name)
	}
}

// readFS 读取文件系统
func readFS(fs embed.FS, dir string) (rp []string) {
	rp = make([]string, 0)
	dirs, err := fs.ReadDir(dir)
	if err != nil {
		return
	}
	for _, item := range dirs {
		name := dir + string(os.PathSeparator) + item.Name()
		if dir == "." {
			name = item.Name()
		}
		if item.IsDir() {
			rp = append(rp, readFS(fs, name)...)
		} else {
			rp = append(rp, name)
		}
	}
	return
}

// Select 切换语言
func (i I18n) Select(lang language.Tag) *I18n {
	if lang.String() == "und" {
		lang = i.ops.language
	}
	return &I18n{
		ops:       i.ops,
		bundle:    i.bundle,
		localizer: i18n.NewLocalizer(i.bundle, lang.String()),
		lang:      lang,
	}
}

// Language 返回当前语言
func (i *I18n) Language() language.Tag {
	return i.lang
}

// LocalizeT 本地化翻译
func (i *I18n) LocalizeT(message *i18n.Message) string {
	if message == nil {
		return ""
	}

	// 本地化翻译
	rp, err := i.localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: message,
	})

	if err != nil {
		// 无法翻译时使用 ID 作为默认消息
		rp = message.ID
	}

	return rp

}

// LocalizeE 本地化翻译 warp error
func (i *I18n) LocalizeE(message *i18n.Message) error {
	return errors.New(i.LocalizeT(message))
}

// T 根据 ID 获取翻译
func (i *I18n) T(id string) string {
	return i.LocalizeT(&i18n.Message{
		ID: id,
	})
}

// E warp error 根据 ID 获取翻译
func (i *I18n) E(id string) error {
	return errors.New(i.T(id))
}
