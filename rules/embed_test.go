package rules

import (
	"strings"
	"testing"
)

func TestNamesAreEnglishMarkdownOnly(t *testing.T) {
	names := Names()
	if len(names) != 10 {
		t.Fatalf("names=%v", names)
	}
	seen := map[string]bool{}
	for _, name := range names {
		if !strings.HasSuffix(name, ".md") {
			t.Fatalf("non-markdown enumerated: %s", name)
		}
		if seen[name] {
			t.Fatalf("duplicate %s", name)
		}
		seen[name] = true
	}
	if !seen["KANDER-AGENTS.md"] || !seen["KANDER-BASE-RULES.md"] {
		t.Fatalf("missing entry files: %v", names)
	}
}

func TestChineseTranslationIsComplete(t *testing.T) {
	for _, name := range Names() {
		cn, actual, err := File(LangCN, name)
		if err != nil || actual != LangCN || len(cn) == 0 {
			t.Fatalf("%s: cn actual=%s err=%v", name, actual, err)
		}
		en, _, err := File(LangEN, name)
		if err != nil {
			t.Fatal(err)
		}
		if string(cn) == string(en) {
			t.Fatalf("%s: cn copy is identical to en", name)
		}
	}
}

func TestEnglishNeverFallsBack(t *testing.T) {
	for _, lang := range []string{LangEN, "ja", "fr", ""} {
		_, actual, err := File(lang, "KANDER-AGENTS.md")
		if err != nil || actual != LangEN {
			t.Fatalf("lang=%q actual=%s err=%v", lang, actual, err)
		}
	}
}

func TestLangFor(t *testing.T) {
	cases := map[string]string{"cn": LangCN, "en": LangEN, "ja": LangEN, "": LangEN, "fr": LangEN}
	for in, want := range cases {
		if got := LangFor(in); got != want {
			t.Fatalf("LangFor(%q)=%s want %s", in, got, want)
		}
	}
}

func TestEntryDeclaresAgentLanguage(t *testing.T) {
	for _, lang := range []string{LangCN, LangEN} {
		data, _, err := File(lang, "KANDER-AGENTS.md")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), "`agent_language`") {
			t.Fatalf("%s rules entry must tell agents to honor agent_language", lang)
		}
	}
}

func TestFileRejectsInvalidNames(t *testing.T) {
	if _, _, err := File(LangEN, "embed.go"); err == nil {
		t.Fatal("non-markdown must not be readable as a rule")
	}
	if _, _, err := File(LangEN, "../embed.go"); err == nil {
		t.Fatal("path escape must fail")
	}
	if _, _, err := File(LangCN, "missing.md"); err == nil {
		t.Fatal("missing file must fail even with fallback")
	}
}

func TestHashStable(t *testing.T) {
	digest, actual, err := Hash(LangEN, "KANDER-BASE-RULES.md")
	if err != nil || actual != LangEN || len(digest) != 64 {
		t.Fatalf("digest=%s actual=%s err=%v", digest, actual, err)
	}
	again, _, err := Hash(LangEN, "KANDER-BASE-RULES.md")
	if err != nil || again != digest {
		t.Fatalf("hash mismatch: %s vs %s", digest, again)
	}
	cn, actual, err := Hash(LangCN, "KANDER-BASE-RULES.md")
	if err != nil || actual != LangCN || cn == digest {
		t.Fatalf("cn digest=%s actual=%s err=%v", cn, actual, err)
	}
}
