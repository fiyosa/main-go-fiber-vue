const secret = {
  APP_ENV: import.meta.env.VITE_APP_ENV || 'development',
  API_URL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
}

export default secret
