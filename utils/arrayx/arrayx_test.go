package arrayx

import "testing"

func TestFindStringInArray(t *testing.T) {
	type args struct {
		obj    string
		target []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				obj:    "abc",
				target: []string{"abc", "bcd"},
			},
			want: true,
		},
		{
			name: "test",
			args: args{
				obj:    "abc",
				target: []string{"cbd", "bcd"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindStringInArray(tt.args.obj, tt.args.target); got != tt.want {
				t.Errorf("FindStringInArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
