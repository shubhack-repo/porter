package main
import (
	"fmt"
	"time"
	"net"
	"strconv"
	"sync"
	"flag"
	"os"
	"bufio"
)

func main(){
	threads := 20
	flag.IntVar(&threads,"t", 20, "Set the threads level")
	flag.Parse()
	var swg sync.WaitGroup
	jobs := make(chan string)
	for t := 0 ; t < threads ; t++{
		swg.Add(1)

		go func(){
			for dom := range jobs{
				portscan(dom)
			}
			swg.Done()
		}()
	}

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan(){
		jobs <- sc.Text()
	}

	close(jobs)

	swg.Wait()
}

func portscan(hostname string){
	concurrency := 60
	ports := make([]string,65535)
	var wg sync.WaitGroup

	for i := 1 ; i <= 65535 ; i++{
		ports[i-1] = strconv.Itoa(i)
	}

	start := 0

	for j := 1; j <= concurrency; j++ {
		end := (65535/concurrency)*j
		wg.Add(1)
		go func(host string, ports []string){
			for _, port := range ports {
        		timeout := time.Second
        		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
        		if err != nil {
            	continue
        		}
        		if conn != nil {
            	fmt.Println(net.JoinHostPort(host, port))
            	conn.Close()
        		}
    		}
    		wg.Done()
		}(hostname,ports[start:end])
	
	start=end+1
	
	}

	wg.Wait()
}
