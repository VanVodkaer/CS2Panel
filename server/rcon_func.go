package server

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	rcon "github.com/forewing/csgo-rcon"
)

// 全局互斥锁，确保RCON命令顺序执行
var rconMutex sync.Mutex

// ExecRconCommand 执行Rcon命令（简化版本）
func ExecRconCommand(name string, command string) (string, error) {
	// 使用互斥锁确保命令顺序执行
	rconMutex.Lock()
	defer rconMutex.Unlock()

	port, err := GetEnvValue(name, "CS2_RCON_PORT")
	if err != nil {
		return "", fmt.Errorf("获取Rcon端口失败，请检查容器是否存在或Rcon端口是否正确配置: %v", err)
	}

	passwd, err := GetEnvValue(name, "CS2_RCONPW")
	if err != nil {
		return "", fmt.Errorf("获取Rcon密码失败: %v", err)
	}

	// 创建Rcon客户端并执行命令
	client := rcon.New(fmt.Sprintf("localhost:%s", port), passwd, 1*time.Second)
	response, err := client.Execute(command)
	if err != nil {
		return "", fmt.Errorf("执行Rcon命令失败: %v", err)
	}

	return response, nil
}

// ExecRconCommands 执行多个Rcon命令
func ExecRconCommands(name string, commands []string) ([]string, error) {
	var responses []string
	for _, cmd := range commands {
		res, err := ExecRconCommand(name, cmd)
		if err != nil {
			return nil, err
		}
		responses = append(responses, res)
	}
	return responses, nil
}

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
				player := parsePlayerLine(line)
				players = append(players, player)
			}
		}
	}

	result.Spawngroups = spawngroups
	result.PlayerList = players
	return result, nil
}

// ====================== 工具函数 ======================

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

	return PlayerInfo{
		ID:      atoi(fields[0]),
		Time:    fields[1],
		Ping:    atoi(fields[2]),
		Loss:    atoi(fields[3]),
		State:   fields[4],
		Rate:    atoi(fields[5]),
		Address: fields[6],
		Name:    strings.Trim(fields[7], "'"),
	}
}
