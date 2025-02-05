package scanner

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

func PortScan(target string, timeout time.Duration, maxConcurrency int) ([]int, error) {
	var openPorts []int
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 使用带缓冲的通道来控制并发
	semaphore := make(chan struct{}, maxConcurrency)
	// 增加缓冲区大小，避免channel阻塞
	resultCh := make(chan int, maxConcurrency)
	errCh := make(chan error, maxConcurrency)

	// 创建一个context来控制超时
	ctx, cancel := context.WithTimeout(context.Background(), timeout*65535)
	defer cancel()

	// 扫描端口
	for port := 1; port <= 65535; port++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()

			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }() // 扫描完毕释放信号量
			case <-ctx.Done():
				return
			}

			addr := fmt.Sprintf("%s:%d", target, port)

			// 使用context控制单个端口扫描的超时
			connCtx, connCancel := context.WithTimeout(ctx, timeout)
			defer connCancel()

			var d net.Dialer
			conn, err := d.DialContext(connCtx, "tcp", addr)
			if err != nil {
				if !errors.Is(err, context.DeadlineExceeded) {
					// 忽略错误输出
					return
				}
				return
			}

			conn.Close()

			select {
			case resultCh <- port:
			case <-ctx.Done():
			}
		}(port)
	}

	// 使用单独的goroutine来收集结果
	go func() {
		wg.Wait()
		close(resultCh)
		close(errCh)
	}()

	// 收集结果
	for {
		select {
		case port, ok := <-resultCh:
			if !ok {
				// 当 resultCh 关闭时，按顺序排序端口
				mu.Lock()
				sort.Ints(openPorts)
				mu.Unlock()
				return openPorts, nil
			}
			mu.Lock()
			openPorts = append(openPorts, port)
			mu.Unlock()

		case <-ctx.Done():
			mu.Lock()
			sort.Ints(openPorts)
			mu.Unlock()
			return openPorts, ctx.Err()
		}
	}
}

// DetectService 根据端口号返回服务名称
func DetectService(port int) string {
	services := map[int]string{
		21:    "FTP",
		22:    "SSH",
		25:    "SMTP",
		53:    "DNS",
		69:    "TFTP",
		80:    "HTTP",
		443:   "HTTPS",
		1433:  "MSSQL",
		3306:  "MySQL",
		5432:  "PostgreSQL",
		1521:  "Oracle",
		3389:  "RDP",
		5900:  "VNC",
		8080:  "HTTP Proxy",
		8443:  "HTTPS Proxy",
		27017: "MongoDB",
		6379:  "Redis",
		9200:  "Elasticsearch",
	}
	if service, exists := services[port]; exists {
		return service
	}
	return "Unknown"
}
