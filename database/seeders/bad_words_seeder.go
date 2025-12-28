package seeders

import (
	"goravel/app/models"

	"github.com/goravel/framework/facades"
)

type BadWordsSeeder struct {
}

// Signature The name and signature of the seeder.
func (s *BadWordsSeeder) Signature() string {
	return "BadWordsSeeder"
}

// Run executes the seeder logic.
func (s *BadWordsSeeder) Run() error {

	judiBrand := []string{"𝐃⁠⁠ 𝐎 ⁠𝐑⁠‌𝐀⁠ 𝟕⁠⁠ 𝟕", "️𝐏𝐋𝐔𝗧𝗢8̲8", "𝓟𝓤𝓛𝓐 𝓤𝓦𝓘𝓝", "ℙ𝕌𝕃𝔸𝕌𝕎𝕀ℕ", "ＰＵＬＡＵＷＩＮ", "𝐌0𝐍𝐀𝟒𝐃", "pulauwin", "pluto88", "mona4d", "dora 77", "pragmatic", "pgsoft", "habanero", "microgaming", "slot88", "idnlive", "toto", "togel", "sbobet", "maxwin", "mahjongways", "olympus", "starlight", "gateofolympus"}

	judiTerms := []string{"gacor", "rungkad", "jp", "jackpot", "sensasional", "wd", "depo", "deposit", "withdraw", "scatters", "tumble", "multiplier", "rtp", "pola", "martingale"}

	pornWords := []string{"vcs", "bokep", "sange", "colmek", "openbo", "lendir", "pasutri", "sepong", "crot", "peju", "nenen"}

	for _, word := range judiBrand {
		s.saveWord(s.generateLeetRegex(word), "judi_brand", true, 1.0)
	}
	for _, word := range judiTerms {
		s.saveWord(s.generateLeetRegex(word), "judi_terms", true, 1.0)
	}
	for _, word := range pornWords {
		s.saveWord(s.generateLeetRegex(word), "pornografi", true, 1.0)
	}

	sensitif := []string{"anjing", "babi", "bangsat", "tolol", "goblok", "cebong", "kadrun", "kampret", "pki", "rezim"}
	for _, word := range sensitif {
		s.saveWord(word, "provokasi", false, 0.7)
	}

	promosi := []string{
		`(?i)cek\s?bio`,
		`(?i)klik\s?link`,
		`(?i)hubungi\s?wa`,
		`(?i)t[.]me\/`,
		`(?i)bit[.]ly\/`,
		`(?i)wa[.]me\/`,
	}
	for _, p := range promosi {
		s.saveWord(p, "spam_link", true, 0.9)
	}

	return nil
}

func (s *BadWordsSeeder) generateLeetRegex(word string) string {
	replacements := map[rune]string{
		'a': "[a|4|@]", 'i': "[i|1|!|l]", 'o': "[o|0]", 'e': "[e|3]",
		's': "[s|5|$]", 'g': "[g|9|6]", 't': "[t|7]", 'b': "[b|8]",
	}
	pattern := "(?i)"
	for _, r := range word {
		if val, ok := replacements[r]; ok {
			pattern += val
		} else {
			pattern += string(r)
		}
		pattern += `[\s\._\-]?`
	}
	return pattern
}

func (s *BadWordsSeeder) saveWord(word string, cat string, isRegex bool, score float64) {
	var bw models.BadWord
	facades.Orm().Query().FirstOrCreate(&bw, models.BadWord{
		Word: word, Category: cat, IsRegex: isRegex, SeverityScore: score, IsActive: true,
	})
}
