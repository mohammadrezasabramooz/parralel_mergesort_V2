package main
import (
	"os"
	"log"
	"bufio"
	"strconv"
	"fmt"
	"time"

	"runtime"
	
	"sync"
)


func main() {
	sizeofarray:=10000000
	var array_database [] int
	array_database=make([]int,sizeofarray)
	//initialize
	fmt.Println("initializing...")
	start := time.Now()
	initializeArray(array_database)
	end := time.Now()
	fmt.Printf("initialize time %v\n",end.Sub(start))
	//sorting process
	sizeofeachcore:=sizeofarray/runtime.NumCPU()
	fmt.Println("Sorting...")
	start=time.Now()
	sort(array_database,sizeofeachcore)

	end = time.Now()
	fmt.Printf("sorting time: %v\n",end.Sub(start))
	fmt.Println("Saving new array in text file...")
	write(array_database,sizeofarray)
}

func initializeArray (array [] int){
	file, err := os.Open("array.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	for i:=0;scanner.Scan();i++  {
		array[i], _ =(strconv.Atoi(scanner.Text()))

	}
}



func sort(m []int,size_of_array_for_each_core int)[]int  {

	useThreshold := !(size_of_array_for_each_core < 0);

	size := len(m)
	middle := size / 2

	if size <= 1 {
		return m
	}

	var left, right []int;

	sortInNewRoutine:= size > size_of_array_for_each_core && useThreshold;

	if !sortInNewRoutine{
		left = sort(m[:middle], size_of_array_for_each_core)
		right = sort(m[middle:],size_of_array_for_each_core)
	} else{
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer func() { wg.Done() }()
			left = sort(m[:middle],size_of_array_for_each_core)

		}()

		go func() {
			defer func() { wg.Done()}()
			right = sort(m[middle:], size_of_array_for_each_core)
		}()

		wg.Wait()
	}

	return merge(left, right)

	}



func merge(left, right  []int)[]int  {
	var result []int
	for len(left) > 0 || len(right) > 0 {
		if len(left) > 0 && len(right) > 0 {
			if left[0] <= right[0] {
				result = append(result, left[0])
				left = left[1:]
			} else {
				result = append(result, right[0])
				right = right[1:]
			}
		} else if len(left) > 0 {
			result = append(result, left[0])
			left = left[1:]
		} else if len(right) > 0 {
			result = append(result, right[0])
			right = right[1:]
		}
	}

	return result
}

var (
	newFile *os.File
	err     error
)

func write(array [] int,size int)  {
	newFile, err = os.Create("sorted_array.txt")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(newFile)
	for i := 0; i < size; i++ {
		_, err=newFile.WriteString(strconv.Itoa(array[i])+"\n")
	}
	newFile.Close()
	fmt.Println("==> done creating file")


}
