package song

import (
	"errors"
	"os/exec"
	"strings"
)

type Song map[string]string


var errFailedToGetNowPlaying = errors.New("failed to get now playing")

func executeAppleScript(script string) (string, error) {
	output, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func scriptResultParser(scriptOutput string) (Song) {
	//class:URL track, id:22363, index:1, name:"나였으면", persistent ID:"423C7D352913BCE3", database ID:22355, date added:date "2023년 6월 27일 화요일 오후 8:33:50", time:"4:22", duration:262.0, artist:"나윤권", album artist:"", composer:"", album:"중독", genre:"K-Pop", track count:0, track number:0, disc count:0, disc number:0, volume adjustment:0, year:2004, comment:"", EQ:"", kind:"", media kind:song, enabled:true, start:0.0, finish:262.0, played count:0, skipped count:0, compilation:false, rating:0, bpm:0, grouping:"", bookmarkable:false, bookmark:0.0, shufflable:true, category:"", description:"", episode number:0, unplayed:true, sort name:"", sort album:"", sort artist:"", sort composer:"", sort album artist:"", release date:date "2004년 7월 7일 수요일 오후 9:00:00", loved:false, disliked:false, album loved:false, album disliked:false, work:"", movement:"", movement number:0, movement count:0
	result := make(Song)
	outputSlice := strings.Split(scriptOutput, ", ")
	
	for _, output := range outputSlice {
		slice := strings.Split(output, ":")

		for i := 0; i < len(slice); i++ {
			slice[i] = strings.ReplaceAll(slice[i], "\n", "")
		}

		// Join the rest of the slice. except the first element [0]
		if len(slice) > 2 {
			slice[1] = strings.Join(slice[1:], ":")
		}

		if len(slice) < 2 {
			continue
		}

		result[slice[0]] = slice[1]
	}

	return result
}

func GetNowPlaying() (Song, error) {
	// script := `tell application "Music" to tell artwork 1 of current track
	// 			return data
	// 			end tell`
	// // artwork, err := executeAppleScript(script)
	// // if err != nil {
	// // 	return "", err
	// // }

	script := `tell application "Music" to get properties of current track`
	scriptOutput, err := executeAppleScript(script)
	if err != nil {
		return Song{}, errFailedToGetNowPlaying
	}
	properties := scriptResultParser(scriptOutput)
	return properties, nil
}
