// Portions of this file are derived from pluralize 8.0.0
// (https://github.com/plurals/pluralize, MIT License, Copyright (c) 2013
// Blake Embrey), the library @dbml/core's schema.rb importer uses.
// Modified: ported to Go (RE2 regular expressions, same rule order).

package importer

import (
	"regexp"
	"strings"
)

type pluralRule struct {
	re   *regexp.Regexp
	repl string
}

var (
	pluralRules      []pluralRule
	singularRules    []pluralRule
	uncountableWords = map[string]bool{}
	irregularPlurals = map[string]string{} // plural -> single
	irregularSingles = map[string]string{} // single -> plural
)

func addRule(rules *[]pluralRule, pattern, repl string) {
	*rules = append(*rules, pluralRule{regexp.MustCompile("(?i)" + pattern), repl})
}

func init() {
	for _, p := range [][2]string{
		{"I", "we"}, {"me", "us"}, {"he", "they"}, {"she", "they"}, {"them", "them"}, {"myself", "ourselves"},
		{"yourself", "yourselves"}, {"itself", "themselves"}, {"herself", "themselves"}, {"himself", "themselves"},
		{"themself", "themselves"}, {"is", "are"}, {"was", "were"}, {"has", "have"}, {"this", "these"}, {"that", "those"},
		{"echo", "echoes"}, {"dingo", "dingoes"}, {"volcano", "volcanoes"}, {"tornado", "tornadoes"}, {"torpedo", "torpedoes"},
		{"genus", "genera"}, {"viscus", "viscera"}, {"stigma", "stigmata"}, {"stoma", "stomata"}, {"dogma", "dogmata"},
		{"lemma", "lemmata"}, {"schema", "schemata"}, {"anathema", "anathemata"}, {"ox", "oxen"}, {"axe", "axes"},
		{"die", "dice"}, {"yes", "yeses"}, {"foot", "feet"}, {"eave", "eaves"}, {"goose", "geese"}, {"tooth", "teeth"},
		{"quiz", "quizzes"}, {"human", "humans"}, {"proof", "proofs"}, {"carve", "carves"}, {"valve", "valves"},
		{"looey", "looies"}, {"thief", "thieves"}, {"groove", "grooves"}, {"pickaxe", "pickaxes"}, {"passerby", "passersby"},
	} {
		s, pl := strings.ToLower(p[0]), strings.ToLower(p[1])
		irregularSingles[s] = pl
		irregularPlurals[pl] = s
	}
	for _, r := range [][2]string{
		{`s?$`, "s"},
		{`[^\x{0000}-\x{007F}]$`, "$0"},
		{`([^aeiou]ese)$`, "$1"},
		{`(ax|test)is$`, "$1es"},
		{`(alias|[^aou]us|t[lm]as|gas|ris)$`, "$1es"},
		{`(e[mn]u)s?$`, "$1s"},
		{`([^l]ias|[aeiou]las|[ejzr]as|[iu]am)$`, "$1"},
		{`(alumn|syllab|vir|radi|nucle|fung|cact|stimul|termin|bacill|foc|uter|loc|strat)(?:us|i)$`, "$1i"},
		{`(alumn|alg|vertebr)(?:a|ae)$`, "$1ae"},
		{`(seraph|cherub)(?:im)?$`, "$1im"},
		{`(her|at|gr)o$`, "$1oes"},
		{`(agend|addend|millenni|dat|extrem|bacteri|desiderat|strat|candelabr|errat|ov|symposi|curricul|automat|quor)(?:a|um)$`, "$1a"},
		{`(apheli|hyperbat|periheli|asyndet|noumen|phenomen|criteri|organ|prolegomen|hedr|automat)(?:a|on)$`, "$1a"},
		{`sis$`, "ses"},
		{`(?:(kni|wi|li)fe|(ar|l|ea|eo|oa|hoo)f)$`, "$1$2ves"},
		{`([^aeiouy]|qu)y$`, "$1ies"},
		{`([^ch][ieo][ln])ey$`, "$1ies"},
		{`(x|ch|ss|sh|zz)$`, "$1es"},
		{`(matr|cod|mur|sil|vert|ind|append)(?:ix|ex)$`, "$1ices"},
		{`\b((?:tit)?m|l)(?:ice|ouse)$`, "$1ice"},
		{`(pe)(?:rson|ople)$`, "$1ople"},
		{`(child)(?:ren)?$`, "$1ren"},
		{`eaux$`, "$0"},
		{`m[ae]n$`, "men"},
		{`^thou$`, "you"},
	} {
		addRule(&pluralRules, r[0], r[1])
	}
	for _, r := range [][2]string{
		{`s$`, ""},
		{`(ss)$`, "$1"},
		{`(wi|kni|(?:after|half|high|low|mid|non|night|[^\w]|^)li)ves$`, "$1fe"},
		{`(ar|(?:wo|[ae])l|[eo][ao])ves$`, "$1f"},
		{`ies$`, "y"},
		{`\b([pl]|zomb|(?:neck|cross)?t|coll|faer|food|gen|goon|group|lass|talk|goal|cut)ies$`, "$1ie"},
		{`\b(mon|smil)ies$`, "$1ey"},
		{`\b((?:tit)?m|l)ice$`, "$1ouse"},
		{`(seraph|cherub)im$`, "$1"},
		{`(x|ch|ss|sh|zz|tto|go|cho|alias|[^aou]us|t[lm]as|gas|(?:her|at|gr)o|[aeiou]ris)(?:es)?$`, "$1"},
		{`(analy|diagno|parenthe|progno|synop|the|empha|cri|ne)(?:sis|ses)$`, "$1sis"},
		{`(movie|twelve|abuse|e[mn]u)s$`, "$1"},
		{`(test)(?:is|es)$`, "$1is"},
		{`(alumn|syllab|vir|radi|nucle|fung|cact|stimul|termin|bacill|foc|uter|loc|strat)(?:us|i)$`, "$1us"},
		{`(agend|addend|millenni|dat|extrem|bacteri|desiderat|strat|candelabr|errat|ov|symposi|curricul|quor)a$`, "$1um"},
		{`(apheli|hyperbat|periheli|asyndet|noumen|phenomen|criteri|organ|prolegomen|hedr|automat)a$`, "$1on"},
		{`(alumn|alg|vertebr)ae$`, "$1a"},
		{`(cod|mur|sil|vert|ind)ices$`, "$1ex"},
		{`(matr|append)ices$`, "$1ix"},
		{`(pe)(rson|ople)$`, "$1rson"},
		{`(child)ren$`, "$1"},
		{`(eau)x?$`, "$1"},
		{`men$`, "man"},
	} {
		addRule(&singularRules, r[0], r[1])
	}
	for _, w := range []string{"adulthood", "advice", "agenda", "aid", "aircraft", "alcohol", "ammo", "analytics", "anime",
		"athletics", "audio", "bison", "blood", "bream", "buffalo", "butter", "carp", "cash", "chassis", "chess", "clothing",
		"cod", "commerce", "cooperation", "corps", "debris", "diabetes", "digestion", "elk", "energy", "equipment",
		"excretion", "expertise", "firmware", "flounder", "fun", "gallows", "garbage", "graffiti", "hardware",
		"headquarters", "health", "herpes", "highjinks", "homework", "housework", "information", "jeans", "justice",
		"kudos", "labour", "literature", "machinery", "mackerel", "mail", "media", "mews", "moose", "music", "mud", "manga",
		"news", "only", "personnel", "pike", "plankton", "pliers", "police", "pollution", "premises", "rain", "research",
		"rice", "salmon", "scissors", "series", "sewage", "shambles", "shrimp", "software", "species", "staff", "swine",
		"tennis", "traffic", "transportation", "trout", "tuna", "wealth", "welfare", "whiting", "wildebeest", "wildlife",
		"you"} {
		uncountableWords[w] = true
	}
	// regexp uncountables are added as identity rules (pluralize.addUncountableRule)
	for _, p := range []string{`pok[eé]mon$`, `[^aeiou]ese$`, `deer$`, `fish$`, `measles$`, `o[iu]s$`, `pox$`, `sheep$`} {
		addRule(&pluralRules, p, "$0")
		addRule(&singularRules, p, "$0")
	}
}

