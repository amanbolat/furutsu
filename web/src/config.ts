console.log(`App environment: ${process.env.VUE_APP_ENV}`)

const config = {
    base_api_url: '',
    authTokenKey: 'auth_token',
}

switch (process.env.VUE_APP_ENV) {
    case 'prod':
        config.base_api_url = process.env.VUE_APP_SERVER_URL
        break
    case 'test':
        config.base_api_url = 'http://localhost:9033/'
        break
    case 'dev':
        config.base_api_url = process.env.VUE_APP_SERVER_URL
        break
    default:
        config.base_api_url = process.env.VUE_APP_SERVER_URL
}

export default config
