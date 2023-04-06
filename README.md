# Phone-Number formatter

## Prerequisites
Need go installed, and get dependencies with `go get ./...` from the root dir

## To Build
build with `go build -o ./server ./cmd/main.go` from the root dir

## To Run
run `./server <<path-to-data-file>>` from the root after building. `./server example.data` works with the given example file. 

## To Test

### Manual
After running server, send get commands to `http://localhost:8081/v1/value?<<uuid>>`

### Unit tests
from the root, simply run `go test ./...`

## Explanations

### Tech Stack
I am most familiar with Golang and is seems suited to the task.  I am using typical HTTP Get method to fetch a key

### Deploying
We would set up a kubernetes and docker framework and deploy to a cloud platform like GCP or AWS.  Would need more infra to keep it from being DDOSed though, such as rate limiting on IP or having a sign-up/login workflow

### Questions 
* How much data can your server handle? How could you improve it so it can handle even larger datasets?
  * As much as memory allows.  To improve, I would move to postgres for a large amount of data or BigTable for huge amounts 
* How many milliseconds it takes for the client to get a response on average? How could you improve the latency?
  * locally it takes about 5ms.  In memory golang Maps powered by hashtables are the key-value lookup fastest I can thing of.  As for real internet application, only way to speed this up would be less cumbersome protocols than HTTP and/or having CDNs duplicating the dataset closer to the request origin.
* What are some failure patterns that you can anticipate?
    * Datasets that are to large to contain in memory.  Data sets with invalid UUID or duplicate, or that want to include newline characters.  DDOS attacks as there is no authentication or anything.  Holding things in memory on a single server is very fragile.  

### Assumptions and improvements
I am under the assumption to NOT use an exists market Key-Data store tech like Bigtable or even Postgresql, so I am doing it in memory with a golang map, which is limited by available memory.  This obviously has limitations for size and durability, and to product-tize it, I would likely use Postgresql for ease of understanding, or BigTable if we were expecting HUGE data sets. 


The testing is also not super thorough, and I'm not handline duplicate UUIDs in the given data file, which are things to address with more time.  

Also would want to return status codes for different errors instead of always 500.  Starting implementing but requires a little more time.  