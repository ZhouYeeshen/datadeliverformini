package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type OCRResult struct {
	RestaurantRevenue     float64 `json:"restaurant_revenue"`
	RetailRevenue         float64 `json:"retail_revenue"`
	AccommodationRevenue  float64 `json:"accommodation_revenue"`
	TobaccoAlcoholRevenue float64 `json:"tobacco_alcohol_revenue"`
	OtherRevenue          float64 `json:"other_revenue"`
	TotalRevenue          float64 `json:"total_revenue"`
	RawText               string  `json:"raw_text"`
	Engine                string  `json:"engine"`
	Confidence            string  `json:"confidence"`
}

type OCRService struct {
	httpClient    *http.Client
	wechatAppID   string
	wechatSecret  string
	tcSecretID    string
	tcSecretKey   string
	tcRegion      string
	tokenMu       sync.Mutex
	accessToken   string
	tokenExpireAt time.Time
}

func NewOCRService(appID, appSecret, tcSecretID, tcSecretKey, tcRegion string) *OCRService {
	return &OCRService{
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		wechatAppID:  appID,
		wechatSecret: appSecret,
		tcSecretID:   tcSecretID,
		tcSecretKey:  tcSecretKey,
		tcRegion:     tcRegion,
	}
}

// ProcessImage reads a file, runs OCR, and parses the result.
func (s *OCRService) ProcessImage(filePath string) (*OCRResult, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}
	rawText, engine := s.recognize(data)

	result := s.parseReceipt(rawText)
	result.RawText = rawText
	result.Engine = engine
	result.TotalRevenue = s.calcTotal(result)
	return result, nil
}

// ParseText parses raw OCR text into structured revenue categories (no image needed).
func (s *OCRService) ParseText(rawText string) *OCRResult {
	result := s.parseReceipt(rawText)
	result.RawText = rawText
	result.Engine = "text-parse"
	result.TotalRevenue = s.calcTotal(result)
	return result
}

// ProcessImageFromBase64 decodes base64, runs OCR, and parses the result.
func (s *OCRService) ProcessImageFromBase64(b64 string) (*OCRResult, error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("base64 解码失败: %w", err)
	}
	rawText, engine := s.recognize(data)

	result := s.parseReceipt(rawText)
	result.RawText = rawText
	result.Engine = engine
	result.TotalRevenue = s.calcTotal(result)
	return result, nil
}

// recognize tries Tencent Cloud OCR first (handwritten + printed), then WeChat OCR fallback (printed only).
func (s *OCRService) recognize(imageData []byte) (string, string) {
	// 1. Tencent Cloud GeneralHandwritingOCR (supports handwritten + printed)
	if s.tcSecretID != "" && s.tcSecretKey != "" {
		text, err := s.callTencentHandwritingOCR(imageData)
		if err != nil {
			fmt.Printf("[OCR] Tencent Cloud error: %v\n", err)
		} else if text != "" {
			return text, "tencent-handwriting"
		} else {
			fmt.Println("[OCR] Tencent Cloud returned empty text")
		}
	} else {
		fmt.Println("[OCR] Tencent Cloud credentials not configured")
	}

	// 2. WeChat OCR API fallback (printed only, requires API permission)
	if s.wechatAppID != "" && s.wechatSecret != "" {
		token, err := s.getAccessToken()
		if err != nil {
			fmt.Printf("[OCR] WeChat access token error: %v\n", err)
			return "", "none"
		}
		text, err := s.callWeChatOCR(token, imageData)
		if err != nil {
			fmt.Printf("[OCR] WeChat OCR error: %v\n", err)
		} else if text != "" {
			return text, "wechat-ocr"
		}
	}

	return "", "none"
}

// ---- WeChat OCR API (cv/ocr/comm) ----

func (s *OCRService) callWeChatOCR(token string, imageData []byte) (string, error) {
	// Compress large images to stay under WeChat API limits
	data := imageData
	if len(data) > 1024*1024 {
		data = compressForOCR(data)
	}

	b64 := base64.StdEncoding.EncodeToString(data)
	body, _ := json.Marshal(map[string]string{"img": b64})

	req, err := http.NewRequest("POST",
		"https://api.weixin.qq.com/cv/ocr/comm?access_token="+token,
		bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("微信 OCR 请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var result struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		Items   []struct {
			Text string `json:"text"`
		} `json:"items"`
	}
	json.Unmarshal(respBody, &result)

	if result.ErrCode != 0 {
		return "", fmt.Errorf("微信 OCR: %s (code=%d)", result.ErrMsg, result.ErrCode)
	}

	var lines []string
	for _, item := range result.Items {
		if item.Text != "" {
			lines = append(lines, item.Text)
		}
	}
	return strings.Join(lines, "\n"), nil
}

