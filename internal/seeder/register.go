package seeder

func init() {
	RegisterSeeder("File", NewFileSeederWithOptions)
	RegisterSeeder("Dict", NewDictSeederWithOptions)
}
