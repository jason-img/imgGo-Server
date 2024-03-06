package templateFilter

import (
	"errors"
	"github.com/flosch/pongo2/v4"
)

// Register 注册pongo2 filter
func Register() {
	_ = pongo2.RegisterFilter("HtmlSafe", HtmlSafe)         // 输出HTML内容
	_ = pongo2.RegisterFilter("FormatDate", FormatDate)     // 时间 to 日期字符串
	_ = pongo2.RegisterFilter("FormatTime", FormatTime)     // 时间 to 日期时间字符串
	_ = pongo2.RegisterFilter("IntToSlice", IntToSlice)     // Int转slice，用于生成页码
	_ = pongo2.RegisterFilter("IntToSliceC5", IntToSliceC5) // Int转slice，用于生成页码 [i-2, i-1, i, i+1, i+2]
}

// HtmlSafe 输出HTML内容
func HtmlSafe(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
	if !in.IsString() {
		return nil, &pongo2.Error{
			OrigError: errors.New("only strings should be sent to the scream filter"),
		}
	}

	s := in.String()
	//s = string(template.HTML(s))

	return pongo2.AsSafeValue(s), nil
}

// FormatDate 时间 to 日期字符串
func FormatDate(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
	if !in.IsTime() {
		return nil, &pongo2.Error{
			OrigError: errors.New("only strings should be sent to the scream filter"),
		}
	}
	s := in.Time()

	return pongo2.AsSafeValue(s.Format("2006-01-02")), nil
}

// FormatTime 时间 to 日期时间字符串
func FormatTime(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
	if !in.IsTime() {
		return nil, &pongo2.Error{
			OrigError: errors.New("only strings should be sent to the scream filter"),
		}
	}
	s := in.Time()

	return pongo2.AsSafeValue(s.Format("2006-01-02 15:04:05")), nil
}

// IntToSlice Int转slice，用于生成页码
func IntToSlice(in *pongo2.Value, max *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
	if !in.IsInteger() {
		return in, &pongo2.Error{
			OrigError: errors.New("只有Integer 型才能转换为切片"),
		}
	}

	var s []int

	number := in.Integer()
	maxNumber := max.Integer()

	if number <= 5 {
		for i := 1; i <= number; i++ {
			if maxNumber > 0 && i > maxNumber {
				break
			}
			s = append(s, i)
		}
	} else if number == 5 {
		s = []int{1, 2, 3, 4, 5}
	} else {
		s = []int{number - 4, number - 3, number - 2, number - 1, number}
	}

	return pongo2.AsSafeValue(s), nil
}

// IntToSliceC5 Int转slice，用于生成页码 [i-2, i-1, i, i+1, i+2]
func IntToSliceC5(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
	if !in.IsInteger() {
		return in, &pongo2.Error{
			OrigError: errors.New("只有Integer 型才能转换为切片"),
		}
	}

	i := in.Integer()
	s := []int{i - 2, i - 1, i, i + 1, i + 2}

	return pongo2.AsSafeValue(s), nil
}
