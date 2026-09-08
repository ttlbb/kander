package install

// previousOfficialHashes are SHA-256 digests of official rule files from earlier eras: the
// copies shipped by install.sh / install.ps1 immediately before self-install, the final Chinese
// originals shipped before the rules became English-only, and the final English-only copies
// shipped before the Chinese translation returned. Doctor uses them to tell outdated official
// copies from local edits when kander-rules-state.json is absent, regardless of language.
var previousOfficialHashes = map[string][]string{
	"KANDER-AGENTS.md":              {"2b4acbb0278eca67217ec1b8a7c0fc023f425992895b6594435becde4ee16656", "ab82abaf8dfcdc274644af77e38815be1a3584d77eecb23dad8de3b58f8f0f82"},
	"KANDER-BASE-RULES.md":          {"519b5b65e3c147b98185b0374e5c51365e3dfbe19be1d57b73694a78183e4935", "a31cefdcfd6530563354255eb18479045d487fc356757e593518c3c30209ed61", "22ebb594841bc5167e923e8e0df96b46b69149a7a7a3593e514166bb656ae8bd"},
	"KANDER-CODE-RULES.md":          {"3b1f82396e4f646721ef8ae4c57ba98e0caa9731faea8ba42309677e583525b4", "9f7a2570202312d0e237caf0ac39305d7f574348ef2b221955b9a52520b0e742"},
	"KANDER-COLLABORATION-RULES.md": {"9a5be4f188c7a5ee5721d074420027c016fc2d4139da18976ebc4c3f8da958fc", "00c45c2117485b7d1d0db0414d2a6271e9fa62f89746891361453a31406ce6ba"},
	"KANDER-GIT-RULES.md":           {"65b51a2ed0f474d275c43715998a50a08f285c469aaf696c1c6867ec79aeab75", "1e2af08df6859cb0ce64f64b27c97f691875f731fd1d6d38d3bb0515d1205fc1"},
	"KANDER-KANBAN-RULES.md":        {"aa8fe364854ccce3b82494d70cb61e53c7eeb21227a1d326cdecb39e396d5662", "c113e309a00931769400b0016f4df8e2e0d5722a9b6a278db34ec97b8f34aaa9", "934d33007713f139342ec60556bd69bf5969e14f1829cd5f5f0e51d43a61a674"},
	"KANDER-REPORTING-RULES.md":     {"5743b7f4ae40a2d665ad89e00243f65b8ddaab6b3c20833be8c100922172ac62", "687c483ba4746c7702fd25bb5fcb65daf4e1fe95fef0165c7c22ba35b7d2f529"},
	"KANDER-REVIEW-RULES.md":        {"84f31f4dcc2cfc4d44eefe811f60f523c08be3e4d62d85ba303ce54f4d177abf", "9a22020d079d00247ed3cfcf6e2cd2c9ab1894170cdb4e1663896f6781cbe938", "78b8e9b1151c8531ba1e5344e18e3c4884736082cbd05079e0e16d432d6f8e11"},
	"KANDER-TASK-GROUP-RULES.md":    {"dfcacf126b2766018d8010055fc7e9e95b88407c17c47b31299788c85d52b251", "0fbb6ca5a497d23769f5e82dcda08803b91ff18caf00d984536daed693428e12", "f642a3b9e6170087b788ed57ec77179c5e6d16446b89bd5cc4d641f8152b16d1"},
	"KANDER-TASK-INTAKE-RULES.md":   {"9596d928014a15746a8d088775045674c09b9461e077458d022cbb306a0c002d", "de14e9509e072ce91e0b889e7d15ba65e3e387442d29562fcea3bd1607ed065c"},
}

func isPreviousOfficial(name, digest string) bool {
	for _, known := range previousOfficialHashes[name] {
		if known == digest {
			return true
		}
	}
	return false
}
