package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/mileusna/useragent"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// -------------------- Structs --------------------

type ServerInfo struct {
	ServerName     string `json:"server_host_name"`
	ServerOS       string `json:"server_os"`
	ServerRAM      string `json:"server_ram_usage"`
	ServerCPU      string `json:"server_cpu_model"`
	ServerCPUUsage string `json:"server_cpu_usage"`
	ServerDisk     string `json:"server_disk_usage"`
	SystemUptime   string `json:"system_uptime"`
	AppUptime      string `json:"app_uptime"`
	ServerPublicIP string `json:"server_public_ip"`
	Environment    string `json:"environment"`
	InternetStatus string `json:"internet_status"`
	HealthScore    int    `json:"health_score"`
	Suggestion     string `json:"suggestion"`
}

type ClientInfo struct {
	UserIP      string `json:"user_ip"`
	OS          string `json:"os"`
	Browser     string `json:"browser"`
	Version     string `json:"browser_version"`
	BrowserEng  string `json:"browser_engine"`
	DeviceType  string `json:"device_type"`
	DeviceBrand string `json:"device_brand"`
	DeviceModel string `json:"device_model"`
	Country     string `json:"country"`
	City        string `json:"city"`
	ISP         string `json:"isp"`
	IsBot       bool   `json:"is_bot"`
	Fingerprint string `json:"fingerprint"`
	RawAgent    string `json:"user_agent"`
}

type DeviceDetails struct {
	Server ServerInfo `json:"server"`
	Client ClientInfo `json:"client"`
}

// -------------------- Global Variables --------------------

var (
	appStartTime      = time.Now()
	publicIPCache     = ""
	publicIPCacheTime time.Time
	internetCache     = ""
	internetCacheTime time.Time
	geoIPCache        = make(map[string]geoIPData)
	geoIPCacheMutex   = sync.Mutex{}
	cacheTTL          = 10 * time.Minute
)

type geoIPData struct {
	Country string
	City    string
	ISP     string
	Timestamp time.Time
}

// -------------------- Server Info --------------------

func GetServerInfo() ServerInfo {
	v, _ := mem.VirtualMemory()
	ramUsage := fmt.Sprintf("%.2f GB Total (Used %.1f%%)", float64(v.Total)/1e9, v.UsedPercent)

	cpuInfo, _ := cpu.Info()
	cpuModel := "Unknown"
	if len(cpuInfo) > 0 {
		cpuModel = cpuInfo[0].ModelName
	}

	cpuPercent, _ := cpu.Percent(0, false)
	cpuVal := 0.0
	if len(cpuPercent) > 0 {
		cpuVal = cpuPercent[0]
	}
	cpuUsage := fmt.Sprintf("%.2f%%", cpuVal)

	diskStat, _ := disk.Usage("/")
	diskUsage := fmt.Sprintf("%.2f GB Free (%.1f%% used)", float64(diskStat.Free)/1e9, diskStat.UsedPercent)

	uptimeSec, _ := host.Uptime()
	systemUptime := (time.Duration(uptimeSec) * time.Second).Truncate(time.Second).String()
	appUptime := time.Since(appStartTime).Truncate(time.Second).String()

	// Parallel fetch Public IP & Internet Status
	var wg sync.WaitGroup
	var publicIP string
	var internetStatus string

	wg.Add(2)
	go func() {
		defer wg.Done()
		publicIP = getPublicIPCached()
	}()
	go func() {
		defer wg.Done()
		internetStatus = checkInternetCached()
	}()
	wg.Wait()

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	healthScore := calculateHealthScore(v.UsedPercent, cpuVal, parseLatency(internetStatus))
	suggestion := "Server operating normally"
	if v.UsedPercent > 85 {
		suggestion = "High RAM usage detected"
	} else if cpuVal > 80 {
		suggestion = "High CPU usage detected"
	} else if parseLatency(internetStatus) > 500*time.Millisecond {
		suggestion = "High internet latency"
	}

	hostname, _ := os.Hostname()

	return ServerInfo{
		ServerName:     hostname,
		ServerOS:       runtime.GOOS,
		ServerRAM:      ramUsage,
		ServerCPU:      cpuModel,
		ServerCPUUsage: cpuUsage,
		ServerDisk:     diskUsage,
		SystemUptime:   systemUptime,
		AppUptime:      appUptime,
		ServerPublicIP: publicIP,
		Environment:    env,
		InternetStatus: internetStatus,
		HealthScore:    healthScore,
		Suggestion:     suggestion,
	}
}

// -------------------- Client Info --------------------

