package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/VanVodkaer/CS2Panel/config"
)

// MapInfo 只保留 name、internal_name 和 playable_modes
type MapInfo struct {
	Name          string   `json:"name"`           // 地图显示名
	InternalName  string   `json:"internal_name"`  // 对应 VMAP/VPK 文件名
	PlayableModes []string `json:"playable_modes"` // 在 CS2 各模式下可玩的模式列表
}

// fetchMaps 抓取地图列表
func fetchMaps(mapclass string) ([]MapInfo, error) {
	var result []MapInfo

	// 1. HTTP GET 页面
	resp, err := http.Get("https://developer.valvesoftware.com/wiki/Counter-Strike_2/Maps")
	if err != nil {
		return nil, fmt.Errorf("请求页面失败：%w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP 返回 %d", resp.StatusCode)
	}

	// 2. 用 goquery 解析 HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("解析 HTML 失败：%w", err)
	}

	// 3. 定位 mapclass 对应的表格
	heading := doc.Find(fmt.Sprintf("span#%s", mapclass)).Closest("h2")
	if heading.Length() == 0 {
		return nil, fmt.Errorf("未找到 %s 标题", mapclass)
	}
	table := heading.NextAll().Filter("table").First()
	if table.Length() == 0 {
		return nil, fmt.Errorf("未找到 %s 表格", mapclass)
	}

	// 4. 拆出两行表头：第一行基础列、第二行各模式列
	headerRows := table.Find("tr").FilterFunction(func(i int, tr *goquery.Selection) bool {
		return tr.Find("th").Length() > 0
	})
	if headerRows.Length() < 2 {
		return nil, fmt.Errorf("表头行不足：找到 %d 行", headerRows.Length())
	}
	modeHeader := headerRows.Eq(1)

	// 5. 读取所有模式名
	modeNames := make([]string, 0)
	modeHeader.Find("th").Each(func(_ int, th *goquery.Selection) {
		if alt, ok := th.Find("img").Attr("alt"); ok && alt != "" {
			modeNames = append(modeNames, strings.TrimSpace(alt))
		} else if txt := strings.TrimSpace(th.Text()); txt != "" {
			modeNames = append(modeNames, txt)
		}
	})
	modeCount := len(modeNames)
	if modeCount == 0 {
		return nil, fmt.Errorf("未检测到任何模式列")
	}

	// 6. 遍历每一行数据
	table.Find("tr").Each(func(_ int, tr *goquery.Selection) {
		cells := tr.Find("td")
		if cells.Length() == 0 {
			return // 跳过表头或空行
		}
		// 收集所有单元格文本
		texts := make([]string, 0, cells.Length())
		cells.Each(func(_ int, td *goquery.Selection) {
			texts = append(texts, strings.TrimSpace(td.Text()))
		})
		// 跳过 Icon 列
		if len(texts) < 2+modeCount {
			return // 列数不够时跳过
		}
		data := texts[1:] // data[0]=MapName, data[1]=Internal, data[2..]=各模式数据

		// 7. 提取可玩的模式
		playable := make([]string, 0, modeCount)
		for i, mode := range modeNames {
			if i+2 < len(data) && strings.EqualFold(data[2+i], "Yes") {
				playable = append(playable, mode)
			}
		}

		// 8. 构造并追加 MapInfo
		mi := MapInfo{
			Name:          data[0],
			InternalName:  data[1],
			PlayableModes: playable,
		}
		result = append(result, mi)
	})

	return result, nil
}

// saveMaps 抓取地图列表并保存到本地 JSON 文件
func saveMaps(mapclass, outputFile string) error {
	maps, err := fetchMaps(mapclass)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(maps, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化为 JSON 失败：%w", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前工作目录失败：%w", err)
	}
	outputPath := filepath.Join(cwd, config.GlobalConfig.Server.PanelDataDir, "maps", outputFile)

	// 关键：先创建目录
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录 %s 失败：%w", dir, err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("写入 %s 失败：%w", outputFile, err)
	}

	return nil
}

// fetchCurrentMaps 抓取并解析“Current Maps”表格，保存为 current_maps.json
func fetchCurrentMaps() error {
	return saveMaps("Current_Maps", "current_maps.json")
}

// fetchFormerMaps 抓取并解析“Former Maps”表格，保存为 former_maps.json
func fetchFormerMaps() error {
	return saveMaps("Former_Maps", "former_maps.json")
}

// getMapList 获取地图列表，如果文件不存在则返回空切片
func getMapList(mapclass string) ([]MapInfo, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("获取当前工作目录失败：%w", err)
	}

	// 用 filepath.Join 生成平台无关的文件路径
	fileName := fmt.Sprintf("%s_maps.json", mapclass)
	filePath := filepath.Join(cwd, config.GlobalConfig.Server.PanelDataDir, "maps", fileName)

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// 如果文件不存在，按需返回空列表
			return []MapInfo{}, nil
		}
		return nil, fmt.Errorf("读取文件 %s 失败：%w", filePath, err)
	}

	var maps []MapInfo
	if err := json.Unmarshal(data, &maps); err != nil {
		return nil, fmt.Errorf("解析 JSON 文件 %s 失败：%w", filePath, err)
	}

	return maps, nil
}

// getCurrentMaps 获取当前地图列表
func getCurrentMaps() ([]MapInfo, error) {
	return getMapList("current")
}

// getFormerMaps 获取历史地图列表
func getFormerMaps() ([]MapInfo, error) {
	return getMapList("former")
}
