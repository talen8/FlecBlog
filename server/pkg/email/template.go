package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

// TemplateManager 邮件模板管理器
type TemplateManager struct {
	mu           sync.RWMutex
	templates    map[string]*template.Template
	templatePath string
}

// NewTemplateManager 创建模板管理器
func NewTemplateManager(templatePath string) *TemplateManager {
	if templatePath == "" {
		templatePath = "templates/email"
	}

	tm := &TemplateManager{
		templates:    make(map[string]*template.Template),
		templatePath: templatePath,
	}

	return tm
}

// LoadTemplate 加载单个模板
// 优先级：自定义模板 > 默认模板
func (tm *TemplateManager) LoadTemplate(name string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// 1. 先尝试加载自定义模板（data/templates）
	customPath := filepath.Join("data/templates/email", name+".tmpl")
	if _, err := os.Stat(customPath); err == nil {
		tmpl, err := template.ParseFiles(customPath)
		if err == nil {
			tm.templates[name] = tmpl
			return nil
		}
	}

	// 2. 加载默认模板（templates）
	defaultPath := filepath.Join(tm.templatePath, name+".tmpl")
	tmpl, err := template.ParseFiles(defaultPath)
	if err != nil {
		return fmt.Errorf("加载模板失败 %s: %w", name, err)
	}

	tm.templates[name] = tmpl
	return nil
}

// LoadTemplates 批量加载模板
func (tm *TemplateManager) LoadTemplates(names []string) {
	for _, name := range names {
		if err := tm.LoadTemplate(name); err != nil {
			slog.Warn("加载邮件模板失败", slog.String("template", name), slog.Any("err", err))
		}
	}
}

// GetTemplate 获取模板
func (tm *TemplateManager) GetTemplate(name string) (*template.Template, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	tmpl, ok := tm.templates[name]
	if !ok {
		return nil, fmt.Errorf("模板不存在: %s", name)
	}

	return tmpl, nil
}

// Render 渲染模板
func (tm *TemplateManager) Render(templateName string, data interface{}) (string, error) {
	tmpl, err := tm.GetTemplate(templateName)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("渲染模板失败 %s: %w", templateName, err)
	}

	return buf.String(), nil
}

// 全局模板管理器实例
var globalTemplateManager *TemplateManager
var once sync.Once

// GetGlobalTemplateManager 获取全局模板管理器
func GetGlobalTemplateManager() *TemplateManager {
	once.Do(func() {
		globalTemplateManager = NewTemplateManager("")
		// 加载所有常用模板
		globalTemplateManager.LoadTemplates([]string{
			"password_reset",
			"comment_reply",
			"comment_new",
			"feedback_new",
			"friend_apply",
			"default",
		})
	})
	return globalTemplateManager
}