func restoreCase(word, token string) string {
	if word == token {
		return token
	}
	if word == strings.ToLower(word) {
		return strings.ToLower(token)
	}
	if word == strings.ToUpper(word) {
		return strings.ToUpper(token)
	}
	if word != "" && token != "" && word[:1] == strings.ToUpper(word[:1]) {
		return strings.ToUpper(token[:1]) + strings.ToLower(token[1:])
	}
	return strings.ToLower(token)
}

// applyRule mirrors replace(word, rule): the first match only.
func applyRule(word string, r pluralRule) string {
	m := r.re.FindStringSubmatchIndex(word)
	if m == nil {
		return word
	}
	var groups []string
	for i := 0; i < len(m); i += 2 {
		if m[i] < 0 {
			groups = append(groups, "")
		} else {
			groups = append(groups, word[m[i]:m[i+1]])
		}
	}
	result := regexp.MustCompile(`\$(\d{1,2})`).ReplaceAllStringFunc(r.repl, func(s string) string {
		n := 0
		for _, c := range s[1:] {
			n = n*10 + int(c-'0')
		}
		if n < len(groups) {
			return groups[n]
		}
		return ""
	})
	match := word[m[0]:m[1]]
	var rc string
	if match == "" {
		prev := ""
		if m[0] > 0 {
			prev = word[m[0]-1 : m[0]]
		}
		rc = restoreCase(prev, result)
	} else {
		rc = restoreCase(match, result)
	}
	return word[:m[0]] + rc + word[m[1]:]
}

func inflect(word string, replaceMap, keepMap map[string]string, rules []pluralRule) string {
	token := strings.ToLower(word)
	if _, ok := keepMap[token]; ok {
		return restoreCase(word, token)
	}
	if r, ok := replaceMap[token]; ok {
		return restoreCase(word, r)
	}
	if token == "" || uncountableWords[token] {
		return word
	}
	for i := len(rules) - 1; i >= 0; i-- {
		if rules[i].re.MatchString(word) {
			return applyRule(word, rules[i])
		}
	}
	return word
}

func pluralize(word string) string {
	return inflect(word, irregularSingles, irregularPlurals, pluralRules)
}

func singularize(word string) string {
	return inflect(word, irregularPlurals, irregularSingles, singularRules)
}
