package analyze

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"testing"
)

func TestColors(t *testing.T) {
	filepath.Walk("images", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		var topColors []uint32
		topColors, err = CrackFile(path, 3)
		if err != nil {
			fmt.Println("error:", err)
			err = nil
		} else {
			fmt.Printf("__file:%v\n", path)
			for idx, this := range topColors {
				fmt.Printf("Top%v: %v\t%v\n", idx+1, HexColor(this), RBA(this))
			}
		}

		return nil
	})
}
