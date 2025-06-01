package server

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/VanVodkaer/CS2Panel/rcon"
)

// 连接池结构
type RconPool struct {
	mu          sync.RWMutex
	connections map[string]*RconConnection
	maxSize     int
}

// RCON连接封装
type RconConnection struct {
	client    *rcon.Client
	lastUsed  time.Time
	mutex     sync.Mutex // 每个连接的互斥锁
	connected bool
}

// 创建新的连接

func NewRconPool(maxSize int) *RconPool {
	pool := &RconPool{
		connections: make(map[string]*RconConnection),
		maxSize:     maxSize,
	}

	// 启动清理协程，定期清理空闲连接
	go pool.cleanupWorker()

	return pool
}

// 获取或创建连接
func (rp *RconPool) GetConnection(name, port, passwd string) (*RconConnection, error) {
	rp.mu.RLock()
	if conn, exists := rp.connections[name]; exists && conn.connected {
		conn.lastUsed = time.Now()
		rp.mu.RUnlock()
		return conn, nil
	}
	rp.mu.RUnlock()

	// 需要创建新连接
	rp.mu.Lock()
	defer rp.mu.Unlock()

	// 双重检查
	if conn, exists := rp.connections[name]; exists && conn.connected {
		conn.lastUsed = time.Now()
		return conn, nil
	}

	// 检查连接池大小
	if len(rp.connections) >= rp.maxSize {
		return nil, fmt.Errorf("连接池已满")
	}

	// 创建新的RCON客户端，减少超时时间
	client := rcon.New(fmt.Sprintf("localhost:%s", port), passwd, 200*time.Millisecond)

	conn := &RconConnection{
		client:    client,
		lastUsed:  time.Now(),
		connected: true,
	}

	rp.connections[name] = conn
	return conn, nil
}

// 清理空闲连接
func (rp *RconPool) cleanupWorker() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rp.mu.Lock()
		now := time.Now()
		for name, conn := range rp.connections {
			// 清理5分钟未使用的连接
			if now.Sub(conn.lastUsed) > 5*time.Minute {
				// forewing/csgo-rcon 包没有Close方法，直接从池中移除
				delete(rp.connections, name)
			}
		}
		rp.mu.Unlock()
	}
}

// 关闭连接池
func (rp *RconPool) Close() {
	rp.mu.Lock()
	defer rp.mu.Unlock()

	// forewing/csgo-rcon 包没有Close方法，直接清空连接池
	rp.connections = make(map[string]*RconConnection)
}

// 全局连接池
var rconPool = NewRconPool(20)

// 执行单个RCON命令 - 优化版本
func ExecRconCommand(name string, command string) (string, error) {
	// 获取环境变量
	port, err := GetEnvValue(name, "CS2_RCON_PORT")
	if err != nil {
		return "", fmt.Errorf("获取Rcon端口失败: %v", err)
	}

	passwd, err := GetEnvValue(name, "CS2_RCONPW")
	if err != nil {
		return "", fmt.Errorf("获取Rcon密码失败: %v", err)
	}

	// 获取连接
	conn, err := rconPool.GetConnection(name, port, passwd)
	if err != nil {
		return "", fmt.Errorf("获取连接失败: %v", err)
	}

	// 对单个连接加锁，而不是全局锁
	conn.mutex.Lock()
	defer conn.mutex.Unlock()

	// 执行命令
	response, err := conn.client.Execute(command)
	if err != nil {
		// 如果执行失败，标记连接为不可用
		conn.connected = false
		return "", fmt.Errorf("执行Rcon命令失败: %v", err)
	}

	return response, nil
}

// 批量执行RCON命令 - 优化版本
func ExecRconCommands(name string, commands []string) ([]string, error) {
	if len(commands) == 0 {
		return []string{}, nil
	}

	// 对于同一个服务器的多个命令，使用单个连接顺序执行
	// 这样可以避免连接竞争，提高效率
	responses := make([]string, len(commands))

	for i, cmd := range commands {
		response, err := ExecRconCommand(name, cmd)
		if err != nil {
			return responses, fmt.Errorf("执行第%d个命令失败: %v", i+1, err)
		}
		responses[i] = response
	}

	return responses, nil
}

