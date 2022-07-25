package firehose

import (
	"context"
	"sync"
	"time"
)

// Connection creates a single connection to the firehose stream.
type Connection struct {
	Since   time.Time
	Stream  func(ctx context.Context, since time.Time, cb func(evt *Event)) error
	Handler func(evt *Event)
}

// NewConnectionManager creates a connection manager with an empty connections list.
func NewConnectionManger() *ConnectionManger {
	return &ConnectionManger{
		conns: []*Connection{},
	}
}

// ConnectionManager helps keep connections open for firehose (realtime),
// satisfies "alt least once" delivery for the event.
type ConnectionManger struct {
	conns []*Connection
}

// Add appends new connection to the list of connections.
func (cm *ConnectionManger) Add(con *Connection) {
	cm.conns = append(cm.conns, con)
}

// Connect opens connections for the list in a blocking call. Leaves concurrency
// to the caller. Pushes connection errors to the provided channel. If you don't want
// error deliveries just pass nil instead of the channel.
func (cm *ConnectionManger) Connect(ctx context.Context, errs chan error) {
	wg := new(sync.WaitGroup)
	wg.Add(len(cm.conns))

	for _, con := range cm.conns {
		go func(con *Connection) {
			defer wg.Done()

			for {
				err := con.Stream(ctx, con.Since, func(evt *Event) {
					if len(evt.ID) > 0 {
						con.Since = evt.ID[0].Dt
					}

					con.Handler(evt)
				})

				if err != nil && errs != nil {
					errs <- err
				}

				if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
					break
				}
			}
		}(con)
	}

	wg.Wait()

	if errs != nil {
		close(errs)
	}
}
