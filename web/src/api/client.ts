import Auth from './auth'
import {serviceInstance} from './base'
import config from '@/config'

const api = serviceInstance(config.base_api_url, Auth)
export default api