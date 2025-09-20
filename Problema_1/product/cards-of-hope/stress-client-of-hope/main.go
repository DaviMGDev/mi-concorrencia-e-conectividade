package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type Dict map[string]any

type Request struct {
	Method string `json:"method"`
	Data   Dict   `json:"data,omitempty"`
}

type Response struct {
	Method string `json:"method"`
	Status string `json:"status"`
	Data   Dict   `json:"data,omitempty"`
}

type stats struct {
	sent     int64
	received int64
	errors   int64
	latSum   int64 // nanos
	latMax   int64 // nanos
}

func doPing(addr string, stop <-chan struct{}, s *stats, interval int) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		atomic.AddInt64(&s.errors, 1)
		return
	}
	defer conn.Close()
	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)
	for {
		select {
		case <-stop:
			return
		default:
		}
		req := Request{Method: "ping"}
		start := time.Now()
		atomic.AddInt64(&s.sent, 1)
		if err := enc.Encode(req); err != nil {
			atomic.AddInt64(&s.errors, 1)
			return
		}
		var resp Response
		if err := dec.Decode(&resp); err != nil || resp.Status != "ok" {
			atomic.AddInt64(&s.errors, 1)
		} else {
			atomic.AddInt64(&s.received, 1)
			elapsed := time.Since(start)
			atomic.AddInt64(&s.latSum, elapsed.Nanoseconds())
			for {
				old := atomic.LoadInt64(&s.latMax)
				if elapsed.Nanoseconds() > old {
					if atomic.CompareAndSwapInt64(&s.latMax, old, elapsed.Nanoseconds()) {
						break
					}
				} else {
					break
				}
			}
		}
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}

// Testa apenas a abertura de conexões simultâneas, sem enviar comandos.
type connResult struct {
	latency time.Duration
	err     error
}

func testConnections(addr string, clients int) {
	var wg sync.WaitGroup
	var openConns int64
	var failedConns int64
	results := make([]connResult, clients)
	stop := make(chan struct{})
	timeout := 2 * time.Second
	progressStep := clients / 100
	if progressStep == 0 {
		progressStep = 1
	}
	for i := 0; i < clients; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			start := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			d := &net.Dialer{}
			conn, err := d.DialContext(ctx, "tcp", addr)
			latency := time.Since(start)
			if err == nil {
				atomic.AddInt64(&openConns, 1)
				results[idx] = connResult{latency: latency, err: nil}
				<-stop // mantém conexão aberta até o final do teste
				conn.Close()
			} else {
				atomic.AddInt64(&failedConns, 1)
				results[idx] = connResult{latency: latency, err: err}
			}
			// Barra de progresso simples
			if (idx+1)%progressStep == 0 || idx == clients-1 {
				done := atomic.LoadInt64(&openConns) + atomic.LoadInt64(&failedConns)
				fmt.Printf("\rTentativas: %d/%d | Sucesso: %d | Falha: %d", done, clients, atomic.LoadInt64(&openConns), atomic.LoadInt64(&failedConns))
			}
		}(i)
	}
	// Aguarda um tempo para todas as conexões serem abertas
	time.Sleep(timeout + 1*time.Second)
	count := int(atomic.LoadInt64(&openConns))
	close(stop)
	wg.Wait()
	fmt.Println()

	// Estatísticas
	var latSum, latMax time.Duration
	for _, r := range results {
		if r.err == nil {
			latSum += r.latency
			if r.latency > latMax {
				latMax = r.latency
			}
		}
	}
	avgLat := float64(latSum.Microseconds()) / float64(count)
	fmt.Println("\n--- Resultados do Teste de Conexões ---")
	fmt.Printf("Conexões abertas com sucesso: %d/%d\n", count, clients)
	fmt.Printf("Falhas ao conectar: %d\n", failedConns)
	fmt.Printf("Latência média de conexão: %.2f ms\n", avgLat/1000)
	fmt.Printf("Latência máxima de conexão: %.2f ms\n", float64(latMax.Microseconds())/1000)
}

func main() {
	var (
		addr     string
		clients  int
		interval int
		duration int
		onlyConn bool
	)
	flag.StringVar(&addr, "addr", "localhost:8080", "Endereço do servidor (host:porta)")
	flag.IntVar(&clients, "clients", 100, "Número de conexões simultâneas")
	flag.IntVar(&interval, "interval", 100, "Intervalo entre pings (ms)")
	flag.IntVar(&duration, "duration", 10, "Duração do teste (segundos)")
	flag.BoolVar(&onlyConn, "onlyconn", false, "Testar apenas conexões simultâneas (sem enviar comandos)")
	flag.Parse()

	if onlyConn {
		fmt.Printf("Testando limite de conexões simultâneas: %d\n", clients)
		testConnections(addr, clients)
		return
	}

	fmt.Printf("Stress test: %d clientes, %d ms entre pings, duração %d s, servidor %s\n", clients, interval, duration, addr)

	var wg sync.WaitGroup
	var s stats
	stop := make(chan struct{})

	for i := 0; i < clients; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			doPing(addr, stop, &s, interval)
		}()
	}

	time.Sleep(time.Duration(duration) * time.Second)
	close(stop)
	wg.Wait()

	fmt.Println("\n--- Resultados do Stress Test ---")
	fmt.Printf("Total de pings enviados: %d\n", s.sent)
	fmt.Printf("Respostas recebidas: %d\n", s.received)
	fmt.Printf("Erros: %d\n", s.errors)
	avgLat := float64(s.latSum) / float64(s.received)
	fmt.Printf("Latência média: %.2f ms\n", avgLat/1e6)
	fmt.Printf("Latência máxima: %.2f ms\n", float64(s.latMax)/1e6)
	fmt.Printf("Pings/s: %.2f\n", float64(s.received)/float64(duration))
	if s.errors > 0 {
		os.Exit(1)
	}
}
