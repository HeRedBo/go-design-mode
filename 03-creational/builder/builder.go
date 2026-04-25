// Package builder 演示建造者模式 (Builder Pattern)
//
// 建造者模式解决的问题：
// 1. 复杂对象的分步骤构建
// 2. 构建过程与表示分离
// 3. 支持不同的构建流程
// 4. 避免构造函数参数过多（伸缩构造函数问题）
//
// 使用场景：
// - 构建复杂对象（HTTP 请求、SQL 查询、文档）
// - 对象有很多可选配置项
// - 构建过程需要多个步骤
// - 需要创建不同表示的对象
//
// Go 中的实现特点：
// - 通常使用链式调用（Method Chaining）
// - 与函数选项模式结合使用
// - Builder 返回最终产品
//
// 与函数选项模式的区别：
// - 函数选项：适合简单配置，一次性设置
// - 建造者：适合复杂构建，分步骤，可验证
//
// 核心组件：
// 1. Product：复杂产品对象
// 2. Builder 接口：定义构建步骤
// 3. ConcreteBuilder：具体建造者
// 4. Director：指挥者（可选），定义构建流程
package builder

import (
	"fmt"
	"strings"
)

// ============ 产品：HTTP 请求 ============

// HTTPRequest HTTP 请求对象
type HTTPRequest struct {
	method      string
	url         string
	headers     map[string]string
	body        string
	timeout     int
	retries     int
	authToken   string
	contentType string
}

// String 返回请求的字符串表示
func (r *HTTPRequest) String() string {
	var parts []string
	parts = append(parts, fmt.Sprintf("%s %s", r.method, r.url))

	if r.contentType != "" {
		parts = append(parts, fmt.Sprintf("Content-Type: %s", r.contentType))
	}
	if r.authToken != "" {
		parts = append(parts, "Authorization: Bearer ***")
	}
	for k, v := range r.headers {
		parts = append(parts, fmt.Sprintf("%s: %s", k, v))
	}
	if r.timeout > 0 {
		parts = append(parts, fmt.Sprintf("Timeout: %ds", r.timeout))
	}
	if r.retries > 0 {
		parts = append(parts, fmt.Sprintf("Retries: %d", r.retries))
	}
	if r.body != "" {
		parts = append(parts, fmt.Sprintf("Body: %s", r.body))
	}

	return strings.Join(parts, "\n")
}

// ============ 建造者接口 ============

// RequestBuilder 请求建造者接口
type RequestBuilder interface {
	WithMethod(method string) RequestBuilder
	WithURL(url string) RequestBuilder
	WithHeader(key, value string) RequestBuilder
	WithBody(body string) RequestBuilder
	WithTimeout(seconds int) RequestBuilder
	WithRetries(count int) RequestBuilder
	WithAuth(token string) RequestBuilder
	WithContentType(contentType string) RequestBuilder
	Build() (*HTTPRequest, error)
}

// ============ 具体建造者 ============

// httpRequestBuilder HTTP 请求建造者实现
type httpRequestBuilder struct {
	request *HTTPRequest
}

// NewRequestBuilder 创建新的请求建造者
func NewRequestBuilder() RequestBuilder {
	return &httpRequestBuilder{
		request: &HTTPRequest{
			headers: make(map[string]string),
			timeout: 30,    // 默认超时 30 秒
			retries: 0,     // 默认不重试
			method:  "GET", // 默认 GET 方法
		},
	}
}

// WithMethod 设置请求方法
func (b *httpRequestBuilder) WithMethod(method string) RequestBuilder {
	b.request.method = strings.ToUpper(method)
	return b
}

// WithURL 设置 URL
func (b *httpRequestBuilder) WithURL(url string) RequestBuilder {
	b.request.url = url
	return b
}

// WithHeader 添加请求头
func (b *httpRequestBuilder) WithHeader(key, value string) RequestBuilder {
	b.request.headers[key] = value
	return b
}

// WithBody 设置请求体
func (b *httpRequestBuilder) WithBody(body string) RequestBuilder {
	b.request.body = body
	return b
}

// WithTimeout 设置超时时间
func (b *httpRequestBuilder) WithTimeout(seconds int) RequestBuilder {
	b.request.timeout = seconds
	return b
}

// WithRetries 设置重试次数
func (b *httpRequestBuilder) WithRetries(count int) RequestBuilder {
	b.request.retries = count
	return b
}