// ---- WeChat access token (cached, 2-hour TTL) ----

func (s *OCRService) getAccessToken() (string, error) {
	s.tokenMu.Lock()
	defer s.tokenMu.Unlock()

	if s.accessToken != "" && time.Now().Before(s.tokenExpireAt) {
		return s.accessToken, nil
	}

	url := fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		s.wechatAppID, s.wechatSecret)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("获取微信 access_token 失败: %w", err)
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}
	json.Unmarshal(data, &result)

	if result.ErrCode != 0 {
		return "", fmt.Errorf("获取 access_token: %s (code=%d)", result.ErrMsg, result.ErrCode)
	}

	s.accessToken = result.AccessToken
	s.tokenExpireAt = time.Now().Add(time.Duration(result.ExpiresIn-300) * time.Second)

	return s.accessToken, nil
}

// === Parsing engine ===

func (s *OCRService) parseReceipt(text string) *OCRResult {
	result := &OCRResult{}
	text = normalizeChineseNums(text)

	lines := splitTextLines(text)

	s.extractByCates(lines, result)
	if s.allZero(result) {
		s.extractByAmountRows(lines, result)
	}
	if s.allZero(result) {
		s.extractByTotalBackfill(text, result)
	}

	return result
}

func (s *OCRService) calcTotal(r *OCRResult) float64 {
	cateSum := math.Round((r.RestaurantRevenue+r.RetailRevenue+
		r.AccommodationRevenue+r.TobaccoAlcoholRevenue+
		r.OtherRevenue)*100) / 100

	rawSum := sumAllAmounts(r.RawText)
	if rawSum > cateSum {
		return rawSum
	}
	return cateSum
}

func (s *OCRService) allZero(r *OCRResult) bool {
	return r.RestaurantRevenue == 0 && r.RetailRevenue == 0 &&
		r.AccommodationRevenue == 0 && r.TobaccoAlcoholRevenue == 0 &&
		r.OtherRevenue == 0
}

func (s *OCRService) extractByCates(lines []string, r *OCRResult) {
	type cateRule struct {
		target *float64
		keys   []string
	}
	rules := []cateRule{
		{&r.RestaurantRevenue, []string{"餐饮", "餐厅", "餐坎", "餐", "伙食", "食堂", "烧烤", "火锅", "外卖", "菜品", "饭菜", "小吃", "早点", "快餐", "堂食"}},
		{&r.RetailRevenue, []string{"零售", "百货", "日用", "商品", "零瘦", "销货", "日杂", "超市", "便利店", "批发"}},
		{&r.AccommodationRevenue, []string{"住宿", "客房", "住", "旅店", "旅馆", "酒店", "民宿", "房间", "房费"}},
		{&r.TobaccoAlcoholRevenue, []string{"烟酒", "烟草", "烟", "酒水", "酒", "白酒", "啤酒", "饮料"}},
		{&r.OtherRevenue, []string{"其他", "其它", "他", "杂项", "其它收入", "服务费"}},
	}

	for _, line := range lines {
		amounts := findAllAmounts(line)
		for _, rule := range rules {
			if *rule.target > 0 {
				continue
			}
			for _, key := range rule.keys {
				if strings.Contains(line, key) && len(amounts) > 0 {
					*rule.target = amounts[0]
					break
				}
			}
		}
	}
}

func (s *OCRService) extractByAmountRows(lines []string, r *OCRResult) {
	type row struct {
		idx    int
		line   string
		amount float64
	}
	var rows []row
	for i, line := range lines {
		amounts := findAllAmounts(line)
		if len(amounts) > 0 {
			rows = append(rows, row{i, line, amounts[len(amounts)-1]})
		}
	}
	if len(rows) == 0 {
		return
	}

	// Sort descending by amount
	for i := 0; i < len(rows)-1; i++ {
		for j := i + 1; j < len(rows); j++ {
			if rows[j].amount > rows[i].amount {
				rows[i], rows[j] = rows[j], rows[i]
			}
		}
	}

	catePtrs := []*float64{&r.RestaurantRevenue, &r.RetailRevenue,
		&r.AccommodationRevenue, &r.TobaccoAlcoholRevenue, &r.OtherRevenue}

	eligible := rows
	if len(rows) > 1 {
		eligible = rows[1:] // skip largest (likely total)
	}

	for _, row := range eligible {
		if s.tryAssignCateByText(row.line, catePtrs, row.amount) {
			continue
		}
		for _, ptr := range catePtrs {
			if *ptr == 0 {
				*ptr = row.amount
				break
			}
		}
	}
}

