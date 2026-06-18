package spidomtr_test

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/spider-pigs/spidomtr"
	"github.com/spider-pigs/spidomtr/pkg/handlers"
	"github.com/spider-pigs/spidomtr/pkg/testunit"
)

func TestRunner(t *testing.T) {
	t.Run("test runner with passing tests", func(t *testing.T) {
		runner := spidomtr.NewRunner(
			spidomtr.Iterations(50),
			spidomtr.ShowLogo(false),
			spidomtr.ShowSummary(false),
		)

		test := testunit.New(
			testunit.Test(func(context.Context, []interface{}) ([]interface{}, error) {
				return nil, nil
			}),
		)

		res := runner.Run(context.Background(), test)
		if res.Stats.Count != 50 {
			t.Fatalf("expected Count 50, got %d", res.Stats.Count)
		}
		if res.Stats.Passed != 50 {
			t.Fatalf("expected Passed 50, got %d", res.Stats.Passed)
		}
		if res.Stats.Skips != 0 {
			t.Fatalf("expected Skips 0, got %d", res.Stats.Skips)
		}
		if res.Stats.Errors != 0 {
			t.Fatalf("expected Errors 0, got %d", res.Stats.Errors)
		}
	})
	t.Run("test runner with error tests", func(t *testing.T) {
		runner := spidomtr.NewRunner(
			spidomtr.Iterations(50),
			spidomtr.ShowLogo(false),
			spidomtr.ShowSummary(false),
		)

		test := testunit.New(
			testunit.Test(func(context.Context, []interface{}) ([]interface{}, error) {
				return nil, errors.New("whooops")
			}),
		)

		res := runner.Run(context.Background(), test)
		if res.Stats.Count != 50 {
			t.Fatalf("expected Count 50, got %d", res.Stats.Count)
		}
		if res.Stats.Passed != 0 {
			t.Fatalf("expected Passed 0, got %d", res.Stats.Passed)
		}
		if res.Stats.Skips != 0 {
			t.Fatalf("expected Skips 0, got %d", res.Stats.Skips)
		}
		if res.Stats.Errors != 50 {
			t.Fatalf("expected Errors 50, got %d", res.Stats.Errors)
		}
	})
	t.Run("test runner with skipped tests", func(t *testing.T) {
		runner := spidomtr.NewRunner(
			spidomtr.Iterations(50),
			spidomtr.ShowLogo(false),
			spidomtr.ShowSummary(false),
		)

		test := testunit.New(
			testunit.Enabled(func() (bool, string) {
				return false, "leave me alone"
			}),
		)

		res := runner.Run(context.Background(), test)
		if res.Stats.Count != 50 {
			t.Fatalf("expected Count 50, got %d", res.Stats.Count)
		}
		if res.Stats.Passed != 0 {
			t.Fatalf("expected Passed 0, got %d", res.Stats.Passed)
		}
		if res.Stats.Skips != 50 {
			t.Fatalf("expected Skips 50, got %d", res.Stats.Skips)
		}
		if res.Stats.Errors != 0 {
			t.Fatalf("expected Errors 0, got %d", res.Stats.Errors)
		}
	})
	t.Run("test runner with pass, skip and error", func(t *testing.T) {
		runner := spidomtr.NewRunner(
			spidomtr.Iterations(50),
			spidomtr.ShowLogo(false),
			spidomtr.ShowSummary(false),
		)

		test1 := testunit.New(
			testunit.Test(func(context.Context, []interface{}) ([]interface{}, error) {
				time.Sleep(1 * time.Millisecond)
				return nil, nil
			}),
		)

		test2 := testunit.New(
			testunit.Test(func(context.Context, []interface{}) ([]interface{}, error) {
				time.Sleep(10 * time.Millisecond)
				return nil, errors.New("whooops")
			}),
		)

		test3 := testunit.New(
			testunit.Enabled(func() (bool, string) {
				return false, "leave me alone"
			}),
		)

		res := runner.Run(context.Background(), test1, test2, test3)
		if res.Stats.Count != 150 {
			t.Fatalf("expected Count 150, got %d", res.Stats.Count)
		}
		if res.Stats.Passed != 50 {
			t.Fatalf("expected Passed 50, got %d", res.Stats.Passed)
		}
		if res.Stats.Skips != 50 {
			t.Fatalf("expected Skips 50, got %d", res.Stats.Skips)
		}
		if res.Stats.Errors != 50 {
			t.Fatalf("expected Errors 50, got %d", res.Stats.Errors)
		}
	})
}

