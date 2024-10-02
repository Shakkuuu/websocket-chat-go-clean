package timefmt

import "time"

// "YYYY-MM-DD HH-MM-SS"に変換
func TimeToStr(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