// WithAuth 设置认证令牌
func (b *httpRequestBuilder) WithAuth(token string) RequestBuilder {
	b.request.authToken = token
	return b
}

// WithContentType 设置内容类型
func (b *httpRequestBuilder) WithContentType(contentType string) RequestBuilder {
	b.request.contentType = contentType
	b.request.headers["Content-Type"] = contentType
	return b
}

// Build 构建最终产品
func (b *httpRequestBuilder) Build() (*HTTPRequest, error) {
	// 验证必填字段
	if b.request.url == "" {
		return nil, fmt.Errorf("URL is required")
	}

	if b.request.method == "" {
		b.request.method = "GET"
	}

	// 返回产品副本
	req := *b.request
	req.headers = make(map[string]string)
	for k, v := range b.request.headers {
		req.headers[k] = v
	}

	return &req, nil
}

// ============ 产品：SQL 查询构建器 ============

// SQLQuery SQL 查询对象
type SQLQuery struct {
	table      string
	columns    []string
	conditions []string
	orderBy    string
	limit      int
	offset     int
	join       string
}

// String 返回 SQL 查询语句
func (q *SQLQuery) String() string {
	var sql strings.Builder

	// SELECT 子句
	sql.WriteString("SELECT ")
	if len(q.columns) == 0 {
		sql.WriteString("*")
	} else {
		sql.WriteString(strings.Join(q.columns, ", "))
	}

	// FROM 子句
	sql.WriteString(fmt.Sprintf(" FROM %s", q.table))

	// JOIN 子句
	if q.join != "" {
		sql.WriteString(" " + q.join)
	}

	// WHERE 子句
	if len(q.conditions) > 0 {
		sql.WriteString(" WHERE " + strings.Join(q.conditions, " AND "))
	}

	// ORDER BY 子句
	if q.orderBy != "" {
		sql.WriteString(" ORDER BY " + q.orderBy)
	}

	// LIMIT 子句
	if q.limit > 0 {
		sql.WriteString(fmt.Sprintf(" LIMIT %d", q.limit))
	}

	// OFFSET 子句
	if q.offset > 0 {
		sql.WriteString(fmt.Sprintf(" OFFSET %d", q.offset))
	}

	return sql.String()
}

// ============ SQL 查询建造者 ============

// SQLQueryBuilder SQL 查询建造者接口
type SQLQueryBuilder interface {
	Select(columns ...string) SQLQueryBuilder
	From(table string) SQLQueryBuilder
	Where(condition string) SQLQueryBuilder
	Join(join string) SQLQueryBuilder
	OrderBy(orderBy string) SQLQueryBuilder
	Limit(limit int) SQLQueryBuilder
	Offset(offset int) SQLQueryBuilder
	Build() *SQLQuery
}

// sqlQueryBuilder SQL 查询建造者实现
type sqlQueryBuilder struct {
	query *SQLQuery
}

// NewSQLQueryBuilder 创建 SQL 查询建造者
func NewSQLQueryBuilder() SQLQueryBuilder {
	return &sqlQueryBuilder{
		query: &SQLQuery{
			columns:    make([]string, 0),
			conditions: make([]string, 0),
		},
	}
}

// Select 设置查询列
func (b *sqlQueryBuilder) Select(columns ...string) SQLQueryBuilder {
	b.query.columns = append(b.query.columns, columns...)
	return b
}

// From 设置表名
func (b *sqlQueryBuilder) From(table string) SQLQueryBuilder {
	b.query.table = table
	return b
}

// Where 添加 WHERE 条件
func (b *sqlQueryBuilder) Where(condition string) SQLQueryBuilder {
	b.query.conditions = append(b.query.conditions, condition)
	return b
}

// Join 添加 JOIN 子句
func (b *sqlQueryBuilder) Join(join string) SQLQueryBuilder {
	b.query.join = join
	return b
}

// OrderBy 设置排序
func (b *sqlQueryBuilder) OrderBy(orderBy string) SQLQueryBuilder {
	b.query.orderBy = orderBy
	return b
}

// Limit 设置限制
func (b *sqlQueryBuilder) Limit(limit int) SQLQueryBuilder {
	b.query.limit = limit
	return b
}

// Offset 设置偏移
func (b *sqlQueryBuilder) Offset(offset int) SQLQueryBuilder {
	b.query.offset = offset
	return b
}

