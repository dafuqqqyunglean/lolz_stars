package main

import "fmt"

func FormatPost(thread Thread, link string) string {
	return fmt.Sprintf(
		"[userids=%d;align=left][USER=%d]@%s[/USER], %s %s  [/userids]",
		thread.CreatorUserID,
		thread.CreatorUserID,
		thread.CreatorUsername,
		thread.Captcha,
		link,
	)
}
