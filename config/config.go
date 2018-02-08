package config

type (
	MongodbConfig struct {
		DBName   string `json:"db_name"`
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
	}

	MysqlConfig struct {
		DBName   string `json:"db_name"`
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json:"port"`
	}

	Settings struct {
		Bind    string        `json:"bind"`
		Domain  string        `json:"domain"`
		BaseURL string        `json:"-"`
		Mongodb MongodbConfig `json:"mongodb"`
		Mysql   MysqlConfig   `json:"mysql"`
	}
)

var (
	local = Settings{
		Bind:   "0.0.0.0:5000",
		Domain: "",
		Mongodb: MongodbConfig{
			DBName: "test",
			Host:   "127.0.0.1:27017",
		},
	}
)

func Load(name string) *Settings {
	// todo
	return &local
}

func Get() *Settings {
	return &local
}
