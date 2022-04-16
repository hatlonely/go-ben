package seeder

import (
	"bufio"
	"io"
	"math/rand"
	"os"

	jsoniter "github.com/json-iterator/go"

	"github.com/pkg/errors"
)

type FileSeederOptions struct {
	Name             string
	IgnoreParseError bool
}

func NewFileSeederWithOptions(options *FileSeederOptions) (*FileSeeder, error) {
	fp, err := os.Open(options.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "os.Open [%s] failed", options.Name)
	}

	var seeds []interface{}
	reader := bufio.NewReader(fp)
	for {
		buf, err := reader.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				return nil, errors.Wrapf(err, "reader.ReadString failed")
			}
			break
		}
		if len(buf) == 1 {
			continue
		}
		var v interface{}
		if err := jsoniter.Unmarshal(buf, &v); err != nil {
			if options.IgnoreParseError {
				continue
			}
			return nil, errors.Wrapf(err, "parse [%s] failed. file [%s]", string(buf), options.Name)
		}

		seeds = append(seeds, v)
	}

	return &FileSeeder{
		seeds: seeds,
	}, nil
}

type FileSeeder struct {
	seeds []interface{}
}

func (s *FileSeeder) Seed() interface{} {
	return s.seeds[rand.Intn(len(s.seeds))]
}
