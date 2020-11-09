<template lang="pug">
  v-container(fluid)
    v-row(dense).justify-center
      v-col(cols="12" md="6")
        h3 Make payment
        span Order ID: {{ orderId }}
        v-form.mt-4(v-model="validForm")
          v-text-field(v-model="cardData.number" label="Card number" outlined dense)
          v-text-field(v-model="cardData.holder" label="Holder name" outlined dense)
          v-container
            v-row
              v-col
                v-text-field(v-model="cardData.cvc" label="CVC" dense type="number")
              v-col
                v-select(v-model="cardData.year" label="Year" dense :items="availableYears" type="number")
              v-col
                v-select(v-model="cardData.month" label="Month" dense :items="availableMonths" type="number")
          v-btn(color="orange darken-2" @click="makePayment").white--text Submit
</template>

<script lang="ts">
import {Component, Mixins, Vue} from 'vue-property-decorator'
import api from '@/api/client'
import AppMixin from '@/mixins/AppMixin'

@Component
export default class PaymentForm extends Mixins(AppMixin) {
  private orderId = ''
  private validForm = false
  private cardData = {
    number: '',
    year: '',
    month: '',
    cvc: '',
    holder: ''
  }

  get availableMonths() {
    return Array.from({length: 12}, (v, k) => k + 1)
  }

  get availableYears() {
    return Array.from({length: 30}, (v, k) => k + 2021)
  }

  private created() {
    this.orderId = this.$route.params.order_id
  }

  private makePayment() {
    api.post(`/order/${this.orderId}/pay`, {
      number: this.cardData.number,
      year: parseInt(this.cardData.year),
      month: parseInt(this.cardData.month),
      cvc: parseInt(this.cardData.cvc),
      holder: this.cardData.holder
    }).then(() => {
      this.showMessage('Payment succeed!')
      this.$router.push('/order')
    }).catch(() => {
      this.showError({message: 'Payment failed', hint: 'Your payment was rejected'})
    })
  }
}
</script>