console.log(`App environment: ${process.env.VUE_APP_ENV}`)

const config = {
    base_api_url: '',
    authTokenKey: 'auth_token',
}

switch (process.env.VUE_APP_ENV) {
    case 'prod':
        config.base_api_url = ''
        break
    case 'test':
        config.base_api_url = 'http://localhost:9033/'
        break
    case 'dev':
        config.base_api_url = ''
        break
}

export default config
