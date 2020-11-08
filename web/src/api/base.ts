import axios, {AxiosError, AxiosRequestConfig} from 'axios'
import eventBust from '../utils/event_bus'
import {AuthManager} from './auth'
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

export function serviceInstance(baseUrl: string, authManager: AuthManager) {
    const instance = axios.create(axiosConfig(baseUrl))

    instance.interceptors.request.use((cfg) => {
        cfg.headers.Authorization = `Bearer ${store.getters.authToken}`
        return cfg
    })

    instance.interceptors.response.use(undefined, (error: AxiosError) => {
        return new Promise(async (resolve, reject) => {
            if (error.response && error.response.status === 401) {
                await store.dispatch('Logout')
                eventBust.$emit('unauthorized_request')
                return
            }

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
