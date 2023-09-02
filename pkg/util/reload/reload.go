package reload

import (
	"bufio"
	"context"
	"fmt"
	"github.com/bmatcuk/doublestar"
	"github.com/pkg/errors"
	"io"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	zlog "github.com/rs/zerolog/log"
	"golang.org/x/xerrors"
)

const (
	waitForTerm = 5 * time.Second
)

type Config struct {
	Cmd      []string
	Targets  []string
	Patterns []string
	Ignores  []string
	Delay    time.Duration
	Restart  bool
	SigOpt   string
	Verbose  bool
}

func Reload(config Config) error {
	if len(config.Cmd) == 0 {
		return errors.New("command is required")
	}
	if len(config.Targets) == 0 {
		config.Targets = []string{"./"}
	}
	if len(config.Patterns) == 0 {
		config.Patterns = []string{"**"}
	}
	sig, sigstr := parseSignalOption(config.SigOpt)
	if sig == nil {
		return errors.Errorf("invalid signal: %q", sigstr)
	}

	zlog.Debug().Interface("targets", config.Targets).Interface("patterns", config.Patterns).Interface("ignores", config.Ignores).Msg("Reload")
	modC, errC, err := watcher(config.Targets, config.Patterns, config.Ignores)
	if err != nil {
		return errors.Wrapf(err, "watcher")
	}

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	reload := runner(ctx, &wg, config.Cmd, config.Delay, sig.(syscall.Signal), config.Restart)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case name, ok := <-modC:
				if !ok {
					cancel()
					wg.Wait()
					log.Fatalf("[RELOAD] wacher closed")
					return
				}
				logVerbose("modified: %q", name)
				reload <- name
			case err := <-errC:
				cancel()
				wg.Wait()
				log.Fatalf("[RELOAD] wacher error: %v", err)
				return
			}
		}
	}()

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	sig = <-s
	log.Printf("[RELOAD] signal: %v", sig)
	cancel()
	wg.Wait()
	return nil
}

func logVerbose(fmt string, args ...interface{}) {
	log.Printf("[RELOAD] "+fmt, args...)
}

func watcher(targets, patterns, ignores []string) (<-chan string, <-chan error, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, nil, err
	}

	if err := addTargets(w, targets, patterns, ignores); err != nil {
		return nil, nil, err
	}

	modC := make(chan string)
	errC := make(chan error)

	go func() {
		defer close(modC)
		for {
			select {
			case event, ok := <-w.Events:
				if !ok {
					errC <- xerrors.Errorf("watcher.Events closed")
					return
				}

				name := filepath.ToSlash(event.Name)

				if ignore, err := matchPatterns(name, ignores); err != nil {
					errC <- xerrors.Errorf("match ignores: %w", err)
					return
				} else if ignore {
					continue
				}

				if match, err := matchPatterns(name, patterns); err != nil {
					errC <- xerrors.Errorf("match patterns: %w", err)
					return
				} else if match {
					modC <- name
				}

				// add watcher if new directory.
				if event.Has(fsnotify.Create) {
					fi, err := os.Stat(name)
					if err != nil {
						// ignore stat errors (notfound, permission, etc.)
						log.Printf("[RELOAD] watcher: %v", err)
					} else if fi.IsDir() {
						err := addDirRecursive(w, fi, name, patterns, ignores, modC)
						if err != nil {
							errC <- err
							return
						}
					}
				}

			case err, ok := <-w.Errors:
				errC <- xerrors.Errorf("watcher.Errors (%v): %w", ok, err)
				return
			}
		}
	}()

	return modC, errC, nil
}

func matchPatterns(t string, pats []string) (bool, error) {
	for _, p := range pats {
		m, err := doublestar.Match(p, t)
		if err != nil {
			return false, xerrors.Errorf("match(%v, %v): %w", p, t, err)
		}
		if m {
			return true, nil
		}
		if strings.HasPrefix(t, "./") {
			m, err = doublestar.Match(p, t[2:])
			if err != nil {
				return false, xerrors.Errorf("match(%v, %v): %w", p, t[2:], err)
			}
			if m {
				return true, nil
			}
		}
	}
	return false, nil
}

func addTargets(w *fsnotify.Watcher, targets, patterns, ignores []string) error {
	for _, t := range targets {
		t = path.Clean(t)
		fi, err := os.Stat(t)
		if err != nil {
			return xerrors.Errorf("stat: %w", err)
		}
		if fi.IsDir() {
			if err := addDirRecursive(w, fi, t, patterns, ignores, nil); err != nil {
				return err
			}
		}
		logVerbose("watching target: %q", t)
		if err := w.Add(t); err != nil {
			return err
		}
	}
	return nil
}

