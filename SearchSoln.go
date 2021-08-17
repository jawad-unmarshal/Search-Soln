package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/sets/treeset"
)

var SearchMap map[string][][2]int
var SubPageMap map[string][][2]int

// var ScoreMap map[int][]int
var MaxScore int = 8
var QueryCount = 0
var PageCount int = 0

//Grab all queries and then compute?

func main() {
	// treemap.NewWithIntComparator()
	SearchMap = make(map[string][][2]int)
	SubPageMap = make(map[string][][2]int)
	GrabInputAndRun()
}

func GrabInputAndRun() {
	scanner := bufio.NewScanner(os.Stdin)
	var wg sync.WaitGroup

	// defer OutPrinter(OpChan)
	// defer fmt.Println(strings.Repeat("billi meow", 10))
	defer wg.Wait()
	// fulStr := ""
	for scanner.Scan() {
		Inp := scanner.Text()
		fullInpArr := strings.Split(Inp, " ")
		// fmt.Println(fullInpArr)
		if fullInpArr[0] == "p" || fullInpArr[0] == "P" {
			// fmt.Println("Page")
			SetUpDataMap(fullInpArr[1:])
		} else if fullInpArr[0] == "q" || fullInpArr[0] == "Q" {
			// fmt.Println("Query")
			wg.Add(1)
			var OpChan chan string = make(chan string, 1000)
			go GetQueryAndCompute(fullInpArr[1:], &wg, OpChan)
			// defer OutPrinter(OpChan)
			// go GetQueryAndCompute(fullInpArr[1:])

		} else if fullInpArr[0] == "pp" || fullInpArr[0] == "PP" {
			// fmt.Println("Subpage")
			SetUpSubPageMap(fullInpArr[1:])
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

// var mut sync.Mutex
func SetUpSubPageMap(data []string) {
	if int(len(data)) > MaxScore {
		MaxScore = int(len(data))
	}
	var i int = 0
	for _, val := range data {
		if SubPageMap[val] == nil {
			SubPageMap[val] = make([][2]int, 0)
		}
		var MapAdds [2]int
		MapAdds[0] = PageCount - 1
		MapAdds[1] = i
		i++
		SubPageMap[val] = append(SubPageMap[val], MapAdds)
	}
	// fmt.Println(SubPageMap)
}

func GetQueryAndCompute(data []string, wg *sync.WaitGroup, OutputChannel chan string) {
	// func GetQueryAndCompute(data []string) {
	defer OutPrinter(OutputChannel)
	defer wg.Done()

	QueryCount++
	ComputeScores(data, QueryCount, OutputChannel)
}

func ComputeScores(queryList []string, QC int, OpChan chan string) {
	var LocationAndScore map[int]int = make(map[int]int)
	QueryPosScore := MaxScore

	for _, query := range queryList {
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
	LocationAndScore = ComputeSubPageScores(queryList, LocationAndScore)
	CalcScoreMap(LocationAndScore, QC, OpChan)
}

func ComputeSubPageScores(queryList []string, LocationAndScore map[int]int) map[int]int {
	QueryPosScore := MaxScore

	for _, query := range queryList {
		if SubPageMap[query] != nil {
			for _, MapVals := range SubPageMap[query] {
				var LocScore int = MaxScore - MapVals[1]
				var CompScore float64 = float64(LocScore * QueryPosScore)
				CompScore = 0.1 * CompScore
				CompScore += float64(LocationAndScore[MapVals[0]])
				LocationAndScore[MapVals[0]] = int(CompScore)
			}
		}
		QueryPosScore--
	}
	return LocationAndScore

}

func calcScoreMap_Deprecated(LocAndSco map[int]int, QC int, OpChan chan string) {
	var ScoreMap = make(map[int][]int)
	// OP := make(chan string, 100)
	for loc, sco := range LocAndSco {
		if ScoreMap[sco] == nil {
			ScoreMap[sco] = make([]int, 0)
		}
		ScoreMap[sco] = append(ScoreMap[sco], loc)
	}
	OutStr := fmt.Sprint("\nQ", QC, ": ")
	ScorePrinter_Deprecated(ScoreMap, OutStr, OpChan)

}

func CalcScoreMap(LocAndSco map[int]int, QC int, OpChan chan string) {
	var ScoreMap = treemap.NewWithIntComparator()
	// OP := make(chan string, 100)
	OutStr := fmt.Sprint("\nQ", QC, ": ")
	for loc, sco := range LocAndSco {
		var val = treeset.NewWithIntComparator()
		val1, hasKey := ScoreMap.Get(sco)
		if !hasKey {
			val1 = treeset.NewWithIntComparator()
		}
		val = val1.(*treeset.Set)
		val.Add(loc)
		ScoreMap.Put(sco, val)

	}
	// fmt.Printf("%#v",ScoreMap.Values())

	ScorePrinter(ScoreMap, OutStr, OpChan)

}

func ScorePrinter(Scoremap *treemap.Map, OutStr string, Out chan string) {
	var ind int = 1
	ScoreIterator := Scoremap.Iterator()
OUTER_FOR:
	for ScoreIterator.End(); ScoreIterator.Prev(); {
		Locations := ScoreIterator.Value()
		LocationItr := Locations.(*treeset.Set).Iterator()

		for LocationItr.Next() {
			OutStr += fmt.Sprint(" P:", (LocationItr.Value().(int) + 1))
			ind++
			if ind > 5 {
				break OUTER_FOR
			}
		}
	}
	Out <- OutStr

}

func OutPrinter(channel chan string) {
	for {
		select {
		case PageValue := <-channel:
			fmt.Print(PageValue)
		default:
			fmt.Println("")
			return
		}
	}
}

func ScorePrinter_Deprecated(ScoreMap map[int][]int, OutStr string, Out chan string) {
	ScoreArr := make([]int, len(ScoreMap))
	var stackJugaad map[int]bool = make(map[int]bool)

	var i = 0
	for scores := range ScoreMap {
		ScoreArr[i] = scores
		i++
	}
	sort.Ints(ScoreArr)
	var printVal = 0
OUTER_FOR:
	for i := len(ScoreArr) - 1; i >= 0; i-- {
		Pages := ScoreMap[ScoreArr[i]]
		sort.Ints(Pages)
		for _, page := range Pages {
			if !stackJugaad[page] {
				stackJugaad[page] = true
				OutStr += fmt.Sprint(" P:", (page + 1))
				printVal++
			}
			if printVal >= 5 {
				break OUTER_FOR
			}
		}

	}
	Out <- OutStr
}
