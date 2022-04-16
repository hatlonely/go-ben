package seeder

import "math/rand"

type DictSeederOptions []map[string]interface{}

func NewDictSeederWithOptions(options *DictSeederOptions) (*DictSeeder, error) {
	return &DictSeeder{
		seeds: *options,
	}, nil
}

type DictSeeder struct {
	seeds []map[string]interface{}
}

func (s *DictSeeder) Seed() interface{} {
	return s.seeds[rand.Intn(len(s.seeds))]
}
