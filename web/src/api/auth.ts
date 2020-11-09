import {AxiosError, AxiosInstance} from 'axios'
import {authServiceInstance} from './base'
import config from '../config'

export interface User {
    id: string;
    username: string;
    password: string;
    full_name: string;
}

export class AuthManager {
    protected authTokenKey: string = 'auth_token'
    protected axiosInstance: AxiosInstance

    constructor(authTokenKey: string, baseUrl: string) {
        this.authTokenKey = authTokenKey
        this.axiosInstance = authServiceInstance(baseUrl)
    }

    public login(user: string, pass: string) {
        const credentials = {
            username: user,
            password: pass
        }

        return this.axiosInstance.post('/auth/login', credentials)
            .then((response: any) => {
                return response.data.data
            }).catch((error: AxiosError) => {
                console.log(error)
                return Promise.reject(error)
            })
    }
}

const srv = new AuthManager(config.authTokenKey, config.base_api_url)

export default srv
