// configs/config.go
package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

// Configという名前の構造体(struct)を作るのだ。
//    .envファイルの中身と対応している
type Config struct {
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSLMODE"`
	ServerPort string `mapstructure:"SERVER_PORT"`
}

// DSNはデータベース接続文字列(DSN)を生成するメソッドなのだ
func (c *Config) DSN() string {
	// "postgres://user:password@host:port/dbname?sslmode=disable" の形に組み立てるのだ
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode)
}

// 5. LoadConfigは設定をファイルや環境変数から読み込む関数なのだ
func LoadConfig() (*Config, error) {
	// 6. Viperに設定ファイルの場所と名前を教えるのだ
	viper.AddConfigPath(".")        // 実行ファイルのある場所を探す
	viper.SetConfigName(".env")     // ".env"という名前のファイルを探す
	viper.SetConfigType("env")      // ファイルの種類は"env"形式だと教える

	// 7. 環境変数からも読み込めるように設定するのだ
	viper.AutomaticEnv()

	// 8. 設定ファイルの読み込みを実行するのだ
	if err := viper.ReadInConfig(); err != nil {
		// ファイルがなくても環境変数があればOKなので、ここではエラーを無視してもいい場合がある
		// でも、今回はファイルがないとエラーにするのだ
		return nil, fmt.Errorf("設定ファイルの読み込みに失敗しました: %w", err)
	}

	// 9. 読み込んだ設定を、Config構造体に流し込む（アンマーシャルする）のだ
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("設定の解析に失敗しました: %w", err)
	}

	// 10. 全て成功したら、設定が詰まったconfigオブジェクトを返すのだ
	return &config, nil
}