// Build 构建 SQL 查询
func (b *sqlQueryBuilder) Build() *SQLQuery {
	q := *b.query
	q.columns = make([]string, len(b.query.columns))
	copy(q.columns, b.query.columns)
	q.conditions = make([]string, len(b.query.conditions))
	copy(q.conditions, b.query.conditions)
	return &q
}

// ============ 产品：HTML 文档构建器 ============

// HTMLDocument HTML 文档对象
type HTMLDocument struct {
	title    string
	head     []string
	body     []string
	styles   []string
	scripts  []string
}

// String 返回 HTML 文档
func (d *HTMLDocument) String() string {
	var html strings.Builder

	html.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
	html.WriteString(fmt.Sprintf("  <title>%s</title>\n", d.title))

	for _, style := range d.styles {
		html.WriteString(fmt.Sprintf("  <style>%s</style>\n", style))
	}
	for _, script := range d.scripts {
		html.WriteString(fmt.Sprintf("  <script>%s</script>\n", script))
	}
	for _, head := range d.head {
		html.WriteString(fmt.Sprintf("  %s\n", head))
	}

	html.WriteString("</head>\n<body>\n")
	for _, body := range d.body {
		html.WriteString(fmt.Sprintf("  %s\n", body))
	}
	html.WriteString("</body>\n</html>")

	return html.String()
}

// ============ HTML 文档建造者 ============

// HTMLBuilder HTML 文档建造者接口
type HTMLBuilder interface {
	WithTitle(title string) HTMLBuilder
	AddStyle(style string) HTMLBuilder
	AddScript(script string) HTMLBuilder
	AddHeadElement(element string) HTMLBuilder
	AddBodyElement(element string) HTMLBuilder
	Build() *HTMLDocument
}

// htmlBuilder HTML 文档建造者实现
type htmlBuilder struct {
	doc *HTMLDocument
}

// NewHTMLBuilder 创建 HTML 文档建造者
func NewHTMLBuilder() HTMLBuilder {
	return &htmlBuilder{
		doc: &HTMLDocument{
			head:    make([]string, 0),
			body:    make([]string, 0),
			styles:  make([]string, 0),
			scripts: make([]string, 0),
		},
	}
}

// WithTitle 设置标题
func (b *htmlBuilder) WithTitle(title string) HTMLBuilder {
	b.doc.title = title
	return b
}

// AddStyle 添加样式
func (b *htmlBuilder) AddStyle(style string) HTMLBuilder {
	b.doc.styles = append(b.doc.styles, style)
	return b
}

// AddScript 添加脚本
func (b *htmlBuilder) AddScript(script string) HTMLBuilder {
	b.doc.scripts = append(b.doc.scripts, script)
	return b
}

// AddHeadElement 添加 head 元素
func (b *htmlBuilder) AddHeadElement(element string) HTMLBuilder {
	b.doc.head = append(b.doc.head, element)
	return b
}

// AddBodyElement 添加 body 元素
func (b *htmlBuilder) AddBodyElement(element string) HTMLBuilder {
	b.doc.body = append(b.doc.body, element)
	return b
}

// Build 构建 HTML 文档
func (b *htmlBuilder) Build() *HTMLDocument {
	doc := *b.doc
	doc.head = make([]string, len(b.doc.head))
	copy(doc.head, b.doc.head)
	doc.body = make([]string, len(b.doc.body))
	copy(doc.body, b.doc.body)
	doc.styles = make([]string, len(b.doc.styles))
	copy(doc.styles, b.doc.styles)
	doc.scripts = make([]string, len(b.doc.scripts))
	copy(doc.scripts, b.doc.scripts)
	return &doc
}

// ============ Director（指挥者） ============

// RequestDirector 请求指挥者，预定义构建流程
type RequestDirector struct {
	builder RequestBuilder
}

// NewRequestDirector 创建请求指挥者
func NewRequestDirector(builder RequestBuilder) *RequestDirector {
	return &RequestDirector{builder: builder}
}

// BuildAPIRequest 构建 API 请求（标准流程）
func (d *RequestDirector) BuildAPIRequest(url, token string) (*HTTPRequest, error) {
	return d.builder.
		WithMethod("POST").
		WithURL(url).
		WithAuth(token).
		WithContentType("application/json").
		WithTimeout(60).
		WithRetries(3).
		Build()
}

// BuildSimpleGet 构建简单 GET 请求
func (d *RequestDirector) BuildSimpleGet(url string) (*HTTPRequest, error) {
	return d.builder.
		WithMethod("GET").
		WithURL(url).
		WithTimeout(30).
		Build()
}
