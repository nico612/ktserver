package main

import (
	"fmt"
	"golang.org/x/text/language"
	"ktserver/internal/pkg/locales"
	"ktserver/pkg/i18n"
)

func main() {

	i18n := i18n.New(i18n.WithFormat("yaml"), i18n.WithFS(locales.Locales))

	fmt.Println(i18n.T(locales.NoPermission))

	i18nC := i18n.Select(language.Chinese)
	fmt.Println(i18nC.T(locales.NoPermission))

}
