import Vue from 'vue'
import dayjs from 'dayjs'

Vue.filter('parseDate', (date?: string) => {
    if (date) {
        return dayjs(date).format('D MMM YYYY')
    }
    return ''
})

Vue.filter('parseTimestamp', (date?: string) => {
    if (date) {
        return dayjs(date).format('DD.MM.YYYY HH:mm:ss')
    }
    return ''
})