func (s *OCRService) tryAssignCateByText(line string, ptrs []*float64, amount float64) bool {
	cateKeys := [][]string{
		{"餐饮", "餐", "伙食"},
		{"零售", "百货", "商品", "零"},
		{"住宿", "客房", "旅店", "酒店"},
		{"烟酒", "烟草", "烟", "酒"},
		{"其他", "其它", "杂"},
	}
	for i, keys := range cateKeys {
		for _, k := range keys {
			if strings.Contains(line, k) && *ptrs[i] == 0 {
				*ptrs[i] = amount
				return true
			}
		}
	}
	return false
}

func (s *OCRService) extractByTotalBackfill(text string, r *OCRResult) {
	re := regexp.MustCompile(`(?i)(合[计针针]|总[计针额]|小[计针])\s*[:：]?\s*(\d+\.?\d*)`)
	match := re.FindStringSubmatch(text)
	if len(match) >= 3 {
		total, _ := strconv.ParseFloat(match[2], 64)
		r.OtherRevenue = total
	}
}

// === Utilities ===

func compressForOCR(data []byte) []byte {
	if len(data) < 100*1024 {
		return data
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return data
	}
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	maxDim := 1200
	if w <= maxDim && h <= maxDim && len(data) < 1024*1024 {
		return data
	}
	var nw, nh int
	if w > h {
		nw, nh = maxDim, h*maxDim/w
	} else {
		nw, nh = w*maxDim/h, maxDim
	}
	if nw < 1 {
		nw = 1
	}
	if nh < 1 {
		nh = 1
	}
	dst := image.NewRGBA(image.Rect(0, 0, nw, nh))
	for y := 0; y < nh; y++ {
		srcY := y * h / nh
		for x := 0; x < nw; x++ {
			srcX := x * w / nw
			dst.Set(x, y, img.At(srcX, srcY))
		}
	}
	var out bytes.Buffer
	jpeg.Encode(&out, dst, &jpeg.Options{Quality: 75})
	return out.Bytes()
}

func normalizeChineseNums(text string) string {
	replacements := map[string]string{
		"〇": "0", "零": "0", "一": "1", "二": "2", "三": "3", "四": "4",
		"五": "5", "六": "6", "七": "7", "八": "8", "九": "9",
		"两": "2", "十": "10", "百": "00", "千": "000", "万": "0000",
		"．": ".", "，": ",", "：": ":",
	}
	for old, nu := range replacements {
		text = strings.ReplaceAll(text, old, nu)
	}
	return text
}

func splitTextLines(text string) []string {
	lines := strings.Split(text, "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}
	return result
}

func sumAllAmounts(text string) float64 {
	re := regexp.MustCompile(`\d+\.?\d*`)
	matches := re.FindAllString(text, -1)
	var total float64
	for _, m := range matches {
		if len(m) > 9 {
			continue
		}
		if len(m) == 4 {
			if v, _ := strconv.Atoi(m); v >= 2020 && v <= 2099 {
				continue
			}
		}
		if v, err := strconv.ParseFloat(m, 64); err == nil && v > 10 && v < 100000000 {
			total += v
		}
	}
	return math.Round(total*100) / 100
}

func findAllAmounts(line string) []float64 {
	re := regexp.MustCompile(`\d+\.?\d*`)
	matches := re.FindAllString(line, -1)
	var result []float64
	for _, m := range matches {
		if len(m) > 8 {
			continue
		}
		if len(m) == 4 {
			if yr, _ := strconv.Atoi(m); yr >= 2020 && yr <= 2099 {
				continue
			}
		}
		if v, err := strconv.ParseFloat(m, 64); err == nil && v > 0 && v < 100000000 {
			result = append(result, v)
		}
	}
	return result
}
