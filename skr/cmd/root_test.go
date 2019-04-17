package cmd

import (
	"os"
	"path"
	"testing"
)

func TestExecute(t *testing.T) {
	if len(rootCmd.Commands()) != 2 {
		t.Fatalf("expect 2 commands")
	}
	rootCmd.SetArgs([]string{"douyin", "-u", "http://v.douyin.com/jf2teV/"})
	if err := rootCmd.Execute(); err != nil {
		t.Error(err)
	}
	filename := path.Join(os.TempDir(), "test_gaoxiaojob.boltdb")
	defer func() {
		if err := os.Remove(filename); err != nil && err != os.ErrNotExist {
			t.Fatal(err)
		}
	}()
	rootCmd.SetArgs([]string{"gaoxiaojob", "-s", filename, "https://oapi.dingtalk.com/robot/send?access_token=fee17fe946f196a86b99e68cef74d8311c5fc020dc1db3454c174df3c58a4409"})
	if err := rootCmd.Execute(); err != nil {
		t.Error(err)
	}
}
