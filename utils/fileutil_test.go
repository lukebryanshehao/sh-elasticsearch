/*
@Time : 2019/3/25 16:10
@Author : lukebryan
@File : fileutil
@Software: GoLand
*/
package utils

import "testing"

func TestPathExists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PathExists(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("PathExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PathExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCopyFile(t *testing.T) {
	type args struct {
		dstName string
		srcName string
	}
	tests := []struct {
		name        string
		args        args
		wantWritten int64
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWritten, err := CopyFile(tt.args.dstName, tt.args.srcName)
			if (err != nil) != tt.wantErr {
				t.Errorf("CopyFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWritten != tt.wantWritten {
				t.Errorf("CopyFile() = %v, want %v", gotWritten, tt.wantWritten)
			}
		})
	}
}

func TestCreateFile(t *testing.T) {
	type args struct {
		file    string
		content string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"test1",args{"C:/work/goworkspace/file/test1.go","package models\r\n\r\nstruct {\r\n	id uint\r\n	name string\r\n}"},true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateFile(tt.args.file, tt.args.content); got != tt.want {
				t.Errorf("CreateFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
