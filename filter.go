package main

const vkPrefixID = 658 // "ВКонтакте" =

func FilterThreads(threads []Thread) []Thread {
	var filtered []Thread

	for _, thread := range threads {
		if thread.ThreadIsClosed {
			continue
		}

		filtered = append(filtered, thread)

		//skip := false
		//for _, prefix := range thread.ThreadPrefixes {
		//	if prefix.PrefixID == vkPrefixID {
		//		skip = true
		//		break
		//	}
		//}
		//
		//if !skip {
		//	filtered = append(filtered, thread)
		//}
	}

	return filtered
}
