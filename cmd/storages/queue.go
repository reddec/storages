package main

import (
	"bufio"
	"github.com/pkg/errors"
	"github.com/reddec/storages"
	"github.com/reddec/storages/queues"
	"io/ioutil"
	"log"
	"os"
)

const statusNoData = 127

type queueCmd struct {
	Put     queuePut     `command:"put" alias:"push" alias:"append" description:"put data to the queue"`
	Peek    queuePeek    `command:"peek" description:"get oldest data from queue but not remove"`
	Get     queueGet     `command:"get" alias:"pop" description:"get oldest data from queue and remove it"`
	Discard queueDiscard `command:"discard" description:"remove oldest data from queue (like silent get)"`
}

type queuePut struct {
	Line bool `short:"l" long:"line" env:"LINE" description:"Line mode for STDIN value - each line is new value"`
	Args struct {
		Values []string `description:"values to put to the queue, if not set - STDIN lines used" positional-arg-name:"values"`
	} `positional-args:"yes"`
}

func (q *queuePut) Execute(args []string) error {
	queue, db := config.getQueue()
	defer db.Close()

	for _, value := range q.Args.Values {
		err := queue.Put([]byte(value))
		if err != nil {
			return errors.Wrap(err, "put data to queue")
		}
	}
	if len(q.Args.Values) > 0 {
		return nil
	}

	if q.Line {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			err := queue.Put(scanner.Bytes())
			if err != nil {
				return errors.Wrap(err, "put data to queue")
			}
		}
		return nil
	}
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	err = queue.Put(data)
	if err != nil {
		return errors.Wrap(err, "put data to queue")
	}
	return nil
}

type queueGet struct {
}

func (q *queueGet) Execute(args []string) error {
	queue, db := config.getQueue()
	defer db.Close()
	data, err := queue.Get()
	if err == os.ErrNotExist {
		db.Close()
		os.Exit(statusNoData)
	} else if err != nil {
		return err
	}
	_, err = os.Stdout.Write(data)
	os.Stdout.Close()
	return err
}

type queuePeek struct {
}

func (q queuePeek) Execute(args []string) error {
	queue, db := config.getQueue()
	defer db.Close()
	data, err := queue.Peek()
	if err == os.ErrNotExist {
		db.Close()
		os.Exit(statusNoData)
	} else if err != nil {
		return err
	}
	_, err = os.Stdout.Write(data)
	return err
}

type queueDiscard struct {
}

func (q queueDiscard) Execute(args []string) error {
	queue, db := config.getQueue()
	defer db.Close()
	err := queue.Discard()
	if err == os.ErrNotExist {
		db.Close()
		os.Exit(statusNoData)
	}
	return err
}

func (cfg Config) getQueue() (storages.Queue, storages.Storage) {
	db := cfg.Storage()
	queue, err := queues.NaiveQueue(db)
	if err != nil {
		db.Close()
		log.Fatal("open queue:", err)
	}
	return queue, db
}
