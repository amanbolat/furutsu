<template lang="pug">
  v-container(fluid)
    v-row(dense).justify-center
      v-col(cols="12" md="6")
        h3 Make payment
        span Order ID: {{ orderId }}
        v-form.mt-4(v-model="validForm")
          v-text-field(label="Card number" outlined dense)
          v-text-field(label="Holder name" outlined dense)
          v-container
            v-row
              v-col
                v-text-field(label="CVC" dense type="number")
              v-col
                v-select(label="Year" dense :items="availableYears")
              v-col
                v-select(label="Month" dense :items="availableMonths")
          v-btn(color="indigo").white--text Submit
</template>

<script lang="ts">
import {Component, Vue} from 'vue-property-decorator'
import api from '@/api/client'
import eventBus from '@/utils/event_bus'

@Component
export default class OrderForm extends Vue {
  private orderId = ''
  private validForm = false
  private cardData = {
    number: '',
    holder: '',
    month: '',
    year: '',
    cvc: ''
  }

  get availableMonths() {
    return Array.from({length: 12}, (v, k) => k + 1 )
  }

  get availableYears() {
    return Array.from({length: 30}, (v, k) => k + 2021 )
  }

  private created() {
    this.orderId = this.$route.params.order_id
  }

  private makePayment() {
    if (!this.validForm) {
      return
    }

    api.post(`/payment/pay/${this.orderId}`).then((response) => {
      console.log(response)
    }).catch((err) => {
      eventBus.$emit('app_err', err)
    })
  }
}
</script>