func TestJoinResults(t *testing.T) {
	test := testunit.New(
		testunit.Test(func(context.Context, []interface{}) ([]interface{}, error) {
			return nil, nil
		}),
	)

	runner1 := spidomtr.NewRunner(
		spidomtr.Iterations(50),
		spidomtr.ShowLogo(false),
		spidomtr.ShowSummary(false),
	)
	res1 := runner1.Run(context.Background(), test)

	runner2 := spidomtr.NewRunner(
		spidomtr.Iterations(50),
		spidomtr.ShowLogo(false),
		spidomtr.ShowSummary(false),
	)
	res2 := runner2.Run(context.Background(), test)

	res := spidomtr.JoinResults(spidomtr.DefaultHistogramBuckets, spidomtr.DefaultPercentiles, res1, res2)
	if res.Stats.Count != 100 {
		t.Fatalf("expected Count 100, got %d", res.Stats.Count)
	}
	if res.Stats.Passed != 100 {
		t.Fatalf("expected Passed 100, got %d", res.Stats.Passed)
	}
	if len(res.Stats.Errorm) != 0 {
		t.Fatalf("expected 0 error messages, got %d", len(res.Stats.Errorm))
	}
	if res.Stats.Skips != 0 {
		t.Fatalf("expected Skips 0, got %d", res.Stats.Skips)
	}
	if res.Stats.Errors != 0 {
		t.Fatalf("expected Errors 0, got %d", res.Stats.Errors)
	}
}

func TestRunnerWithUsers(t *testing.T) {
	runner := spidomtr.NewRunner(
		spidomtr.ID("stupid"),
		spidomtr.Description("just running some stupid tests"),
		spidomtr.ShowLogo(false),
		spidomtr.ShowSummary(false),
		spidomtr.Iterations(50),
		spidomtr.Timeout(time.Second*10),
		spidomtr.Users(10),
	)

	test1 := testunit.New(
		testunit.ID("awesome_test"),
		testunit.Test(func(context.Context, []interface{}) ([]interface{}, error) {
			time.Sleep(1 * time.Millisecond)
			return nil, nil
		}),
	)

	test2 := testunit.New(
		testunit.ID("not_always_an_awesome_test"),
		testunit.Test(func(context.Context, []interface{}) ([]interface{}, error) {
			time.Sleep(1 * time.Millisecond)
			if coinflip() == "heads" {
				return nil, errors.New("whooops")
			}
			return nil, nil
		}),
	)

	test3 := testunit.New(
		testunit.ID("skipped_test"),
		testunit.Enabled(func() (bool, string) {
			return false, "leave me alone"
		}),
	)

	res := runner.Run(context.Background(), test1, test2, test3)
	if res.Stats.Count != 1500 {
		t.Fatalf("expected Count 1500, got %d", res.Stats.Count)
	}
}

func TestHandlers(t *testing.T) {
	runner := spidomtr.NewRunner(
		spidomtr.HistogramBuckets(5),
		spidomtr.Description("just running some stupid tests"),
		spidomtr.Handlers(handlers.ProgressBar()),
		spidomtr.ID("stupid"),
		spidomtr.Iterations(50),
		spidomtr.ShowSummary(true),
		spidomtr.Timeout(time.Second*10),
		spidomtr.Users(10),
	)

	min := 10
	max := 50

	test1 := testunit.New(
		testunit.ID("awesome_test"),
		testunit.Test(func(context.Context, []interface{}) ([]interface{}, error) {
			latency := rand.Intn(max-min) + min
			time.Sleep(time.Duration(latency) * time.Millisecond)
			return nil, nil
		}),
	)

	test2 := testunit.New(
		testunit.ID("not_always_an_awesome_test"),
		testunit.Test(func(context.Context, []interface{}) ([]interface{}, error) {
			time.Sleep(20 * time.Millisecond)
			if coinflip() == "heads" {
				return nil, errors.New("whooops")
			}
			return nil, nil
		}),
	)

	test3 := testunit.New(
		testunit.ID("skipped_test"),
		testunit.Enabled(func() (bool, string) {
			return false, "leave me alone"
		}),
	)

	res := runner.Run(context.Background(), test1, test2, test3)
	if res.Stats.Count != 1500 {
		t.Fatalf("expected Count 1500, got %d", res.Stats.Count)
	}
	if res.Stats.Skips != 500 {
		t.Fatalf("expected Skips 500, got %d", res.Stats.Skips)
	}
	if res.Stats.Passed <= 500 {
		t.Fatalf("expected Passed > 500, got %d", res.Stats.Passed)
	}
}

func coinflip() string {
	coin := []string{
		"heads",
		"tails",
	}
	rand.Seed(time.Now().UnixNano())
	return coin[rand.Intn(len(coin))]
}
