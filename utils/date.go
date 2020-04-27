package utils

import "fmt"

const (
    Minute = 60
    Hour   = 60 * Minute
    Day    = 24 * Hour
    Week   = 7 * Day
    Month  = 30 * Day
    Year   = 12 * Month
)

func ComputeTimeDiff(diff int64) (int64, string) {
    diffStr := ""
    switch {
    case diff <= 0:
        diff = 0
        diffStr = "刚刚"
    case diff < 2:
        diff = 0
        diffStr = "1 秒"
    case diff < 1*Minute:
        diffStr = fmt.Sprintf("%d 秒", diff)
        diff = 0

    case diff < 2*Minute:
        diff -= 1 * Minute
        diffStr = "1 分钟"
    case diff < 1*Hour:
        diffStr = fmt.Sprintf("%d 分钟", diff/Minute)
        diff -= diff / Minute * Minute

    case diff < 2*Hour:
        diff -= 1 * Hour
        diffStr = "1 小时"
    case diff < 1*Day:
        diffStr = fmt.Sprintf("%d 小时", diff/Hour)
        diff -= diff / Hour * Hour

    case diff < 2*Day:
        diff -= 1 * Day
        diffStr = "1 天"
    case diff < 1*Week:
        diffStr = fmt.Sprintf("%d 天", diff/Day)
        diff -= diff / Day * Day

    case diff < 2*Week:
        diff -= 1 * Week
        diffStr = "1 星期"
    case diff < 1*Month:
        diffStr = fmt.Sprintf("%d 星期", diff/Week)
        diff -= diff / Week * Week

    case diff < 2*Month:
        diff -= 1 * Month
        diffStr = "1 月"
    case diff < 1*Year:
        diffStr = fmt.Sprintf("%d 月", diff/Month)
        diff -= diff / Month * Month

    case diff < 2*Year:
        diff -= 1 * Year
        diffStr = "1 年"
    default:
        diffStr = fmt.Sprintf("%d 年", diff/Year)
        diff = 0
    }
    return diff, diffStr
}
