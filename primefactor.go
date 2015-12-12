package main

import (
    "fmt"
    "net/http"
	"log"
	"strconv"
)

//CREDIT: http://edapx.com/2014/04/12/how-to-get-the-prime-factors-of-a-number-in-golang/

// Generate numbers until the limit max.
// after the 2, all the prime numbers are odd
// Send a channel signal when the limit is reached
func Generate(max int, ch chan<- int) {
    ch <- 2
    for i := 3; i <= max; i += 2 {
        ch <- i
    }
    ch <- -1 // signal that the limit is reached
}

// Copy the values from channel 'in' to channel 'out',
// removing those divisible by 'prime'.
func Filter(in <-chan int, out chan<- int, prime int) {
    for i := <-in; i != -1; i = <-in {
        if i%prime != 0 {
            out <- i
        }
    }
    out <- -1
}

func CalcPrimeFactors(number_to_factorize int) []int {
    rv := []int{}
    ch := make(chan int)
    go Generate(number_to_factorize, ch)
    for prime := <-ch; (prime != -1) && (number_to_factorize > 1); prime = <-ch {
        for number_to_factorize%prime == 0 {
            number_to_factorize = number_to_factorize / prime
            rv = append(rv, prime)
        }
        ch1 := make(chan int)
        go Filter(ch, ch1, prime)
        ch = ch1
    }
    return rv
}

func main(){
    http.HandleFunc("/nsa",func (w http.ResponseWriter, r *http.Request){
        defer r.Body.Close()
        nstr := r.FormValue("num")
        if nstr == "" {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("Please include num paramter"))
            return
        }
        num, err := strconv.Atoi(string(nstr))
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(fmt.Sprintf("error parsing number: %v\n",err)))
            return
        }
        factors := CalcPrimeFactors(num)
        w.Write([]byte(fmt.Sprintf("%v\n",factors)))
    })

    http.HandleFunc("/health", func ( w http.ResponseWriter, r *http.Request){
        w.Write([]byte("ok"))
    })

    log.Fatal(http.ListenAndServe(":8080",nil))
}
