package stringx

import (
	"fmt"
	"testing"
)

func Test_isDigits(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			"test",
			args{
				data: "1",
			},
			true,
		},
		{
			"test",
			args{
				data: "2",
			},
			true,
		},
		{
			"test",
			args{
				data: "22",
			},
			false,
		},
		{
			"test",
			args{
				data: "a",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDigit(tt.args.data); got != tt.want {
				t.Errorf("isDigits() = %v, want %v", got, tt.want)
			}
		})
	}
}


func ExampleMustUnquote() {
	test := "aaa\\u6388\\u6743\\u9a8c\\u8bc1\\u5931\\u8d25\\uff1a\\u8bf7\\u6c42\\u53d7\\u9650"
	fmt.Println(UnicodeToUtf8(test))
	//Output:
}