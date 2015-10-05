package main


func routine() {
    c := make(chan int)
    <- c
}

