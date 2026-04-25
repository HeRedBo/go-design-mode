// Package facade 演示外观模式 (Facade Pattern)
//
// 外观模式解决的问题：
// 1. 为复杂子系统提供简化接口
// 2. 降低客户端与子系统的耦合
// 3. 隐藏子系统的复杂性
//
// 使用场景：
// - 简化复杂库的使用
// - 统一多个服务的接口
// - 提供默认配置和流程
// - 封装复杂的初始化过程
//
// 与适配器的区别：
// - 适配器：改变接口以匹配客户端期望
// - 外观：简化接口，提供便捷访问
//
// 核心组件：
// 1. Facade 外观：简化接口
// 2. Subsystem 子系统：复杂的功能模块
package facade

import "fmt"

// ============ 示例 1: 家庭影院系统 ============

// TV 电视
type TV struct{}

// On 打开电视
func (t *TV) On() {
	fmt.Println("[TV] TV is on")
}

// Off 关闭电视
func (t *TV) Off() {
	fmt.Println("[TV] TV is off")
}

// SetChannel 设置频道
func (t *TV) SetChannel(channel int) {
	fmt.Printf("[TV] Channel set to %d\n", channel)
}

// DVDPlayer DVD 播放器
type DVDPlayer struct{}

// On 打开播放器
func (d *DVDPlayer) On() {
	fmt.Println("[DVD] DVD player is on")
}

// Off 关闭播放器
func (d *DVDPlayer) Off() {
	fmt.Println("[DVD] DVD player is off")
}

// Play 播放
func (d *DVDPlayer) Play(movie string) {
	fmt.Printf("[DVD] Playing: %s\n", movie)
}

// Stop 停止
func (d *DVDPlayer) Stop() {
	fmt.Println("[DVD] Stopped")
}

// SoundSystem 音响系统
type SoundSystem struct{}

// On 打开音响
func (s *SoundSystem) On() {
	fmt.Println("[Sound] Sound system is on")
}

// Off 关闭音响
func (s *SoundSystem) Off() {
	fmt.Println("[Sound] Sound system is off")
}

// SetVolume 设置音量
func (s *SoundSystem) SetVolume(level int) {
	fmt.Printf("[Sound] Volume set to %d\n", level)
}

// Lights 灯光系统
type Lights struct{}

// Dim 调暗
func (l *Lights) Dim(level int) {
	fmt.Printf("[Lights] Dimmed to %d%%\n", level)
}

// On 打开
func (l *Lights) On() {
	fmt.Println("[Lights] Lights are on")
}

// HomeTheaterFacade 家庭影院外观
type HomeTheaterFacade struct {
	tv    *TV
	dvd   *DVDPlayer
	sound *SoundSystem
	lights *Lights
}

// NewHomeTheaterFacade 创建家庭影院外观
func NewHomeTheaterFacade() *HomeTheaterFacade {
	return &HomeTheaterFacade{
		tv:     &TV{},
		dvd:    &DVDPlayer{},
		sound:  &SoundSystem{},
		lights: &Lights{},
	}
}

// WatchMovie 观看电影（一键启动）
func (f *HomeTheaterFacade) WatchMovie(movie string) {
	fmt.Println("=== Starting Movie Night ===")
	f.lights.Dim(20)
	f.sound.On()
	f.sound.SetVolume(70)
	f.tv.On()
	f.tv.SetChannel(1)
	f.dvd.On()
	f.dvd.Play(movie)
	fmt.Println("=== Enjoy your movie! ===")
}

// EndMovie 结束电影（一键关闭）
func (f *HomeTheaterFacade) EndMovie() {
	fmt.Println("=== Ending Movie Night ===")
	f.dvd.Stop()
	f.dvd.Off()
	f.tv.Off()
	f.sound.Off()
	f.lights.On()
	fmt.Println("=== Good night! ===")
}

// ============ 示例 2: 订单处理系统 ============

// Inventory 库存系统
type Inventory struct{}

// CheckStock 检查库存
func (i *Inventory) CheckStock(productID string) bool {
	fmt.Printf("[Inventory] Checking stock for %s\n", productID)
	return true
}

// ReserveStock 预留库存
func (i *Inventory) ReserveStock(productID string, quantity int) {
	fmt.Printf("[Inventory] Reserved %d units of %s\n", quantity, productID)
}

// PaymentGateway 支付网关
type PaymentGateway struct{}

// ProcessPayment 处理支付
func (p *PaymentGateway) ProcessPayment(amount float64, card string) bool {
	fmt.Printf("[Payment] Processing payment: $%.2f with card %s\n", amount, card)
	return true
}

