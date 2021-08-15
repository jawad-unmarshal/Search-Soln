package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

var SearchMap map[string][][2]int
var ScoreMap map[int][]int
var MaxScore int = 8
var QueryCount = 1
var PageCount int = 0

func main() {
	SearchMap = make(map[string][][2]int)
	GrabInputAndRun()
}

func GrabInputAndRun() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		Inp := scanner.Text()
		fullInpArr := strings.Split(Inp, " ")
		// fmt.Println(fullInpArr)
		if fullInpArr[0] == "p" || fullInpArr[0] == "P" {
			// fmt.Println("Page")
			SetUpDataMap(fullInpArr[1:])
		} else if fullInpArr[0] == "q" || fullInpArr[0] == "Q" {
			// fmt.Println("Query")
			GetQueryAndCompute(fullInpArr[1:])
		} else if fullInpArr[0] == "pp" || fullInpArr[0] == "PP" {
			fmt.Println("Subpage")
		} else {
			return
		}

	}
}

func SetUpDataMap(data []string) {
	if int(len(data)) > MaxScore {
		MaxScore = int(len(data))
	}
	var i int = 0
	// i = 0

	for _, val := range data {
		if SearchMap[val] == nil {
			SearchMap[val] = make([][2]int, 0)
		}
		var MapAdds [2]int
		MapAdds[0] = PageCount
		MapAdds[1] = i
		i++
		SearchMap[val] = append(SearchMap[val], MapAdds)
	}
	PageCount++
}

func GetQueryAndCompute(data []string) {
	fmt.Print("Q:", QueryCount)
	QueryCount++
	ComputeScores(data)
}
func ComputeScores(queryList []string) {
	var LocationAndScore map[int]int = make(map[int]int)
	// LocationAndScore =
	// var Locs [][]int
	QueryPosScore := MaxScore

	for _, query := range queryList {
		// fmt.Println(query, " found in map: ", SearchMap[query], "\nAt Location: ", SearchMap[query][0], "With penalty: ", SearchMap[query][1])
		if SearchMap[query] != nil {
			for _, MapVals := range SearchMap[query] {
				var LocScore int = MaxScore - MapVals[1]
				var CompScore int = LocScore * QueryPosScore
				CompScore += LocationAndScore[MapVals[0]]
				LocationAndScore[MapVals[0]] = CompScore
			}
		}
		QueryPosScore--
	}
	CalcScoreMap(LocationAndScore)
}

func CalcScoreMap(LocAndSco map[int]int) {
	ScoreMap = make(map[int][]int)

	for loc, sco := range LocAndSco {
		if ScoreMap[sco] == nil {
			ScoreMap[sco] = make([]int, 0)
		}
		ScoreMap[sco] = append(ScoreMap[sco], loc)
	}

	if ScoreMap != nil {
		ScorePrinter()
		// fmt.Println(ScoreMap)
	}
	ScoreMap = make(map[int][]int)
	fmt.Println()

}
func ScorePrinter() {
	ScoreArr := make([]int, len(ScoreMap))
	var stackJugaad map [int]bool = make(map[int]bool)

	var i = 0
	for scores := range ScoreMap{
		ScoreArr[i] = scores
		i++
	}
	sort.Ints(ScoreArr)
	//Iterate in reverse
	var printVal = 0
	for i := len(ScoreArr)-1;i>=0;i-- {
		// fmt.Println("For a Score of ",ScoreArr[i],"we have values: ",ScoreMap[ScoreArr[i]])
		Pages := ScoreMap[ScoreArr[i]]
		sort.Ints(Pages)
		for _,page := range Pages {
			if !stackJugaad[page] {
				stackJugaad[page] = true
				fmt.Print(" P:",(page+1))
				printVal++
			}
			if printVal >= 5 {
				return
			}
		}

	}
}
