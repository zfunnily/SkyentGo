package components

import (
	"fmt"
	"testing"
	"time"
)

func PRINT() {
	fmt.Println("12312312312")
}

func TestTimeWheel_Start(t *testing.T) {
	TimeOut(1*time.Second, func() {
		fmt.Println("12312313123")
	})
	select {}
}
