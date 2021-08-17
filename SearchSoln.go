package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/sets/treeset"
)

var mainPageSearchMap map[string][][2]int
var subPageSearchMap map[string][][2]int

type pageDetails struct {
	MaxScore    int
	QueryNum    int
	NextPageNum int
}

func main() {
	mainPageSearchMap = make(map[string][][2]int)
	subPageSearchMap = make(map[string][][2]int)
	var MetaData pageDetails
	MetaData.MaxScore = 8
	Start(&MetaData)
}

func Start(MetaData *pageDetails) {
	scanner := bufio.NewScanner(os.Stdin)
	var wg sync.WaitGroup
	defer wg.Wait()
	for scanner.Scan() {
		Inp := scanner.Text()
		fullInpArr := strings.Split(Inp, " ")
		if fullInpArr[0] == "p" || fullInpArr[0] == "P" {
			addToMap(MetaData.NextPageNum,MetaData,fullInpArr[1:], mainPageSearchMap)
			MetaData.NextPageNum++
		} else if fullInpArr[0] == "q" || fullInpArr[0] == "Q" {
			wg.Add(1)
			var OpChan chan string = make(chan string, 1000)
			go getQueryAndCompute(fullInpArr[1:], &wg, OpChan, MetaData)

		} else if fullInpArr[0] == "pp" || fullInpArr[0] == "PP" {
			addToMap(MetaData.NextPageNum,MetaData,fullInpArr[1:], subPageSearchMap)
		} else {
			return
		}

	}

}


func addToMap(PageNum int,MetaData *pageDetails, data []string, MapRef map[string][][2]int) {
	if int(len(data)) > MetaData.MaxScore {
		MetaData.MaxScore = int(len(data))
	}
	var i int = 0
	for _, val := range data {
		if MapRef[val] == nil {
			MapRef[val] = make([][2]int, 0)
		}
		var MapAdds [2]int
		MapAdds[0] = PageNum
		MapAdds[1] = i
		i++
		MapRef[val] = append(MapRef[val], MapAdds)
	}

}

func getQueryAndCompute(data []string, wg *sync.WaitGroup, OutputChannel chan string, MetaData *pageDetails) {
	defer outPrinter(OutputChannel)
	defer wg.Done()
	MetaData.QueryNum++
	findRank(data, MetaData.QueryNum, OutputChannel, MetaData)
}

func computeScore(MaxScore int, queryList []string, MapRef map[string][][2]int, LocAndScore map[int]int) {

	QueryPosScore := MaxScore
	for _, query := range queryList {
		if MapRef[query] != nil {
			for _, MapVals := range MapRef[query] {
				var LocScore int = MaxScore - MapVals[1]
				var CompScore int = LocScore * QueryPosScore
				CompScore += LocAndScore[MapVals[0]]
				LocAndScore[MapVals[0]] = CompScore
			}
		}
		QueryPosScore--
	}
}

func findRank(queryList []string, QC int, OpChan chan string, MetaData *pageDetails) {
	var LocationAndScore map[int]int = make(map[int]int)
	computeScore(MetaData.MaxScore, queryList, mainPageSearchMap, LocationAndScore)
	computeScore(MetaData.MaxScore, queryList, subPageSearchMap, LocationAndScore)
	calcScoreMap(LocationAndScore, QC, OpChan)
}

func calcScoreMap(LocAndSco map[int]int, QC int, OpChan chan string) {
	var ScoreMap = treemap.NewWithIntComparator()
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

	scorePrinter(ScoreMap, OutStr, OpChan)

}

func scorePrinter(Scoremap *treemap.Map, OutStr string, Out chan string) {
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

func outPrinter(channel chan string) {
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
