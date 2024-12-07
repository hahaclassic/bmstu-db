package controller

import (
	"context"
	"fmt"

	"github.com/hahaclassic/databases/09_redis/internal/service"
	tableoutput "github.com/hahaclassic/databases/09_redis/pkg/table"
	"github.com/jedib0t/go-pretty/table"
	"golang.org/x/exp/slog"
)

type optionHandler struct {
	name string
	f    func(ctx context.Context) error
}

type Controller struct {
	s        *service.Service
	handlers []*optionHandler
}

func New(s *service.Service) *Controller {
	optionHandlers := []*optionHandler{
		{
			name: "Top 10 most streamed tracks",
			f:    s.Top10MostStreamedTracks(service.DB, true),
		},
		{
			name: "Top 10 most streamed tracks in cache",
			f:    s.Top10MostStreamedTracks(service.Cache, true),
		},
		{
			name: "Top 10 most streamed tracks in DB with cache",
			f:    s.Top10MostStreamedTracks(service.CachedDB, true),
		},
		{
			name: "Bench without changes",
			f:    s.Bench(nil, "basic"),
		},
		{
			name: "Bench with insertions",
			f:    s.Bench(s.AddTrack, "insert"),
		},
		{
			name: "Bench with deletions",
			f:    s.Bench(s.DeleteRandomTrack, "delete"),
		},
		{
			name: "Bench with updates",
			f:    s.Bench(s.UpdateRandomTrackStreams, "update"),
		},
	}

	return &Controller{
		s:        s,
		handlers: optionHandlers,
	}
}

func (c *Controller) Run(ctx context.Context) {
	const exitNumber = 0

	for {
		option := c.getOperationNumber()

		if option == exitNumber {
			break
		}

		if err := c.handlers[option-1].f(ctx); err != nil {
			slog.Error("handler", "err", err)
			continue
		}

		fmt.Println()
	}
}

func (c Controller) getOperationNumber() int {
	c.showMenu()
	fmt.Println()

	var option int
	for {
		fmt.Print("Enter operation number: ")
		if _, err := fmt.Scan(&option); err == nil &&
			option >= 0 && option <= len(c.handlers) {
			break
		}

		fmt.Println("Invalid number. Try again.")
	}

	return option
}

func (c Controller) showMenu() {
	headers := []string{"#", "Operation"}
	rows := make([][]any, 0, len(c.handlers))

	for i := range c.handlers {
		rows = append(rows, []any{i + 1, c.handlers[i].name})
	}
	rows = append(rows, []any{0, "Exit"})

	tableoutput.PrintTable(table.StyleDefault, headers, rows)
}
