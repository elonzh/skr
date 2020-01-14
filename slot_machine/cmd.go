package slot_machine

import (
	"container/ring"
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCommand(_ *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "slot_machine",
		Version: "v20200114",
		RunE: func(cmd *cobra.Command, args []string) error {
			tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
			s, err := tcell.NewScreen()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return err
			}

			if err = s.Init(); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return err
			}
			l := []string{"0️⃣", "1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣", "7️⃣", "8️⃣", "9️⃣"}
			r := ring.New(len(l))
			for i := 0; i < len(l); i++ {
				r.Value = i
				r = r.Next()
			}
			for i := 0; i < 23; i++ {
				for j := 0; j < 3; j++ {
					s.SetContent(0, j, []rune(l[r.Value.(int)])[0], []rune(l[r.Value.(int)])[1:], tcell.StyleDefault)
					r = r.Next()
				}
				s.Show()
				time.Sleep(1000 * time.Millisecond)
				r = r.Move(1 - 3)
				fmt.Println()
			}

			s.Fini()
			return nil
		},
	}
	return cmd
}
