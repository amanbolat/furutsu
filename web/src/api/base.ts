import axios, {AxiosError, AxiosRequestConfig} from 'axios'
import eventBust from '../utils/event_bus'
import store from '@/store/index'

function axiosConfig(baseUrl: string): AxiosRequestConfig {
    return {
        baseURL: baseUrl,
        method: 'POST',
        timeout: 30000,
        responseType: 'json',
        withCredentials: true
    }
}

export function serviceInstance(baseUrl: string) {
    const instance = axios.create(axiosConfig(baseUrl))

    instance.interceptors.request.use((cfg) => {
        cfg.headers.Authorization = `Bearer ${store.getters.authToken}`
        return cfg
    })

    instance.interceptors.response.use((response: any) => {
        return response
    }, (error: AxiosError) => {
        if (error.response && error.response.status === 401) {
            eventBust.$emit('unauthorized_request', !!store.getters.isAuthenticated)
            store.dispatch('Logout').catch((err) => {
                console.log('failed to dispatch Logout', err)
            })
            throw error
        }

        const err = {message: 'Internet connection problem'} as AppError
        if (error.response && error.response.data) {
            err.message = error.response.data.message
            err.hint = error.response.data.hint
        }

        throw err
    })

    return instance
}

export function authServiceInstance(baseUrl: string) {
    const instance = axios.create(axiosConfig(baseUrl))

    instance.interceptors.response.use(undefined, (error: AxiosError) => {
        return new Promise((resolve, reject) => {
            const err = {} as AppError
            err.message = "Internet connection problem"
            if (error.response && error.response.data) {
                err.message = error.response.data.message
                err.hint = error.response.data.hint
            }

            reject(err)
        })
    })


    return instance
}

export interface AppError {
    message: string;
    hint: string;
}
