import Vue from 'vue'
import {Component} from 'vue-property-decorator'
import {AppError} from '@/api/base'

@Component
export default class AppMixin extends Vue {
    public showError(error: AppError) {
        let err = {} as AppError
        if (!error.message) {
            err.message = 'Unknown error'
        } else {
            err = error
        }
        this.$notify({
            group: 'notify',
            title: 'Error',
            type: 'error',
            text: err.message + '.\n' + (err.hint ? err.hint : ''),
            duration: 2000
        })
    }

    public showMessage(message: string) {
        this.$notify({
            group: 'notify',
            title: 'Success',
            type: 'success',
            text: message,
            duration: 3500
        })
    }

    private sumFromCents(cents: number): number {
        return (cents / 100)
    }
}