func addDirRecursive(w *fsnotify.Watcher, fi fs.FileInfo, t string, patterns, ignores []string, ch chan<- string) error {
	logVerbose("watching target: %q", t)
	err := w.Add(t)
	if err != nil {
		return xerrors.Errorf("wacher add: %w", err)
	}
	des, err := os.ReadDir(t)
	if err != nil {
		return xerrors.Errorf("read dir: %w", err)
	}
	for _, de := range des {
		name := path.Join(t, de.Name())
		if ignore, err := matchPatterns(name, ignores); err != nil {
			return xerrors.Errorf("match ignores: %w", err)
		} else if ignore {
			continue
		}
		if ch != nil {
			if match, err := matchPatterns(name, patterns); err != nil {
				return xerrors.Errorf("match patterns: %w", err)
			} else if match {
				ch <- name
			}
		}
		if de.IsDir() {
			fi, err := de.Info()
			if err != nil {
				return err
			}
			err = addDirRecursive(w, fi, name, patterns, ignores, ch)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type bytesErr struct {
	bytes []byte
	err   error
}

// stdinReader bypasses stdin to child processes
//
// cmd.Wait() blocks until stdin.Read() returns.
// so stdinReader.Read() returns EOF when the child process exited.
type stdinReader struct {
	input    <-chan bytesErr
	chldDone <-chan struct{}
}

func (s *stdinReader) Read(b []byte) (int, error) {
	select {
	case be, ok := <-s.input:
		if !ok {
			return 0, io.EOF
		}
		return copy(b, be.bytes), be.err
	case <-s.chldDone:
		return 0, io.EOF
	}
}

func clearChBuf[T any](c <-chan T) {
	for {
		select {
		case <-c:
		default:
			return
		}
	}
}

func runner(ctx context.Context, wg *sync.WaitGroup, cmd []string, delay time.Duration, sig syscall.Signal, autorestart bool) chan<- string {
	reload := make(chan string)
	trigger := make(chan string)

	go func() {
		for name := range reload {
			// ignore restart when the trigger is not waiting
			select {
			case trigger <- name:
			default:
			}
		}
	}()

	var pcmd string // command string for display.
	for _, s := range cmd {
		if strings.ContainsAny(s, " \t\"'") {
			s = fmt.Sprintf("%q", s)
		}
		pcmd += " " + s
	}
	pcmd = pcmd[1:]

	stdinC := make(chan bytesErr, 1)
	go func() {
		b1 := make([]byte, 255)
		b2 := make([]byte, 255)
		for {
			n, err := os.Stdin.Read(b1)
			stdinC <- bytesErr{b1[:n], err}
			b1, b2 = b2, b1
		}
	}()

	chldDone := makeChildDoneChan()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			cmdctx, cancel := context.WithCancel(ctx)
			restart := make(chan struct{})
			done := make(chan struct{})

			go func() {
				log.Printf("[RELOAD] start: %s", pcmd)
				clearChBuf(chldDone)
				stdin := &stdinReader{stdinC, chldDone}
				err := runCmd(cmdctx, cmd, sig, stdin)
				if err != nil {
					log.Printf("[RELOAD] command error: %v", err)
				} else {
					log.Printf("[RELOAD] command exit status 0")
				}
				if autorestart {
					close(restart)
				}

				close(done)
			}()

			select {
			case <-ctx.Done():
				cancel()
				<-done
				return
			case name := <-trigger:
				log.Printf("[RELOAD] triggered: %q", name)
			case <-restart:
				logVerbose("auto restart")
			}

			logVerbose("wait %v", delay)
			t := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				cancel()
				<-done
				return
			case <-t.C:
			}
			cancel()
			<-done // wait process closed
		}
	}()

	return reload
}

func runCmd(ctx context.Context, cmd []string, sig syscall.Signal, stdin *stdinReader) error {
	c := prepareCommand(cmd)
	c.Stdin = bufio.NewReader(stdin)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Start(); err != nil {
		return err
	}

	var cerr error
	done := make(chan struct{})
	go func() {
		cerr = waitCmd(c)
		close(done)
	}()

	select {
	case <-done:
		if cerr != nil {
			cerr = xerrors.Errorf("process exit: %w", cerr)
		}
		return cerr
	case <-ctx.Done():
		if err := killChilds(c, sig); err != nil {
			return xerrors.Errorf("kill childs: %w", err)
		}
	}

	select {
	case <-done:
	case <-time.NewTimer(waitForTerm).C:
		if err := killChilds(c, syscall.SIGKILL); err != nil {
			return xerrors.Errorf("kill childs (SIGKILL): %w", err)
		}
		<-done
	}

	if cerr != nil {
		return xerrors.Errorf("process canceled: %w", cerr)
	}
	return nil
}
