package config

import (
	"log"

	"github.com/spf13/viper"
)

var (
	IsProd              bool
	AppHost             string
	AppPort             int
	DBHost              string
	DBUser              string
	DBPassword          string
	DBName              string
	DBPort              int
	ProductTokenExpDays string
	LogMealBaseUrl      string
	LogMealApiKey       string
	JWTSecret           string
	JWTAccessExp        int
	JWTRefreshExp       int
	JWTResetPasswordExp int
	JWTVerifyEmailExp   int
	SMTPHost            string
	SMTPPort            int
	SMTPUsername        string
	SMTPPassword        string
	EmailFrom           string
	GoogleClientID      string
	GoogleClientSecret  string
	RedirectURL         string
	MidtransServerKey   string
	MidtransStatus      string
)

func init() {
	loadConfig()

	// server configuration
	IsProd = viper.GetString("APP_ENV") == "prod"
	AppHost = viper.GetString("APP_HOST")
	AppPort = viper.GetInt("APP_PORT")

	// database configuration
	DBHost = viper.GetString("DB_HOST")
	DBUser = viper.GetString("DB_USER")
	DBPassword = viper.GetString("DB_PASSWORD")
	DBName = viper.GetString("DB_NAME")
	DBPort = viper.GetInt("DB_PORT")

	// product token
	ProductTokenExpDays = viper.GetString("PRODUCT_TOKEN_EXP_DAYS")

	// log meal
	LogMealBaseUrl = viper.GetString("LOG_MEAL_BASE_URL")
	LogMealApiKey = viper.GetString("LOG_MEAL_API_KEY")

	// jwt configuration
	JWTSecret = viper.GetString("JWT_SECRET")
	JWTAccessExp = viper.GetInt("JWT_ACCESS_EXP_MINUTES")
	JWTRefreshExp = viper.GetInt("JWT_REFRESH_EXP_DAYS")
	JWTResetPasswordExp = viper.GetInt("JWT_RESET_PASSWORD_EXP_MINUTES")
	JWTVerifyEmailExp = viper.GetInt("JWT_VERIFY_EMAIL_EXP_MINUTES")

	// SMTP configuration
	SMTPHost = viper.GetString("SMTP_HOST")
	SMTPPort = viper.GetInt("SMTP_PORT")
	SMTPUsername = viper.GetString("SMTP_USERNAME")
	SMTPPassword = viper.GetString("SMTP_PASSWORD")
	EmailFrom = viper.GetString("EMAIL_FROM")

	// oauth2 configuration
	GoogleClientID = viper.GetString("GOOGLE_CLIENT_ID")
	GoogleClientSecret = viper.GetString("GOOGLE_CLIENT_SECRET")
	RedirectURL = viper.GetString("REDIRECT_URL")

	// Midtrans configuration
	MidtransServerKey = viper.GetString("MIDTRANS_SERVER_KEY")
	MidtransStatus = viper.GetString("MIDTRANS_STATUS")
}

func loadConfig() {
	configPaths := []string{
		"./",     // For app
		"../../", // For test folder
	}

	for _, path := range configPaths {
		viper.SetConfigFile(path + ".env")

		if err := viper.ReadInConfig(); err == nil {
			log.Printf("Config file loaded from %s", path)
			return
		}
	}

	log.Println("Failed to load any config file")
}
