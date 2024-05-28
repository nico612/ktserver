package i18n

import (
	"embed"
	"golang.org/x/text/language"
)

type Options struct {
	format   string       // 文件格式
	language language.Tag // 语言
	files    []string     // 语言文件
	fs       embed.FS     // 文件系统，用于读取文件
}

// WithFormat 设置文件格式
func WithFormat(format string) func(*Options) {
	return func(options *Options) {
		if format != "" {
			getOptionsOrSetDefault(options).format = format
		}
	}
}

// WithLanguage 设置语言
func WithLanguage(lang language.Tag) func(*Options) {
	return func(options *Options) {
		if lang.String() != "und" { // 如果语言不是未知的
			getOptionsOrSetDefault(options).language = lang
		}
	}
}

// WithFile 设置语言文件
func WithFile(f string) func(*Options) {
	return func(options *Options) {
		if f != "" {
			getOptionsOrSetDefault(options).files = append(getOptionsOrSetDefault(options).files, f)
		}
	}

}

// WithFS 设置文件系统
func WithFS(fs embed.FS) func(*Options) {
	return func(options *Options) {
		getOptionsOrSetDefault(options).fs = fs
	}
}

func getOptionsOrSetDefault(options *Options) *Options {
	if options == nil {
		return &Options{
			format:   "yml",            // 默认文件格式为 yml
			language: language.English, // 默认语言为英语
			files:    []string{},
		}
	}
	return options
}