func GetClientInfo(clientIP, uaString string) ClientInfo {
	ua := useragent.Parse(uaString)

	deviceType := "desktop"
	if ua.Mobile {
		deviceType = "mobile"
	} else if ua.Tablet {
		deviceType = "tablet"
	}

	brand := ua.Device
	if brand == "" {
		brand = manualBrandDetection(uaString)
	}

	country, city, isp := getGeoIPCached(clientIP)
	engine := detectBrowserEngine(uaString)
	fingerprint := generateFingerprint(clientIP, uaString, ua.OS, ua.Name)

	return ClientInfo{
		UserIP:      clientIP,
		OS:          ua.OS,
		Browser:     ua.Name,
		Version:     ua.Version,
		BrowserEng:  engine,
		DeviceType:  deviceType,
		DeviceBrand: brand,
		DeviceModel: ua.Device,
		Country:     country,
		City:        city,
		ISP:         isp,
		IsBot:       ua.Bot,
		Fingerprint: fingerprint,
		RawAgent:    uaString,
	}
}

// -------------------- Cached Helpers --------------------

func getPublicIPCached() string {
	if time.Since(publicIPCacheTime) < cacheTTL && publicIPCache != "" {
		return publicIPCache
	}

	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("https://checkip.amazonaws.com")
	if err != nil {
		return "unknown"
	}
	defer resp.Body.Close()

	ipBytes, _ := io.ReadAll(resp.Body)
	ip := strings.TrimSpace(string(ipBytes))
	publicIPCache = ip
	publicIPCacheTime = time.Now()
	return ip
}

func checkInternetCached() string {
	if time.Since(internetCacheTime) < 3*time.Second && internetCache != "" {
		return internetCache
	}

	start := time.Now()
	client := http.Client{Timeout: 3 * time.Second}
	status := "Disconnected"
	if resp, err := client.Get("https://google.com"); err == nil {
		defer resp.Body.Close()
		latency := time.Since(start)
		status = fmt.Sprintf("Stable (%v)", latency.Truncate(time.Millisecond))
	}

	internetCache = status
	internetCacheTime = time.Now()
	return status
}

func getGeoIPCached(ip string) (string, string, string) {
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return "Unknown", "Unknown", "Unknown"
	}

	geoIPCacheMutex.Lock()
	defer geoIPCacheMutex.Unlock()

	if data, ok := geoIPCache[ip]; ok && time.Since(data.Timestamp) < cacheTTL {
		return data.Country, data.City, data.ISP
	}

	country, city, isp := fetchGeoIP(ip)
	geoIPCache[ip] = geoIPData{Country: country, City: city, ISP: isp, Timestamp: time.Now()}
	return country, city, isp
}

func fetchGeoIP(ip string) (string, string, string) {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil || parsedIP.IsLoopback() || parsedIP.IsPrivate() {
		return "Local", "Local", "Internal Network"
	}

	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "Unknown", "Unknown", "Unknown"
	}
	defer resp.Body.Close()

	var data struct {
		Status  string `json:"status"`
		Country string `json:"country"`
		City    string `json:"city"`
		ISP     string `json:"isp"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil || data.Status != "success" {
		return "Unknown", "Unknown", "Unknown"
	}

	return data.Country, data.City, data.ISP
}

// -------------------- Utility --------------------

func generateFingerprint(ip, ua, os, browser string) string {
	raw := ip + ua + os + browser
	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:])
}

func detectBrowserEngine(ua string) string {
	ua = strings.ToLower(ua)
	switch {
	case strings.Contains(ua, "chrome"):
		return "Blink"
	case strings.Contains(ua, "firefox"):
		return "Gecko"
	case strings.Contains(ua, "safari"):
		return "WebKit"
	default:
		return "Unknown"
	}
}

func calculateHealthScore(ram float64, cpu float64, latency time.Duration) int {
	score := 100
	if ram > 85 {
		score -= 20
	}
	if cpu > 80 {
		score -= 20
	}
	if latency > 500*time.Millisecond {
		score -= 10
	}
	if score < 0 {
		score = 0
	}
	return score
}

func parseLatency(status string) time.Duration {
	if strings.HasPrefix(status, "Stable") {
		var d time.Duration
		fmt.Sscanf(status, "Stable (%v)", &d)
		return d
	}
	return 0
}

func manualBrandDetection(ua string) string {
	ua = strings.ToLower(ua)
	switch {
	case strings.Contains(ua, "iphone"):
		return "Apple"
	case strings.Contains(ua, "samsung") || strings.Contains(ua, "sm-"):
		return "Samsung"
	case strings.Contains(ua, "xiaomi") || strings.Contains(ua, "redmi"):
		return "Xiaomi"
	case strings.Contains(ua, "oppo") || strings.Contains(ua, "cph"):
		return "Oppo"
	case strings.Contains(ua, "vivo"):
		return "Vivo"
	case strings.Contains(ua, "infinix"):
		return "Infinix"
	case strings.Contains(ua, "tecno"):
		return "Tecno"
	case strings.Contains(ua, "itel"):
		return "Itel"
	default:
		return "Generic Device"
	}
}