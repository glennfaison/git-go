package ls_tree

import (
	"reflect"
	"testing"
)

func TestParseTreeObjectFromString(t *testing.T) {
	type args struct {
		file_content string
	}
	tests := []struct {
		name string
		args args
		want []TreeObjectEntry
	}{
		{
			name: "test1",
			args: args{
				file_content: "tree 136\x00040000 folder1\x007f21f4d392c2d79987c1100644 file1\x00d51d366274410103d3ec040000 folder2\x009f94e0d364e0ea74f579100644 file2\x00d8329fc1cc938780ffdd",
			},
			want: []TreeObjectEntry{
				{
					mode: "040000", name: "folder1", _20byteSha: "7f21f4d392c2d79987c1",
				},
				{
					mode: "100644", name: "file1", _20byteSha: "d51d366274410103d3ec",
				},
				{
					mode: "040000", name: "folder2", _20byteSha: "9f94e0d364e0ea74f579",
				},
				{
					mode: "100644", name: "file2", _20byteSha: "d8329fc1cc938780ffdd",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseTreeObjectFromString(tt.args.file_content)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTreeObjectFromString() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}
