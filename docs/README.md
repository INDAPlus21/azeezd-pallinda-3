# Answers to "Task 1 - Matching Behaviour"

### What happens if you remove the `go-command` from the `Seek` call in the `main` function?
The calls on the names will be in order because without the concurrency every sender sends their message which then is received by the next name in the list. With the concurrency another name might receive the message before the next person in line and therefore the names do not appear in any particular order.

### What happens if you switch the declaration `wg := new(sync.WaitGroup)` to `var wg sync.WaitGroup` and the parameter `wg *sync.WaitGroup` to `wg sync.WaitGroup`?
This makes the `WaitGroup` be passed by value not reference because we omitted the pointer. This makes the `wg.Done` inside the `Seek` redundant because it is not changing the `wg` in the main goroutine making it wait forever and creating a deadlock.

### What happens if you remove the buffer on the channel match?
This will create a deadlock because there will be sent data nowhere to be received in the currently running goroutines. Buffering makes it possible to have a message sit in the channel later to be received which in this case is the `select` after the `wg.Wait`

### What happens if you remove the default-case from the case-statement in the main function?
The default case was there if all sendings were received which in this program's case is never going to happen due to the fact that there is an odd number of names that send and receive. Thus, removing the `select`'s default case will correspond to no change at all.

# Analysis of SingleWorker vs MapReduce for word count
The time is collected by running `go run words.go` several times finding an approximate interval of the `average time/run` value.

| Variant       | Avg time/run (ms)     |
|---            |---                    |
| Single Worker |   ~9-10ms             |
| Map Reduce    |   ~7-8ms              |