package connection

import (
		"time"
		"fmt"
		h "client/util/header"
	)

//*Timer for clients active status
func (c *Client) checkActive(amtTime time.Duration, pause chan bool) {

	timer := time.NewTimer(amtTime * time.Minute)

	for {
		select {
		case <-c.Status:
			// Reset the timer when there is activity
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(amtTime * time.Minute)

		case <-timer.C:
			// No activity for 1 minute, close the connection
			fmt.Println(h.I, "Connection closed due to inactivity")
			c.Conn.Close()
			return
		case <-pause:
			fmt.Println(h.I, "timer paused")
			<-pause
			fmt.Println(h.I, "timer resumed")
		}
	}
}