// 并发执行多个服务器的命令
func ExecRconCommandsConcurrent(serverCommands map[string][]string) (map[string][]string, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make(map[string][]string)
	errors := make(map[string]error)

	for serverName, commands := range serverCommands {
		wg.Add(1)
		go func(name string, cmds []string) {
			defer wg.Done()

			responses, err := ExecRconCommands(name, cmds)

			mu.Lock()
			if err != nil {
				errors[name] = err
			} else {
				results[name] = responses
			}
			mu.Unlock()
		}(serverName, commands)
	}

	wg.Wait()

	// 检查是否有错误
	if len(errors) > 0 {
		return results, fmt.Errorf("部分服务器执行失败: %v", errors)
	}

	return results, nil
}

// 健康检查 - 测试连接是否正常
func (rp *RconPool) HealthCheck(name, port, passwd string) error {
	conn, err := rp.GetConnection(name, port, passwd)
	if err != nil {
		return err
	}

	conn.mutex.Lock()
	defer conn.mutex.Unlock()

	// 执行一个简单的命令测试连接
	_, err = conn.client.Execute("echo test")
	if err != nil {
		conn.connected = false
		return err
	}

	return nil
}

// ====================== 服务器状态相关 ======================
type ServerStatus struct {
	ServerAddress string        `json:"server_address"`
	ClientStatus  string        `json:"client_status"`
	CurrentState  string        `json:"current_state"`
	Source        string        `json:"source"`
	Hostname      string        `json:"hostname"`
	SpawnGroup    int           `json:"spawn_group"`
	Version       string        `json:"version"`
	SteamID       string        `json:"steam_id"`
	SteamID64     string        `json:"steam_id_64"`
	LocalIP       string        `json:"local_ip"`
	PublicIP      string        `json:"public_ip"`
	OS            string        `json:"os"`
	PlayerSummary PlayerSummary `json:"player_summary"`
	Spawngroups   []Spawngroup  `json:"spawngroups"`
	PlayerList    []PlayerInfo  `json:"player_list"`
}

type PlayerSummary struct {
	Humans       int  `json:"humans"`
	Bots         int  `json:"bots"`
	MaxPlayers   int  `json:"max_players"`
	Hibernating  bool `json:"hibernating"`
	ReservedSlot bool `json:"reserved_slot"`
}

type Spawngroup struct {
	ID    int      `json:"id"`
	Path  string   `json:"path"`
	Type  string   `json:"type"`
	Flags []string `json:"flags"`
}

type PlayerInfo struct {
	ID      int    `json:"id"`
	Time    string `json:"time"`
	Ping    int    `json:"ping"`
	Loss    int    `json:"loss"`
	State   string `json:"state"`
	Rate    int    `json:"rate"`
	Address string `json:"address"`
	Name    string `json:"name"`
}

// 获取服务器状态（主函数调用）
func GetServerStatus(name string) (ServerStatus, error) {
	statusOutput, err := ExecRconCommand(name, "status")
	if err != nil {
		return ServerStatus{}, fmt.Errorf("获取服务器状态失败: %v", err)
	}
	return ParseCS2Status(statusOutput)
}

// ====================== 工具函数 ======================

