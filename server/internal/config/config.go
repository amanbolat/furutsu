package config

type Config struct {
	DbConnString string `split_words:"true" required:"true"`
	Port         int    `split_words:"true" required:"true"`
	MigratesDir  string `split_words:"true" required:"true"`
}
