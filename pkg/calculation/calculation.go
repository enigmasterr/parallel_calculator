package calculation

import (
	"strconv"
)

func Calc(expression string) (float64, error) {
	prior := map[string]int{
		"(": 0,
		")": 1,
		"+": 2,
		"-": 2,
		"*": 3,
		"/": 3,
	}
	var ans []string
	var st []string
	num := ""
	for _, sim := range expression {
		if sim == '(' {
			if len(num) > 0 {
				ans = append(ans, num)
			}
			st = append(st, string(sim))
		} else {
			if sim == '+' || sim == '-' || sim == '*' || sim == '/' {
				if num != "" {
					ans = append(ans, num)
					num = ""
				}
				if len(st) == 0 {
					st = append(st, string(sim))
				} else {
					if prior[string(sim)] > prior[st[len(st)-1]] {
						st = append(st, string(sim))
					} else {
						for len(st) > 0 && prior[string(sim)] <= prior[st[len(st)-1]] {
							ans = append(ans, st[len(st)-1])
							st = st[:len(st)-1]
						}
						st = append(st, string(sim))
					}
				}
			} else if sim == ')' {
				if len(num) > 0 {
					ans = append(ans, num)
					num = ""
				}
				for st[len(st)-1] != "(" {
					ans = append(ans, st[len(st)-1])
					st = st[:len(st)-1]
				}
				st = st[:len(st)-1]
			} else {
				num += string(sim)
			}
		}
	}
	if num != "" {
		ans = append(ans, num)
		num = ""
	}
	for len(st) > 0 {
		if st[len(st)-1] == "(" || st[len(st)-1] == ")" {
			return 0, ErrInvalidExpression
		} else {
			ans = append(ans, st[len(st)-1])
			st = st[:len(st)-1]
		}
	}
	var stk []float64
	for _, v := range ans {
		if v == "+" || v == "-" || v == "*" || v == "/" {
			if len(stk) < 2 {
				return 0, ErrInvalidExpression
			}
			a := stk[len(stk)-1]
			stk = stk[:len(stk)-1]
			b := stk[len(stk)-1]
			stk = stk[:len(stk)-1]
			if v == "+" {
				stk = append(stk, b+a)
			} else if v == "-" {
				stk = append(stk, b-a)
			} else if v == "*" {
				stk = append(stk, b*a)
			} else if v == "/" {
				if a == 0 {
					return 0, ErrDivisionByZero
				}
				stk = append(stk, b/a)
			}
		} else {
			num, _ := strconv.ParseFloat(v, 64)
			stk = append(stk, num)
		}
	}
	if len(stk) != 1 {
		return 0, ErrInvalidExpression
	}
	return stk[0], nil
}