// Shipping 物流系统
type Shipping struct{}

// CalculateCost 计算运费
func (s *Shipping) CalculateCost(address string) float64 {
	fmt.Printf("[Shipping] Calculating cost to: %s\n", address)
	return 10.0
}

// Ship 发货
func (s *Shipping) Ship(productID string, address string) {
	fmt.Printf("[Shipping] Shipping %s to %s\n", productID, address)
}

// Notification 通知系统
type Notification struct{}

// SendConfirmation 发送确认通知
func (n *Notification) SendConfirmation(email string, orderID string) {
	fmt.Printf("[Notification] Sending confirmation to %s for order %s\n", email, orderID)
}

// OrderFacade 订单处理外观
type OrderFacade struct {
	inventory *Inventory
	payment   *PaymentGateway
	shipping  *Shipping
	notification *Notification
}

// NewOrderFacade 创建订单处理外观
func NewOrderFacade() *OrderFacade {
	return &OrderFacade{
		inventory:    &Inventory{},
		payment:      &PaymentGateway{},
		shipping:     &Shipping{},
		notification: &Notification{},
	}
}

// PlaceOrder 下单（简化接口）
func (f *OrderFacade) PlaceOrder(productID string, quantity int, amount float64, card string, address string, email string) (string, error) {
	fmt.Println("=== Placing Order ===")

	// 检查库存
	if !f.inventory.CheckStock(productID) {
		return "", fmt.Errorf("product out of stock")
	}

	// 处理支付
	if !f.payment.ProcessPayment(amount, card) {
		return "", fmt.Errorf("payment failed")
	}

	// 预留库存
	f.inventory.ReserveStock(productID, quantity)

	// 计算运费并发货
	cost := f.shipping.CalculateCost(address)
	fmt.Printf("[Order] Shipping cost: $%.2f\n", cost)
	f.shipping.Ship(productID, address)

	// 发送确认
	orderID := fmt.Sprintf("ORD-%s-%d", productID, quantity)
	f.notification.SendConfirmation(email, orderID)

	fmt.Println("=== Order Completed ===")
	return orderID, nil
}

// ============ 示例 3: 数据库操作外观 ============

// Connection 数据库连接
type Connection struct{}

// Open 打开连接
func (c *Connection) Open(connectionString string) error {
	fmt.Printf("[DB] Opening connection: %s\n", connectionString)
	return nil
}

// Close 关闭连接
func (c *Connection) Close() {
	fmt.Println("[DB] Connection closed")
}

// Query 执行查询
type Query struct{}

// Execute 执行查询
func (q *Query) Execute(sql string) ([]map[string]interface{}, error) {
	fmt.Printf("[Query] Executing: %s\n", sql)
	return []map[string]interface{}{{"id": 1, "name": "test"}}, nil
}

// Transaction 事务
type Transaction struct{}

// Begin 开始事务
func (t *Transaction) Begin() {
	fmt.Println("[Transaction] Transaction started")
}

// Commit 提交事务
func (t *Transaction) Commit() {
	fmt.Println("[Transaction] Transaction committed")
}

// Rollback 回滚事务
func (t *Transaction) Rollback() {
	fmt.Println("[Transaction] Transaction rolled back")
}

// DatabaseFacade 数据库操作外观
type DatabaseFacade struct {
	connection  *Connection
	query       *Query
	transaction *Transaction
}

// NewDatabaseFacade 创建数据库外观
func NewDatabaseFacade() *DatabaseFacade {
	return &DatabaseFacade{
		connection:  &Connection{},
		query:       &Query{},
		transaction: &Transaction{},
	}
}

// ExecuteQuery 执行查询（简化接口）
func (f *DatabaseFacade) ExecuteQuery(connectionString string, sql string) ([]map[string]interface{}, error) {
	if err := f.connection.Open(connectionString); err != nil {
		return nil, err
	}
	defer f.connection.Close()

	return f.query.Execute(sql)
}

// ExecuteTransaction 执行事务（简化接口）
func (f *DatabaseFacade) ExecuteTransaction(connectionString string, queries []string) error {
	if err := f.connection.Open(connectionString); err != nil {
		return err
	}
	defer f.connection.Close()

	f.transaction.Begin()

	for _, sql := range queries {
		if _, err := f.query.Execute(sql); err != nil {
			f.transaction.Rollback()
			return err
		}
	}

	f.transaction.Commit()
	return nil
}