// 解析 status 输出文本为结构体
func ParseCS2Status(output string) (ServerStatus, error) {
	var result ServerStatus
	lines := strings.Split(output, "\n")

	var spawngroups []Spawngroup
	var players []PlayerInfo

	inSpawngroup := false
	inPlayers := false

	var playerLineRegex = regexp.MustCompile(`^\d+`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "Server:"):
			result.ServerAddress = extractAfter(line, "Server:  Running [", "]")

		case strings.HasPrefix(line, "Client:"):
			result.ClientStatus = extractAfter(line, "Client:  ", "")

		case strings.HasPrefix(line, "@ Current"):
			result.CurrentState = extractAfter(line, "@ Current  :  ", "")

		case strings.HasPrefix(line, "source"):
			result.Source = extractAfter(line, "source   : ", "")

		case strings.HasPrefix(line, "hostname"):
			result.Hostname = extractAfter(line, "hostname : ", "")

		case strings.HasPrefix(line, "spawn"):
			result.SpawnGroup = atoi(extractAfter(line, "spawn    : ", ""))

		case strings.HasPrefix(line, "version"):
			result.Version = extractAfter(line, "version  : ", "")

		case strings.HasPrefix(line, "steamid"):
			raw := extractAfter(line, "steamid  : ", "")
			parts := strings.Split(raw, " ")
			if len(parts) >= 2 {
				result.SteamID = parts[0]
				result.SteamID64 = strings.Trim(parts[1], "()")
			}

		case strings.HasPrefix(line, "udp/ip"):
			result.PublicIP = extractAfter(line, "public ", "")
			result.LocalIP = extractBetween(line, "udp/ip   : ", " (public")

		case strings.HasPrefix(line, "os/type"):
			result.OS = extractAfter(line, "os/type  : ", "")

		case strings.HasPrefix(line, "players  :"):
			playerSummary := extractAfter(line, "players  : ", "")
			humans, bots, max := parsePlayerCounts(playerSummary)
			result.PlayerSummary = PlayerSummary{
				Humans:       humans,
				Bots:         bots,
				MaxPlayers:   max,
				Hibernating:  !strings.Contains(playerSummary, "not hibernating"),
				ReservedSlot: !strings.Contains(playerSummary, "unreserved"),
			}

		case strings.HasPrefix(line, "---------spawngroups"):
			inSpawngroup = true
			continue

		case strings.HasPrefix(line, "---------players--------"):
			inSpawngroup = false
			inPlayers = true
			continue

		case strings.HasPrefix(line, "#end"):
			break

		case inSpawngroup && strings.HasPrefix(line, "loaded spawngroup"):
			sg := parseSpawngroupLine(line)
			spawngroups = append(spawngroups, sg)

		case inPlayers:
			// 跳过标题行，例如 "id     time ping loss      state   rate adr name"
			if strings.HasPrefix(strings.ToLower(strings.TrimSpace(line)), "id") {
				continue
			}
			// 如果行以数字开头，则认为是有效的玩家数据行
			if playerLineRegex.MatchString(strings.TrimSpace(line)) {
				// 检查行是否有足够的字段再进行解析
				fields := strings.Fields(line)
				if len(fields) >= 8 {
					player := parsePlayerLine(line)
					// 只有当解析成功时才添加到列表中
					if player.ID > 0 || player.Name != "" {
						players = append(players, player)
					}
				}
			}
		}
	}

	result.Spawngroups = spawngroups
	result.PlayerList = players
	return result, nil
}

func extractAfter(s, prefix, suffix string) string {
	s = strings.TrimPrefix(s, prefix)
	if suffix != "" && strings.Contains(s, suffix) {
		s = strings.SplitN(s, suffix, 2)[0]
	}
	return strings.TrimSpace(s)
}

func extractBetween(s, start, end string) string {
	s = strings.TrimPrefix(s, start)
	if idx := strings.Index(s, end); idx != -1 {
		return strings.TrimSpace(s[:idx])
	}
	return ""
}

func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func parsePlayerCounts(s string) (humans, bots, max int) {
	re := regexp.MustCompile(`(\d+) humans, (\d+) bots \((\d+) max`)
	match := re.FindStringSubmatch(s)
	if len(match) == 4 {
		humans = atoi(match[1])
		bots = atoi(match[2])
		max = atoi(match[3])
	}
	return
}

func parseSpawngroupLine(line string) Spawngroup {
	re := regexp.MustCompile(`spawngroup\(\s*(\d+)\)\s*:\s*SV:\s*\[(\d+):\s*([^|]+)\s*\|\s*(.*)\]`)
	match := re.FindStringSubmatch(line)
	if len(match) < 5 {
		return Spawngroup{}
	}
	flags := strings.Split(match[4], "|")
	for i := range flags {
		flags[i] = strings.TrimSpace(flags[i])
	}
	return Spawngroup{
		ID:    atoi(match[1]),
		Path:  strings.TrimSpace(match[3]),
		Type:  "main lump",
		Flags: flags,
	}
}

func parsePlayerLine(line string) PlayerInfo {
	fields := strings.Fields(line)

	// 检查字段数量，确保至少有8个字段
	if len(fields) < 8 {
		// 如果字段不足，返回一个空的 PlayerInfo 或者带有部分信息的结构
		return PlayerInfo{
			ID:      0,
			Time:    "",
			Ping:    0,
			Loss:    0,
			State:   "",
			Rate:    0,
			Address: "",
			Name:    "",
		}
	}

	// 安全地解析各个字段
	playerInfo := PlayerInfo{
		ID:      atoi(fields[0]),
		Time:    fields[1],
		Ping:    atoi(fields[2]),
		Loss:    atoi(fields[3]),
		State:   fields[4],
		Rate:    atoi(fields[5]),
		Address: fields[6],
		Name:    strings.Trim(fields[7], "'"),
	}

	return playerInfo
}
