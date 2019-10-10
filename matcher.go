package marker

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type MatcherFunc func(string) Match

type state int

const (
	ready state = iota
	s
	m
	t
	w
	f
	sa
	sumo
	th
	tu
	we
	fr
	sat
	ready_for_day
	thu
	thurtue
	wed
	satu
	d
	da
	day
)

// Match contains information about found patterns by MatcherFunc
type Match struct {
	Template string
	Patterns []string
}

// MatchAll returns a MatcherFunc that matches all patterns in given string
func MatchAll(pattern string) MatcherFunc {
	return func(str string) Match {
		count := strings.Count(str, pattern)
		return Match{
			Template: strings.ReplaceAll(str, pattern, "%s"),
			Patterns: fillSlice(make([]string, count), pattern),
		}
	}
}

// MatchN returns a MatcherFunc that matches first n patterns in given string
func MatchN(pattern string, n int) MatcherFunc {
	return func(str string) Match {
		count := min(n, strings.Count(str, pattern))
		return Match{
			Template: strings.Replace(str, pattern, "%s", n),
			Patterns: fillSlice(make([]string, count), pattern),
		}
	}
}

// MatchRegexp returns a MatcherFunc that matches regexp in given string
func MatchRegexp(r *regexp.Regexp) MatcherFunc {
	return func(str string) Match {
		return Match{
			Template: r.ReplaceAllString(str, "%s"),
			Patterns: r.FindAllString(str, -1),
		}
	}
}

// MatchSurrounded takes in characters surrounding a given expected match and returns the match findings
func MatchSurrounded(charOne string, charTwo string) MatcherFunc {
	return func(str string) Match {
		quoteCharOne := regexp.QuoteMeta(charOne)
		quoteCharTwo := regexp.QuoteMeta(charTwo)
		matchPattern := fmt.Sprintf("%s[^%s]*%s", quoteCharOne, quoteCharOne, quoteCharTwo)
		r, _ := regexp.Compile(matchPattern)
		return MatchRegexp(r)(str)
	}
}

// MatchBracketSurrounded is a helper utility for easy matching of bracket surrounded text
func MatchBracketSurrounded() MatcherFunc {
	return MatchSurrounded("[", "]")
}

// MatchParensSurrounded is a helper utility for easy matching text surrounded in parentheses
func MatchParensSurrounded() MatcherFunc {
	return MatchSurrounded("(", ")")
}

var daysOfWeek = [14]string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday",
	"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

// MatchDaysOfWeek returns a MatcherFunc that matches days of the week in given string
func MatchDaysOfWeek() MatcherFunc {
	return func(str string) Match {
		patternMatchIndexes := make(map[int]string)
		for _, day := range daysOfWeek {
			for strings.Contains(str, day) {
				matchIndex := strings.Index(str, day)
				str = strings.Replace(str, day, "%s", 1)
				patternMatchIndexes[matchIndex] = day
			}
		}
		matchIndexes := make([]int, 0, len(patternMatchIndexes))
		for matchKey := range patternMatchIndexes {
			matchIndexes = append(matchIndexes, matchKey)
		}
		sort.Ints(matchIndexes)
		pattern := make([]string, 0, len(patternMatchIndexes))
		for _, index := range matchIndexes {
			pattern = append(pattern, patternMatchIndexes[index])
		}
		return Match{
			Template: str,
			Patterns: pattern,
		}
	}
}

func MatchDaysOfWeek2() MatcherFunc {
	return func(str string) Match {
		state := ready
		strbuf := bytes.NewBufferString(str).Bytes()
		var startInd int
		var matchStr []string
		out := str
		var offset int
		for ind, char := range strbuf {
			switch state {
			case ready:
				startInd = ind
				switch char {
				case 'm', 'M':
					state = m
				case 't', 'T':
					state = t
				case 'w', 'W':
					state = w
				case 'f', 'F':
					state = f
				case 's', 'S':
					state = s
				default:
					state = ready
				}
			case s:
				switch char {
				case 'u':
					state = sumo
				case 'a':
					state = sa
				default:
					state = ready
				}
			case m:
				switch char {
				case 'o':
					state = sumo
				default:
					state = ready

				}
			case t:
				switch char {
				case 'u':
					state = tu
				case 'h':
					state = th
				default:
					state = ready
				}
			case w:
				switch char {
				case 'e':
					state = we
				default:
					state = ready
				}
			case f:
				switch char {
				case 'r':
					state = fr
				default:
					state = ready
				}
			case sa:
				switch char {
				case 't':
					state = sat
				default:
					state = ready
				}
			case sumo:
				switch char {
				case 'n':
					state = ready_for_day
				default:
					state = ready
				}
			case th:
				switch char {
				case 'u':
					state = thu
				default:
					state = ready
				}
			case tu:
				switch char {
				case 'e':
					state = thurtue
				default:
					state = ready
				}
			case we:
				switch char {
				case 'd':
					state = wed
				default:
					state = ready
				}
			case fr:
				switch char {
				case 'i':
					state = ready_for_day
				default:
					state = ready
				}
			case sat:
				switch char {
				case 'u':
					state = satu
				default:
					state = ready
				}
			case ready_for_day:
				switch char {
				case 'd':
					state = d
				default:
					state = ready
				}
			case thu:
				switch char {
				case 'r':
					state = thurtue
				default:
					state = ready
				}
			case thurtue:
				switch char {
				case 's':
					state = ready_for_day
				default:
					state = ready
				}
			case wed:
				switch char {
				case 'n':
					state = tu
				default:
					state = ready
				}
			case satu:
				switch char {
				case 'r':
					state = ready_for_day
				default:
					state = ready
				}
			case d:
				switch char {
				case 'a':
					state = da
				default:
					state = ready
				}
			case da:
				switch char {
				case 'y':
					state = day
				default:
					state = ready
				}
			case day:
				matchStr = append(matchStr, str[startInd:ind])
				out = out[:startInd-offset] + "%s" + out[ind-offset:]
				offset += ind - startInd - 2
				state = ready
			}
		}
		return Match{
			Template: out,
			Patterns: matchStr,
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func fillSlice(s []string, v string) []string {
	for i := range s {
		s[i] = v
	}
	return s
}
