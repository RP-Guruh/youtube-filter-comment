package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/goravel/framework/facades"
	"google.golang.org/api/option"
)

// CheckJudolComments menggunakan Gemini AI untuk mendeteksi komentar judi online.
// Input: map[comment_id]comment_text
// Output: map[comment_id]is_judol (true jika mengandung judol)
func CheckJudolComments(comments map[string]string) map[string]bool {
	// 1. INISIALISASI MAP HASIL DAN CEK INPUT KOSONG
	results := make(map[string]bool)
	if len(comments) == 0 {
		return results
	}

	// 2. AMBIL API KEY DARI CONFIG ATAU ENV
	apiKey := facades.Config().GetString("app.gemini_api_key")
	if apiKey == "" {
		apiKey = fmt.Sprintf("%v", facades.Config().Env("GEMINI_API_KEY", ""))
	}

	if apiKey == "" {
		log.Println("[GEMINI] API Key not found (GEMINI_API_KEY)")
		return results
	}

	// 3. BUAT CLIENT BARU UNTUK KONEKSI KE GOOGLE AI
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Printf("[GEMINI] Failed to create client: %v", err)
		return results
	}
	defer client.Close()

	// 4. SETTING MODEL YANG DIGUNAKAN (FLASH LITE) DAN FORMAT OUTPUT JSON
	// Menggunakan gemini-2.5-flash-lite untuk kecepatan dan efisiensi
	model := client.GenerativeModel("gemini-2.5-flash-lite")

	// Set response format ke JSON if supported, but manual cleaning is safer
	model.ResponseMIMEType = "application/json"

	// 5. SUSUN PROMPT INSTRUKSI DAN DAFTAR KOMENTAR YANG AKAN DIPERIKSA
	// Construct the prompt
	var builder strings.Builder
	builder.WriteString("Tugas: Identifikasi komentar yang mempromosikan judi online (judol).\n")
	builder.WriteString("Kriteria Judol: Promosi situs judi, slot gacor, depo/wd, ajakan bergabung ke grup judi, atau spam kata kunci terkait judi.\n")
	builder.WriteString("Tujuan: Hapus komentar yang berkaitan dengan judol.\n\n")
	builder.WriteString("Format Output: Harus berupa JSON valid object map dengan format { \"ID_KOMENTAR\": boolean }.\n")
	builder.WriteString("Gunakan ID_KOMENTAR yang diberikan dari daftar di bawah sebagai KEY.\n")
	builder.WriteString("Value true jika itu judol, false jika aman.\n")
	builder.WriteString("Contoh Output: { \"comment_1\": true, \"comment_2\": false }\n\n")
	builder.WriteString("Daftar Komentar:\n")

	for id, text := range comments {
		cleanText := strings.ReplaceAll(text, "\n", " ")
		builder.WriteString(fmt.Sprintf("- [%s]: %s\n", id, cleanText))
	}

	// 6. KIRIM PERMINTAAN GENERATE CONTENT KE AI
	resp, err := model.GenerateContent(ctx, genai.Text(builder.String()))
	if err != nil {
		log.Printf("[GEMINI] Failed to generate content: %v", err)
		return results
	}

	// 7. VALIDASI APAKAH AI MEMBERIKAN RESPON YANG VALID
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		log.Println("[GEMINI] No response candidates found")
		return results
	}

	// 8. AMBIL TEKS HASIL GENERATE DARI PART PERTAMA
	// Extract text from response
	var responseText string
	if part, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
		responseText = string(part)
	}

	// 9. BERSIHKAN TEKS DARI BLOCK MARKDOWN JSON JIKA ADA
	// Clean JSON if it's wrapped in markdown code blocks
	responseText = strings.TrimSpace(responseText)
	if strings.HasPrefix(responseText, "```") {
		responseText = strings.TrimPrefix(responseText, "```json")
		responseText = strings.TrimPrefix(responseText, "```")
		responseText = strings.TrimSuffix(responseText, "```")
		responseText = strings.TrimSpace(responseText)
	}

	// 10. KONVERSI (UNMARSHAL) HASIL JSON DARI AI KE DALAM MAP GO
	err = json.Unmarshal([]byte(responseText), &results)
	if err != nil {
		// FALLBACK: CEK JIKA AI MENGIRIM FORMAT ARRAY BUKAN MAP
		var arrayResults []map[string]any
		if errArray := json.Unmarshal([]byte(responseText), &arrayResults); errArray == nil {
			log.Printf("[GEMINI] Received array format instead of map. Prompt might need refinement. Raw: %s", responseText)
		}
		log.Printf("[GEMINI] Failed to unmarshal response: %v. Raw: %s", err, responseText)
	}

	// 11. CETAK LOG HASIL ANALISA UNTUK DEBUGGING
	// Debug results
	log.Printf("[GEMINI] AI Analysis Results: %+v", results)

	return results
}

// GeminiAI adalah placeholder fungsi lama, mungkin bisa di-update atau dihapus
// Tergantung kebutuhan user ke depannya.
func GeminiAI(text string) (string, error) {
	apiKey := facades.Config().GetString("app.gemini_api_key")
	if apiKey == "" {
		apiKey = fmt.Sprintf("%v", facades.Config().Env("GEMINI_API_KEY", ""))
	}

	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY not found")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text(text))
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response")
	}

	if part, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
		return string(part), nil
	}

	return "", fmt.Errorf("unexpected response type")
}
