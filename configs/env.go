package configs

import "time"

const (
	Development = "development"
	Staging     = "staging"
	Production  = "production"
)

type Environment struct {
	Env       string `env:"ENVIRONMENT,default=development"`
	Server    Server
	MongoDB   MongoDB
	Security  Security
	Cors      Cors
	RateLimit RateLimit
	SMTP      SMTP
	URLs      URLs
	Key       Key
	OTP       OTP
}

type Server struct {
	Port int    `env:"SERVER_PORT,default=8080"`
	Host string `env:"SERVER_HOST,default=localhost"`
}

type MongoDB struct {
	ConnectionURI     string        `env:"MONGODB_CONNECTION_URI,required"`
	DatabaseName      string        `env:"MONGODB_DATABASE_NAME,required"`
	MaxPoolSize       uint64        `env:"MONGODB_MAX_POOL_SIZE,default=100"`
	MaxConnIdleTime   time.Duration `env:"MONGODB_MAX_CONN_IDLE_TIME,default=5s"`
	ConnectionTimeout time.Duration `env:"MONGODB_CONNECTION_TIMEOUT,default=10s"`
	Timeout           time.Duration `env:"MONGODB_TIMEOUT,default=3s"`
}

type Security struct {
	JWTExpirationHours time.Duration `env:"JWT_EXPIRATION_HOURS,default=2h"`
	BcryptCost         int           `env:"BCRYPT_COST,default=12"`
	CookieName         string        `env:"COOKIE_NAME,default=auth_token"`
}

type Cors struct {
	AllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS,default=*"`
	AllowedMethods []string `env:"CORS_ALLOWED_METHODS,default=GET|POST|PUT|DELETE|OPTIONS"`
	AllowedHeaders []string `env:"CORS_ALLOWED_HEADERS,default=Content-Type,Authorization"`
}

type RateLimit struct {
	MaxRequests int           `env:"RATE_LIMIT_MAX_REQUESTS,default=100"`
	Window      time.Duration `env:"RATE_LIMIT_WINDOW,default=1m"`
}

type SMTP struct {
	User string `env:"SMTP_USER,required"`
	Pass string `env:"SMTP_PASS,required"`
	Host string `env:"SMTP_HOST,required"`
	Port int    `env:"SMTP_PORT,default=587"`
}

type URLs struct {
	APIBaseURL     string `env:"API_BASE_URL,required"`
	ClientLoginURL string `env:"CLIENT_LOGIN_URL,required"`
}

type OTP struct {
	ExpirationMinutes    time.Duration `env:"OTP_EXPIRATION_MINUTES,default=10m"`
	JWTExpirationMinutes time.Duration `env:"JWT_EXPIRATION_MINUTES,default=10m"`
}

type Key struct {
	PrivateKey string
	PublicKey  string
}
