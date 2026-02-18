package i18n

import (
	"errors"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gitlab-digital.tele2.kz/digital/core/backend/golang/logger"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"plugin"
	"strings"
)

var (
	ErrDirNotFound         = errors.New("directory not found")
	ErrDestinationNotFound = errors.New("destination not found")
	ErrFailedLoadFile      = errors.New("failed to load file")
	ErrFileInfoNull        = errors.New("file info is nil")
)

// Localizer interface of plugin
type Localizer interface {
	GetMessage(code, lng string) string
}

type Options struct {
	Name            string
	Configurator    plugin.Plugin
	DefaultLanguage language.Tag
	TemplateDir     string
	TemplateExt     string
}

type Plugin struct {
	Options *Options
	Bundle  *i18n.Bundle
	Logger  logger.Logger
}

// New  initialize plugin with default  options and load the content file
// for more developing  this plugin you can read here  https://github.com/nicksnyder/go-i18n
func New(contentFile, directory string, lg logger.Logger) (*Plugin, error) {
	//nolint:exhaustruct
	p := &Plugin{
		Options: &Options{
			Name:            "plugin.i18n.goi18n.default",
			DefaultLanguage: language.Russian,
			TemplateDir:     directory,
			TemplateExt:     "json",
		},
		Logger: lg,
	}
	p.Bundle = i18n.NewBundle(p.Options.DefaultLanguage)

	if err := p.loadMessages(contentFile); err != nil {
		return nil, err
	}

	return p, nil
}

// loadMessages function which need upload file with content and store into bundle of plugin
func (p *Plugin) loadMessages(contentFile string) error {
	if _, err := os.Stat(p.Options.TemplateDir); err != nil {
		return errors.Join(ErrDirNotFound, err)
	}

	cleanRoot := filepath.Clean(filepath.Join(p.Options.TemplateDir, contentFile))

	return filepath.Walk(cleanRoot, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return ErrFileInfoNull
		}
		if !info.IsDir() && strings.HasSuffix(path, p.Options.TemplateExt) {
			if err != nil {
				return errors.Join(ErrDestinationNotFound, err)
			}

			_, err = p.Bundle.LoadMessageFile(path)
			if err != nil {
				return errors.Join(ErrFailedLoadFile, err)
			}
		}

		return nil
	})
}

// GetMessage  function that getting message from the bundle by code and language
func (p *Plugin) GetMessage(code, lng string) string {
	localizer := i18n.NewLocalizer(p.Bundle, lng)

	//nolint:exhaustruct
	text, err := localizer.LocalizeMessage(&i18n.Message{
		ID: code,
	})

	if err != nil {
		p.Logger.Errorf("failed to localize message: %w", err)
	}

	return text
}
