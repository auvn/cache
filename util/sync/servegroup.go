package sync

import (
	"sync"
	"sync/atomic"
)

const (
	unlockedShutdown uint32 = iota
	lockedShutdown
)

type Quit <-chan struct{}

type Server interface {
	Serve(Quit) error
}

func serve(wg *sync.WaitGroup, errCh chan error, quit Quit, s Server) {
	wg.Add(1)
	defer wg.Done()
	err := s.Serve(quit)
	if err != nil {
		select {
		case <-quit:
			return
		case errCh <- err:
		default:
		}
	}
}

type ServeGroup interface {
	Serve(...Server) ServeGroup
	Shutdown()
	Quit() Quit
	Err() <-chan error
}

type serveGroup struct {
	shutdownLock uint32
	wg           *sync.WaitGroup

	mu   sync.Mutex
	quit chan struct{}
	err  chan error
}

type ServeFn func() error

func (self *serveGroup) refresh() {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.err = make(chan error, 1)
	self.quit = make(chan struct{}, 1)
}

func (self *serveGroup) Serve(ss ...Server) ServeGroup {
	for _, s := range ss {
		self.mu.Lock()
		go serve(self.wg, self.err, self.quit, s)
		self.mu.Unlock()
	}
	return self
}

func (self *serveGroup) Quit() Quit {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.quit
}

func (self *serveGroup) Err() <-chan error {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.err
}

func (self *serveGroup) Shutdown() {
	if !atomic.CompareAndSwapUint32(
		&self.shutdownLock,
		unlockedShutdown,
		lockedShutdown,
	) {
		return
	}

	close(self.quit)
	self.wg.Wait()

	close(self.err)
	self.refresh()
	atomic.StoreUint32(&self.shutdownLock, unlockedShutdown)
}

func NewServeGroup() ServeGroup {
	sg := &serveGroup{
		wg:           &sync.WaitGroup{},
		shutdownLock: unlockedShutdown,
	}
	sg.refresh()
	return sg
}

func Serve(ss ...Server) ServeGroup {
	group := NewServeGroup()
	group.Serve(ss...)
	return group
}
