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
            title: 'Error',
            message: err.message + '.\n' + err.hint,
            type: 'error',
            duration: 2000
        })
    }

    public showMessage(message: string) {
        this.$notify({
            title: `Success`,
            message: message,
            duration: 1500,
        })
    }

    private sumFromCents(cents: number): number {
        return (cents / 100)
    }
}