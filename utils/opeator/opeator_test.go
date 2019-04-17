package opeator

import "testing"

func TestIntIn(t *testing.T) {
	type args struct {
		src int
		to  []int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name:"test1",
			args:args{
				src: 7,
				to:  []int{7,8},
			},
			want:true,
		},
		{
			name:"test2",
			args:args{
				src: 7,
				to:  []int{8,9},
			},
			want:false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntIn(tt.args.src, tt.args.to); got != tt.want {
				t.Errorf("IntIn() = %v, want %v", got, tt.want)
			}
		})
	}
}
