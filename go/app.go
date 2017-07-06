package app

import (
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/pata/", handlePata)
}

func handlePata(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Fprintln(w, r.Form)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	a := r.Form.Get("a")
	b := r.Form.Get("b")
	fmt.Fprintln(w, a)
	fmt.Fprintln(w, b)
	ans := generatePatatokakushi(a, b)
	fmt.Fprintln(w, ans)
}

// https://mrekucci.blogspot.jp/2015/07/dont-abuse-mathmax-mathmin.html
func Max(x, y int) int {
    if x > y {
        return x
    }
    return y
}

func generatePatatokakushi(a, b string) string {
/*
	if a == nil || b == nil {
		return ""
	}
*/
	var ans string;
	ra := []rune(a)
	rb := []rune(b)
	maxLen := Max(len(ra), len(rb))
	for i := 0; i < maxLen; i++ {
        if i < len(ra) {
			ans += string(ra[i])
		}
        if i < len(rb) {
			ans += string(rb[i])
		}
    }
	return ans
}
