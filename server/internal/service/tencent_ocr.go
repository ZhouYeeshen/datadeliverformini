package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// callTencentHandwritingOCR calls Tencent Cloud GeneralHandwritingOCR API.
// Supports both printed and handwritten Chinese text recognition.
func (s *OCRService) callTencentHandwritingOCR(imageData []byte) (string, error) {
	if s.tcSecretID == "" || s.tcSecretKey == "" {
		return "", fmt.Errorf("腾讯云密钥未配置")
	}

	// Compress large images
	data := imageData
	if len(data) > 1024*1024 {
		data = compressForOCR(data)
	}

	b64 := base64.StdEncoding.EncodeToString(data)
	body, _ := json.Marshal(map[string]string{"ImageBase64": b64})

	host := "ocr.tencentcloudapi.com"
	service := "ocr"
	action := "GeneralHandwritingOCR"
	version := "2018-11-19"
	region := s.tcRegion
	algorithm := "TC3-HMAC-SHA256"
	timestamp := time.Now().Unix()
	date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")

	// Step 1: canonical request
	httpRequestMethod := "POST"
	canonicalURI := "/"
	canonicalQueryString := ""
	canonicalHeaders := "content-type:application/json\n" + "host:" + host + "\n"
	signedHeaders := "content-type;host"
	hashedPayload := sha256Hex(body)
	canonicalRequest := strings.Join([]string{
		httpRequestMethod, canonicalURI, canonicalQueryString,
		canonicalHeaders, signedHeaders, hashedPayload,
	}, "\n")

	// Step 2: string to sign
	credentialScope := date + "/" + service + "/tc3_request"
	hashedCanonicalRequest := sha256Hex([]byte(canonicalRequest))
	stringToSign := strings.Join([]string{
		algorithm, fmt.Sprintf("%d", timestamp), credentialScope, hashedCanonicalRequest,
	}, "\n")

	// Step 3: signature
	secretDate := hmacSHA256([]byte("TC3"+s.tcSecretKey), date)
	secretService := hmacSHA256(secretDate, service)
	secretSigning := hmacSHA256(secretService, "tc3_request")
	signature := hex.EncodeToString(hmacSHA256(secretSigning, stringToSign))

	// Step 4: authorization header
	authorization := fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		algorithm, s.tcSecretID, credentialScope, signedHeaders, signature)

	req, err := http.NewRequest("POST", "https://"+host, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", host)
	req.Header.Set("X-TC-Action", action)
	req.Header.Set("X-TC-Version", version)
	req.Header.Set("X-TC-Timestamp", fmt.Sprintf("%d", timestamp))
	req.Header.Set("X-TC-Region", region)
	req.Header.Set("Authorization", authorization)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("腾讯云 OCR 请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var result struct {
		Response struct {
			Error struct {
				Code    string `json:"Code"`
				Message string `json:"Message"`
			} `json:"Error"`
			HandwritingResult []struct {
				DetectedText string  `json:"DetectedText"`
				Confidence   float64 `json:"Confidence"`
			} `json:"HandwritingResult"`
			TextDetections []struct {
				DetectedText string  `json:"DetectedText"`
				Confidence   float64 `json:"Confidence"`
			} `json:"TextDetections"`
			RequestID string `json:"RequestId"`
		} `json:"Response"`
	}
	json.Unmarshal(respBody, &result)

	if result.Response.Error.Code != "" {
		return "", fmt.Errorf("腾讯云 OCR: %s (%s)", result.Response.Error.Message, result.Response.Error.Code)
	}

	var lines []string
	for _, item := range result.Response.HandwritingResult {
		if item.DetectedText != "" {
			lines = append(lines, item.DetectedText)
		}
	}
	for _, item := range result.Response.TextDetections {
		if item.DetectedText != "" {
			lines = append(lines, item.DetectedText)
		}
	}
	return strings.Join(lines, "\n"), nil
}

func sha256Hex(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

func hmacSHA256(key []byte, data string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(data))
	return mac.Sum(nil)
}
