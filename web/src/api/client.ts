import {serviceInstance} from './base'
import config from '@/config'

const api = serviceInstance(config.base_api_url)
export default api