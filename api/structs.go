package api

type MazeConfig struct {
	Player             int `yaml:"player"`
	Wall               int `yaml:"wall"`
	Space              int `yaml:"space"`
	Maks_element_count int `yaml:"maks_element_count"`
}

type Configuration struct {
	Maze                    MazeConfig `yaml:"maze"`
	Port                    string     `yaml:"port"`
	Mongo_connection_string string     `yaml:"mongo_connection_string"`
	LogLevel                string     `yaml:"log_level"`
	Database                string     `yaml:"database"`
}
