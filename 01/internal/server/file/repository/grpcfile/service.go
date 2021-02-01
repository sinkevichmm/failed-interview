package grpcservice

import (
	"errors"
	"failed-interview/01/internal/models"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
)

type GRPCService struct {
	fileClient []*FileClient
	Mutex      *sync.Mutex
	capOrder   []freeCapacity
}

type freeCapacity struct {
	managerNum   int
	freeCapacity uint
}

func NewGRPCService(grpcaddress string, auth string) *GRPCService {
	addrs := strings.Split(grpcaddress, " ")

	fileClients := []*FileClient{}

	for _, a := range addrs {
		fc := NewFileClient(a, auth)
		if fc == nil {
			fmt.Printf("cannot dial server %s\n", a)
			continue
		}

		fileClients = append(fileClients, fc)
	}

	if len(fileClients) == 0 {
		log.Fatal("no available managers")
	}

	cap := make([]freeCapacity, len(fileClients))
	for i := 0; i < len(cap); i++ {
		cap[i].managerNum = i
	}

	return &GRPCService{fileClient: fileClients, Mutex: &sync.Mutex{}, capOrder: cap}
}

func (g *GRPCService) getFreeCapacity(index int, wg *sync.WaitGroup) {
	defer wg.Done()

	free := g.fileClient[index].GetFreeCapacity()

	g.Mutex.Lock()
	defer g.Mutex.Unlock()

	g.capOrder[index].freeCapacity = free
	g.capOrder[index].managerNum = index
}

func (g *GRPCService) sortByFreeCapacity() {
	var wg sync.WaitGroup

	for i := 0; i < len(g.fileClient); i++ {
		wg.Add(1)

		go g.getFreeCapacity(i, &wg)
	}

	wg.Wait()

	g.Mutex.Lock()
	defer g.Mutex.Unlock()

	sort.Slice(g.capOrder, func(i, j int) bool {
		return g.capOrder[i].freeCapacity > g.capOrder[j].freeCapacity
	})
}

func (g *GRPCService) SaveFile(file *models.File) (id string, err error) {
	g.sortByFreeCapacity()

	if g.capOrder[0].freeCapacity == 0 {
		return id, errors.New("storage is full")
	}

	id, err = g.fileClient[g.capOrder[0].managerNum].SaveFile(file)
	if err != nil {
		return id, err
	}

	return id, err
}

func (g *GRPCService) GetFileIDs() (ids []string) {
	ch := make(chan []string, len(g.fileClient))

	var wg sync.WaitGroup

	for i := 0; i < len(g.fileClient); i++ {
		wg.Add(1)

		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()

			ch <- g.fileClient[i].GetFileIDs()
		}(i, &wg)
	}

	wg.Wait()
	close(ch)

	for id := range ch {
		ids = append(ids, id...)
	}

	return ids
}

func (g *GRPCService) GetFileInfoByID(id string) (meta *models.Meta, err error) {
	ch := make(chan *models.Meta, len(g.fileClient))

	var wg sync.WaitGroup

	for i := 0; i < len(g.fileClient); i++ {
		wg.Add(1)

		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()

			m, err := g.fileClient[i].GetFileInfoByID(id)

			if err == nil {
				ch <- m
				return
			}
			ch <- nil
		}(i, &wg)
	}

	wg.Wait()
	close(ch)

	for m := range ch {
		if m != nil {
			meta = m
		}
	}

	if meta == nil {
		return nil, fmt.Errorf(" %s %w", id, ErrIDNotFound)
	}

	return meta, err
}

func (g *GRPCService) DeleteFile(id string) (err error) {
	ch := make(chan error, len(g.fileClient))

	var wg sync.WaitGroup

	for i := 0; i < len(g.fileClient); i++ {
		wg.Add(1)

		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()
			ch <- g.fileClient[i].DeleteFile(id)
		}(i, &wg)
	}

	wg.Wait()
	close(ch)

	deleted := false

	for e := range ch {
		if e == nil {
			deleted = true
		}
	}

	if deleted {
		return nil
	}

	return fmt.Errorf(" %s %w", id, ErrIDNotFound)
}

func (g *GRPCService) GetFileByID(id string) (file *models.File, err error) {
	ch := make(chan *models.File, len(g.fileClient))

	var wg sync.WaitGroup

	for i := 0; i < len(g.fileClient); i++ {
		wg.Add(1)

		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()

			f, err := g.fileClient[i].GetFileByID(id)
			if err == nil {
				ch <- f
				return
			}
			ch <- nil
		}(i, &wg)
	}

	wg.Wait()
	close(ch)

	for f := range ch {
		if f != nil {
			file = f
		}
	}

	if file == nil {
		return nil, fmt.Errorf(" %s %w", id, ErrIDNotFound)
	}

	return file, err
}
