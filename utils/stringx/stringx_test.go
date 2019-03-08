package stringx

import "testing"

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
			if got := isDigits(tt.args.data); got != tt.want {
				t.Errorf("isDigits() = %v, want %v", got, tt.want)
			}
		})
	}
}
