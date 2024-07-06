package ls_tree

import (
	"reflect"
	"testing"

	"github.com/codecrafters-io/git-starter-go/pkg"
)

func TestParseTreeObjectFromString(t *testing.T) {
	type args struct {
		file_content string
	}
	tests := []struct {
		name string
		args args
		want []pkg.TreeObjectEntry
	}{
		{
			name: "test1",
			args: args{
				file_content: "tree 136\x00040000 folder1\x007f21f4d392c2d79987c1100644 file1\x00d51d366274410103d3ec040000 folder2\x009f94e0d364e0ea74f579100644 file2\x00d8329fc1cc938780ffdd",
			},
			want: []pkg.TreeObjectEntry{
				{
					Mode: "040000", Name: "folder1", ShaAs20Bytes: "7f21f4d392c2d79987c1",
				},
				{
					Mode: "100644", Name: "file1", ShaAs20Bytes: "d51d366274410103d3ec",
				},
				{
					Mode: "040000", Name: "folder2", ShaAs20Bytes: "9f94e0d364e0ea74f579",
				},
				{
					Mode: "100644", Name: "file2", ShaAs20Bytes: "d8329fc1cc938780ffdd",
				},
			},
		},
		{
			name: "test2",
			args: args{
				file_content: "tree 99\x0040000 donkey\x00\x04\xac\x17ͱ\xe4?\x06\x02\xae\xd4Xoc\xd8p\xb3t\xc7\a100644 horsey\x00=\u03a2\x116\xec\x14ʊ\xf3\xabG\xa0z'ȜK\xe8\xaf40000 yikes\x00D\xb1\"t\xdd\x13^J\xcd%\xe8@h^\xb4\xf8\xdaO\xbaI",
			},
			want: []pkg.TreeObjectEntry{
				{
					Mode: "040000", Name: "donkey", ShaAs20Bytes: "\x04\xac\x17ͱ\xe4?\x06\x02\xae\xd4Xoc\xd8p\xb3t\xc7\a",
				},
				{
					Mode: "100644", Name: "horsey", ShaAs20Bytes: "=\u03a2\x116\xec\x14ʊ\xf3\xabG\xa0z'ȜK\xe8\xaf",
				},
				{
					Mode: "040000", Name: "yikes", ShaAs20Bytes: "D\xb1\"t\xdd\x13^J\xcd%\xe8@h^\xb4\xf8\xdaO\xbaI",
				},
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			got := pkg.ParseTreeObjectFromString(tt.args.file_content)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTreeObjectFromString() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}
